package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"go.vocdoni.io/dvote/api"
	"go.vocdoni.io/dvote/crypto/ethereum"
	"go.vocdoni.io/dvote/crypto/zk"
	"go.vocdoni.io/dvote/crypto/zk/circuit"
	"go.vocdoni.io/dvote/httprouter/apirest"
	"go.vocdoni.io/dvote/log"
)

const (
	// HTTPGET is the method string used for calling Request()
	HTTPGET = "GET"
	// HTTPPOST is the method string used for calling Request()
	HTTPPOST = "POST"
	// HTTPDELETE is the method string used for calling
	HTTPDELETE = "DELETE"

	errCodeNot200 = "API error"
)

// HTTPclient is the Vocdoni API HTTP client.
type HTTPclient struct {
	c       *http.Client
	token   *uuid.UUID
	addr    *url.URL
	account *ethereum.SignKeys
	chainID string
	circuit circuit.ZkCircuitConfig
	zkAddr  *zk.ZkAddress
}

// NewHTTPclient creates a new HTTP(s) API Vocdoni client.
func NewHTTPclient(addr *url.URL, bearerToken *uuid.UUID) (*HTTPclient, error) {
	tr := &http.Transport{
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: false,
		WriteBufferSize:    1 * 1024 * 1024, // 1 MiB
		ReadBufferSize:     1 * 1024 * 1024, // 1 MiB
	}
	c := &HTTPclient{
		c:     &http.Client{Transport: tr, Timeout: time.Second * 8},
		token: bearerToken,
		addr:  addr,
	}
	data, status, err := c.Request(HTTPGET, nil, "chain", "info")
	if err != nil {
		return nil, err
	}
	if status != apirest.HTTPstatusOK {
		return nil, fmt.Errorf("%s: %d (%s)", errCodeNot200, status, data)
	}
	info := &api.ChainInfo{}
	if err := json.Unmarshal(data, info); err != nil {
		return nil, fmt.Errorf("cannot get chain ID from API server")
	}
	c.chainID = info.ID

	// Get the default circuit config
	circuitConf, exists := circuit.CircuitsConfigurations[info.CircuitConfigurationTag]
	if !exists {
		return nil, fmt.Errorf("empty or wrong circui configuration tag provided")
	}
	c.circuit = circuitConf

	return c, nil
}

// ChainID returns the chain identifier name in which the API backend is connected.
func (c *HTTPclient) ChainID() string {
	return c.chainID
}

// SetAccount sets the Vocdoni account used for signing transactions and assign
// a new ZkAddress based on the provided private key account.
func (c *HTTPclient) SetAccount(accountPrivateKey string) error {
	c.account = new(ethereum.SignKeys)
	err := c.account.AddHexKey(accountPrivateKey)
	if err != nil {
		return err
	}

	c.zkAddr, err = zk.AddressFromString(accountPrivateKey)
	return err
}

// Clone returns a copy of the HTTPclient with the accountPrivateKey set as the account key.
// Panics if the accountPrivateKey is not valid.
func (c *HTTPclient) Clone(accountPrivateKey string) *HTTPclient {
	clone := *c
	clone.account = new(ethereum.SignKeys)
	if err := clone.account.AddHexKey(accountPrivateKey); err != nil {
		panic(err)
	}
	return &clone
}

// MyAddress returns the address of the account used for signing transactions.
func (c *HTTPclient) MyAddress() common.Address {
	return c.account.Address()
}

// MyZkAddress returns the zkAddress of the current account used for anonymous
// voting
func (c *HTTPclient) MyZkAddress() *zk.ZkAddress {
	return c.zkAddr
}

// SetAuthToken configures the bearer authentication token.
func (c *HTTPclient) SetAuthToken(token *uuid.UUID) {
	c.token = token
}

// SetHostAddr configures the host address of the API server.
func (c *HTTPclient) SetHostAddr(addr *url.URL) error {
	c.addr = addr
	data, status, err := c.Request(HTTPGET, nil, "chain", "info")
	if err != nil {
		return err
	}
	if status != apirest.HTTPstatusOK {
		return fmt.Errorf("%s: %d (%s)", errCodeNot200, status, data)
	}
	info := &api.ChainInfo{}
	if err := json.Unmarshal(data, info); err != nil {
		return fmt.Errorf("cannot get chain ID from API server")
	}
	c.chainID = info.ID
	return nil
}

// Request performs a `method` type raw request to the endpoint specified in urlPath parameter.
// Method is either GET or POST. If POST, a JSON struct should be attached.  Returns the response,
// the status code and an error.
func (c *HTTPclient) Request(method string, jsonBody any, urlPath ...string) ([]byte, int, error) {
	body, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, 0, err
	}
	u, err := url.Parse(c.addr.String())
	if err != nil {
		return nil, 0, err
	}
	u.Path = path.Join(u.Path, path.Join(urlPath...))
	headers := http.Header{}
	if c.token != nil {
		headers = http.Header{
			"Authorization": []string{"Bearer " + c.token.String()},
			"User-Agent":    []string{"Vocdoni API client / 1.0"},
			"Content-Type":  []string{"application/json"},
		}
	}

	log.Debugw("http request", "type", method, "path", u.Path, "body", jsonBody)
	resp, err := c.c.Do(&http.Request{
		Method: method,
		URL:    u,
		Header: headers,
		Body: func() io.ReadCloser {
			if jsonBody == nil {
				return nil
			}
			return io.NopCloser(bytes.NewBuffer(body))
		}(),
	})
	if err != nil {
		return nil, 0, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return data, resp.StatusCode, nil
}

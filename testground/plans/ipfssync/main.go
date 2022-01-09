// This plan demonstrates chaging the network on a per-host or per-net basis to make some nodes
// unreachable to a subset of the test. The nodes are in three "regions", which you might imagine to
// be countries with restrictive policies, corporate firewalls, misconfigured routers, etc. In this
// plan, all the nodes in "regionA" cannot reach "regionB" because the network between them is
// broken. We should expect to see nodes in "regionC" can reach all nodes on the network, while A
// and B cannot.
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/testground/sdk-go/network"
	"github.com/testground/sdk-go/ptypes"
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
	"github.com/testground/sdk-go/sync"

	flag "github.com/spf13/pflag"

	"go.vocdoni.io/dvote/crypto/ethereum"
	"go.vocdoni.io/dvote/data"
	"go.vocdoni.io/dvote/ipfssync"
	"go.vocdoni.io/dvote/log"
)

type region int

const (
	regionA = iota
	regionB
	regionC
)

const RunDuration = 5 * time.Minute

func (r region) String() string {
	return [...]string{"region_A", "region_B", "region_C"}[r]
}

type node struct {
	Region region
	IP     *net.IP
}

func main() {
	testcases := map[string]interface{}{
		"drop":   routeFilter(network.Drop),
		"reject": routeFilter(network.Reject),
		"accept": routeFilter(network.Accept),
	}
	run.InvokeMap(testcases)
}

func expectErrors(runenv *runtime.RunEnv, a *node, b *node) bool {
	if runenv.TestCase == "accept" || a.Region == regionC || b.Region == regionC {
		return false
	}
	if (a.Region == regionA && b.Region == regionB) || (a.Region == regionB && b.Region == regionA) {
		return true
	}
	return false
}

func routeFilter(action network.FilterAction) run.TestCaseFn {

	return func(runenv *runtime.RunEnv) error {

		ctx, cancel := context.WithTimeout(context.Background(), RunDuration+300*time.Second)
		defer cancel()

		client := sync.MustBoundClient(ctx, runenv)

		if !runenv.TestSidecar {
			return fmt.Errorf("this plan must be run with sidecar enabled")
		}

		netclient := network.NewClient(client, runenv)
		netclient.MustWaitNetworkInitialized(ctx)

		config := &network.Config{
			// Control the "default" network. At the moment, this is the only network.
			Network: "default",

			// Enable this network. Setting this to false will disconnect this test
			// instance from this network. You probably don't want to do that.
			Enable: true,
			Default: network.LinkShape{
				Latency:   100 * time.Millisecond,
				Bandwidth: 1 << 20, // 1Mib
			},
			CallbackState: "network-configured",
			RoutingPolicy: network.AllowAll,
		}

		runenv.RecordMessage("before netclient.MustConfigureNetwork")
		netclient.MustConfigureNetwork(ctx, config)

		// Start ipfsSync node
		runenv.RecordMessage("Starting ipfsSync node")

		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		userDir := home + "/.ipfs"
		logLevel := flag.String("logLevel", "info", "log level")
		dataDir := flag.String("dataDir", userDir, "directory for storing data")
		key := flag.String("key", "vocdoni", "secret shared group key for the sync cluster")
		nodeKey := flag.String("nodeKey", "", "custom private hexadecimal 256 bit key for p2p identity")
		port := flag.Int16("port", 4171, "port for the sync network")
		helloInterval := flag.Int("helloInterval", 40, "period in seconds for sending hello messages")
		updateInterval := flag.Int("updateInterval", 20, "period in seconds for sending update messages")
		peers := flag.StringArray("peers", []string{},
			"custom list of peers to connect to (multiaddresses separated by commas)")
		private := flag.Bool("private", false,
			"if enabled a private libp2p network will be created (using the secret key at transport layer)")
		bootnodes := flag.StringArray("bootnodes", []string{},
			"list of bootnodes (multiaddress separated by commas)")
		bootnode := flag.Bool("bootnode", false,
			"act as a bootstrap node (will not try to connect with other bootnodes)")

		flag.CommandLine.SortFlags = false
		flag.Parse()

		log.Init(*logLevel, "stdout")
		ipfsStore := data.IPFSNewConfig(*dataDir)
		storage, err := data.Init(data.StorageIDFromString("IPFS"), ipfsStore)
		if err != nil {
			log.Fatal(err)
		}

		sk := ethereum.NewSignKeys()
		var privKey string

		if len(*nodeKey) > 0 {
			if err := sk.AddHexKey(*nodeKey); err != nil {
				log.Fatal(err)
			}
			_, privKey = sk.HexString()
		} else {
			pk := make([]byte, 64)
			kfile, err := os.OpenFile(*dataDir+"/.ipfsSync.key", os.O_CREATE|os.O_RDWR, 0o600)
			if err != nil {
				log.Fatal(err)
			}

			if n, err := kfile.Read(pk); err != nil || n == 0 {
				log.Info("generating new node private key")
				if err := sk.Generate(); err != nil {
					log.Fatal(err)
				}
				_, privKey = sk.HexString()
				if _, err := kfile.WriteString(privKey); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Info("loaded saved node private key")
				if err := sk.AddHexKey(string(pk)); err != nil {
					log.Fatal(err)
				}
				_, privKey = sk.HexString()
			}
			if err := kfile.Close(); err != nil {
				log.Fatal(err)
			}
		}

		p2pType := "libp2p"
		if *private {
			p2pType = "privlibp2p"
		}

		is := ipfssync.NewIPFSsync(*dataDir, *key, privKey, p2pType, storage)
		is.HelloInterval = time.Second * time.Duration(*helloInterval)
		is.UpdateInterval = time.Second * time.Duration(*updateInterval)
		is.Port = *port
		if *bootnode {
			is.Bootnodes = []string{""}
		} else {
			is.Bootnodes = *bootnodes
		}
		is.Start()
		for _, peer := range *peers {
			time.Sleep(2 * time.Second)
			log.Infof("connecting to peer %s", peer)
			if err := is.Transport.AddPeer(peer); err != nil {
				log.Warnf("cannot connect to custom peer: (%s)", err)
			}
		}

		// Race to signal this point, the sequence ID determines to which region this node belongs.
		seq := client.MustSignalEntry(ctx, "region-select")
		ip := netclient.MustGetDataNetworkIP()
		me := node{region(int(seq) % 3), &ip}
		runenv.RecordMessage("my ip is %s and I am in region %s", ip, me.Region)

		// instead of blocking forever, sleep for RunDuration
		time.Sleep(RunDuration)

		// publish my address so other nodes know how to reach me.
		nodeTopic := sync.NewTopic("nodes", node{})
		nodeCh := make(chan *node)
		_, _ = client.MustPublishSubscribe(ctx, nodeTopic, &me, nodeCh)

		// Wait until we have received all addresses
		nodes := make([]*node, 0)
		for found := 1; found <= runenv.TestInstanceCount; found++ {
			n := <-nodeCh
			runenv.RecordMessage("received node (%s) %s", n.Region.String(), n.IP.String())
			if !me.IP.Equal(*n.IP) {
				nodes = append(nodes, n)
			}
		}

		// nodes from regionA apply a network policy for the nodes in regionB
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}
		if me.Region == regionA {
			cfg := network.Config{
				Network:        "default",
				CallbackState:  sync.State("reconfigured" + hostname),
				CallbackTarget: 1,
				Enable:         true,
			}

			for _, p := range nodes {
				if p.Region == regionB {
					pnet := ptypes.IPNet{
						IPNet: net.IPNet{
							IP:   *p.IP,
							Mask: net.IPMask([]byte{255, 255, 255, 255}),
						},
					}
					cfg.Rules = append(cfg.Rules, network.LinkRule{
						Subnet: pnet,
						LinkShape: network.LinkShape{
							Filter: action,
						},
					})
				}
			}
			netclient.MustConfigureNetwork(ctx, &cfg)
		}

		runenv.RecordMessage("waiting for all nodes to receive all addresses")
		// Wait until *all* nodes have received all addresses.
		_, err = client.SignalAndWait(ctx, "nodeRoundup", runenv.TestInstanceCount)
		if err != nil {
			return err
		}

		// The http doesn't start instantly, just hang on a sec.
		time.Sleep(10 * time.Second)

		client.MustSignalAndWait(ctx, "testcomplete", runenv.TestInstanceCount)

		client.Close()
		return nil
	}
}

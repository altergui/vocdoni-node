package zk

import (
	"math/big"
	"testing"

	qt "github.com/frankban/quicktest"
	"go.vocdoni.io/dvote/crypto/zk/prover"
	"go.vocdoni.io/proto/build/go/models"
)

func TestProtobufZKProofToProverProof(t *testing.T) {
	c := qt.New(t)

	badInput := &models.ProofZkSNARK{
		A:            []string{},
		B:            []string{},
		C:            []string{},
		PublicInputs: []string{},
	}
	_, err := ProtobufZKProofToProverProof(badInput)
	c.Assert(err, qt.IsNotNil)

	input := &models.ProofZkSNARK{
		A:            []string{"0", "1", "2"},
		B:            []string{"0", "1", "2", "3", "4", "5"},
		C:            []string{"0", "1", "2"},
		PublicInputs: []string{"0", "1", "2"},
	}
	expected := &prover.Proof{
		Data: prover.ProofData{
			A: []string{"0", "1", "2"},
			B: [][]string{
				{"0", "1"},
				{"2", "3"},
				{"4", "5"},
			},
			C: []string{"0", "1", "2"},
		},
		PubSignals: []string{"0", "1", "2"},
	}
	result, err := ProtobufZKProofToProverProof(input)
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.DeepEquals, expected)
}

func TestProverProofToProtobufZKProof(t *testing.T) {
	c := qt.New(t)

	badInput := &prover.Proof{
		Data:       prover.ProofData{},
		PubSignals: []string{},
	}
	_, err := ProverProofToProtobufZKProof(0, badInput, nil, nil, nil, nil)
	c.Assert(err, qt.IsNotNil)

	input := &prover.Proof{
		Data: prover.ProofData{
			A: []string{"0", "1", "2"},
			B: [][]string{
				{"0", "1"},
				{"2", "3"},
				{"4", "5"},
			},
			C: []string{"0", "1", "2"},
		},
		PubSignals: []string{},
	}
	_, err = ProverProofToProtobufZKProof(0, input, nil, nil, nil, nil)
	c.Assert(err, qt.IsNotNil)

	expected := &models.ProofZkSNARK{
		A: []string{"0", "1", "2"},
		B: []string{"0", "1", "2", "3", "4", "5"},
		C: []string{"0", "1", "2"},
		PublicInputs: []string{
			"0", "0", "0", "0", "1",
			"302689215824177652345211539748426020171",
			"205062086841587857568430695525160476881",
		},
	}
	mockData := make([]byte, 32)
	result, err := ProverProofToProtobufZKProof(0, input, mockData, mockData, mockData, new(big.Int).SetInt64(1))
	c.Assert(err, qt.IsNil)
	c.Assert(result.A, qt.ContentEquals, expected.A)
	c.Assert(result.B, qt.ContentEquals, expected.B)
	c.Assert(result.C, qt.ContentEquals, expected.C)
	c.Assert(result.PublicInputs, qt.ContentEquals, expected.PublicInputs)

	input = &prover.Proof{
		Data: prover.ProofData{
			A: []string{"0", "1", "2"},
			B: [][]string{
				{"0", "1"},
				{"2", "3"},
				{"4", "5"},
			},
			C: []string{"0", "1", "2"},
		},
		PubSignals: []string{"0", "1", "2", "3", "4", "5", "6"},
	}
	expected = &models.ProofZkSNARK{
		A:            []string{"0", "1", "2"},
		B:            []string{"0", "1", "2", "3", "4", "5"},
		C:            []string{"0", "1", "2"},
		PublicInputs: []string{"0", "1", "2", "3", "4", "5", "6"},
	}
	result, err = ProverProofToProtobufZKProof(0, input, mockData, mockData, mockData, new(big.Int).SetInt64(1))
	c.Assert(err, qt.IsNil)
	c.Assert(result.A, qt.ContentEquals, expected.A)
	c.Assert(result.B, qt.ContentEquals, expected.B)
	c.Assert(result.C, qt.ContentEquals, expected.C)
	c.Assert(result.PublicInputs, qt.ContentEquals, expected.PublicInputs)
}

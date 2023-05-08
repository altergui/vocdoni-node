package indexer

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"go.vocdoni.io/dvote/crypto/ethereum"
	"go.vocdoni.io/dvote/log"
	"go.vocdoni.io/dvote/util"
	"go.vocdoni.io/dvote/vochain"
	"go.vocdoni.io/dvote/vochain/state"
	"go.vocdoni.io/dvote/vochain/transaction/vochaintx"
	"go.vocdoni.io/proto/build/go/models"
)

func BenchmarkIndexVotes(b *testing.B) {
	app := vochain.TestBaseApplication(b)

	idx := newTestIndexer(b, app, true)
	pid := util.RandomBytes(32)
	if err := app.State.AddProcess(&models.Process{
		ProcessId:    pid,
		EntityId:     util.RandomBytes(20),
		EnvelopeType: &models.EnvelopeType{EncryptedVotes: false},
		Status:       models.ProcessStatus_READY,
		BlockCount:   100000000,
		VoteOptions:  &models.ProcessVoteOptions{MaxCount: 3, MaxValue: 1},
		Mode:         &models.ProcessMode{AutoStart: true},
	}); err != nil {
		b.Fatal(err)
	}
	app.AdvanceTestBlock()

	s := new(ethereum.SignKeys)
	err := s.Generate()
	qt.Assert(b, err, qt.IsNil)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := int32(0); j < 2000; j++ {
			vote := &state.Vote{
				Height:      uint32(util.RandomInt(10, 10000)),
				ProcessID:   pid,
				Nullifier:   util.RandomBytes(32),
				VotePackage: []byte(`{"votes": [1,0,1]}`),
				Weight:      new(big.Int).SetUint64(uint64(util.RandomInt(1, 10000))),
			}
			idx.OnVote(vote, j)
		}
		app.AdvanceTestBlock()
	}
}

func BenchmarkFetchTx(b *testing.B) {
	numTxs := 1000
	app := vochain.TestBaseApplication(b)

	idx := newTestIndexer(b, app, true)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < numTxs; j++ {
			idx.OnNewTx(&vochaintx.VochainTx{TxID: util.Random32()}, uint32(i), int32(j))
		}
		err := idx.Commit(uint32(i))
		qt.Assert(b, err, qt.IsNil)

		time.Sleep(time.Second * 2)

		startTime := time.Now()
		for j := 0; j < numTxs; j++ {
			_, err = idx.GetTxReference(uint64((i * numTxs) + j + 1))
			qt.Assert(b, err, qt.IsNil)
		}
		log.Infof("fetched %d transactions (out of %d total) by index, took %s",
			numTxs, (i+1)*numTxs, time.Since(startTime))
		startTime = time.Now()
		for j := 0; j < numTxs; j++ {
			_, err = idx.GetTxHashReference([]byte(fmt.Sprintf("hash%d%d", i, j)))
			qt.Assert(b, err, qt.IsNil)
		}
		log.Infof("fetched %d transactions (out of %d total) by hash, took %s",
			numTxs, (i+1)*numTxs, time.Since(startTime))
	}
}

func BenchmarkNewProcess(b *testing.B) {
	app := vochain.TestBaseApplication(b)

	idx := newTestIndexer(b, app, true)
	_ = idx // used via the callbacks; we want to benchmark it too
	startTime := time.Now()
	numProcesses := b.N
	entityID := util.RandomBytes(20)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < numProcesses; i++ {
		pid := util.RandomBytes(32)
		if err := app.State.AddProcess(&models.Process{
			ProcessId: pid,
			// EntityId:     util.RandomBytes(20),
			EntityId:     entityID,
			EnvelopeType: &models.EnvelopeType{EncryptedVotes: false},
			Status:       models.ProcessStatus_READY,
			BlockCount:   100000000,
			VoteOptions:  &models.ProcessVoteOptions{MaxCount: 3, MaxValue: 1},
			Mode:         &models.ProcessMode{AutoStart: true},
		}); err != nil {
			b.Fatal(err)
		}
	}
	app.AdvanceTestBlock()
	log.Infof("indexed %d new processes, took %s",
		numProcesses, time.Since(startTime))
}

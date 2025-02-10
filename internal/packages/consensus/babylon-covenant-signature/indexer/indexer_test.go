package indexer

import (
	"testing"
	"time"

	"github.com/cosmostation/cvms/internal/common"
	"github.com/cosmostation/cvms/internal/helper/logger"
	"github.com/cosmostation/cvms/internal/packages/consensus/babylon-covenant-signature/repository"

	"github.com/stretchr/testify/assert"
)

var (
	syncStartHeight int64 = 245219

	p = common.Packager{
		ChainName:    "babylon",
		ChainID:      "bbn-test-5",
		ProtocolType: "cosmos",
		Endpoints: common.Endpoints{
			RPCs: []string{"https://rpc-office.cosmostation.io/babylon-testnet"},
			APIs: []string{"https://lcd-office.cosmostation.io/babylon-testnet"},
		},
		Logger: logger.GetTestLogger(),
	}
)

func TestStart(t *testing.T) {
	waitingDuration := 60 * time.Second
	// Step 1: Set up the database
	tempDBName := "temp"
	indexerDB, err := common.NewTestLoaclIndexerDB(tempDBName)
	assert.NoError(t, err)

	p.SetIndexerDB(indexerDB)
	p.SetRetentionTime("1h")
	p.IsConsumerChain = false

	// Step 2: Initialize the CovenantSignatureIndexer
	idx, err := NewCovenantSignatureIndexer(p)
	assert.NoError(t, err)

	//test target height
	idx.earliestBlockHeight = syncStartHeight

	// Modify Start() to accept a callback for testing purposes
	err = idx.Start()
	assert.NoError(t, err)

	// Step 4: Wait for the goroutine to finish
	done := make(chan bool)

	go func() {
		// Simulate some work
		time.Sleep(waitingDuration)
		done <- true
	}()

	<-done
	t.Log("Goroutine finished work")
}

func TestBatchSync(t *testing.T) {
	waitingDuration := 60 * time.Second
	tempDBName := "temp"
	indexerDB, err := common.NewTestLoaclIndexerDB(tempDBName)
	p.SetIndexerDB(indexerDB)
	p.IsConsumerChain = false
	assert.NoError(t, err)

	idx, err := NewCovenantSignatureIndexer(p)
	assert.NoError(t, err)

	err = idx.InitChainInfoID()
	assert.NoError(t, err)

	err = idx.repo.InitPartitionTablesByChainInfoID(repository.IndexName, idx.ChainID, 100)
	assert.NoError(t, err)

	// Step 4: Wait for the goroutine to finish
	done := make(chan bool)

	newIndexPointer, err := idx.batchSync(idx.earliestBlockHeight, idx.earliestBlockHeight+1)

	go func() {
		// Simulate some work
		time.Sleep(waitingDuration)
		done <- true
	}()

	assert.NoError(t, err)
	t.Logf("new index point: %d", newIndexPointer)
}

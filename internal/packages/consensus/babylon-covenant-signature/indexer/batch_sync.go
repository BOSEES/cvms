package indexer

import (
	"context"
	"time"

	"sync"

	indexertypes "github.com/cosmostation/cvms/internal/common/indexer/types"
	"github.com/cosmostation/cvms/internal/helper"
	"github.com/cosmostation/cvms/internal/packages/consensus/babylon-covenant-signature/model"
)

// NOTE: babylon covenant signature will be created at random block (When Delegation TX occurs)
// first we can get the block txs from the chain
// 1. get block txs /cosmos/tx/v1beta1/txs/block/{height}
// and then, query the first decoded tx in the block
// 2. filtering "/babylon.btcstaking.v1.MsgAddCovenantSigs" message in block txs
// 3. make a list for covenantSigs of committee members
func (idx *CovenantSignatureIndexer) batchSync(lastIndexPointerHeight, newIndexPointerHeight int64) (
	/* new index pointer */ int64,
	/* error */ error,
) {
	if lastIndexPointerHeight >= idx.Lh.LatestHeight {
		idx.Infof("current height is %d and latest height is %d both of them are same, so it'll skip the logic", lastIndexPointerHeight, idx.Lh.LatestHeight)
		return lastIndexPointerHeight, nil
	}

	// set starntHeight and endHeight for batch sync
	startHeight := newIndexPointerHeight
	endHeight := idx.Lh.LatestHeight
	blockSyncSize := idx.Lh.LatestHeight - newIndexPointerHeight

	// set limit at end-height in this batch sync logic
	if (idx.Lh.LatestHeight - newIndexPointerHeight) > indexertypes.BatchSyncLimit {
		endHeight = newIndexPointerHeight + indexertypes.BatchSyncLimit
		blockSyncSize = indexertypes.BatchSyncLimit
		idx.Debugf("by batch sync limit, end height will change to %d", endHeight)
	}

	// init channel and waitgroup for go-routine
	ch := make(chan helper.Result)
	var wg sync.WaitGroup

	// init covenant signature list
	covenantSignatureList := make([][]model.BabylonCovenantSignature, blockSyncSize)
	var endBlockTimestamp time.Time

	for height := startHeight; height < endHeight; height++ {
		wg.Add(1)
		height := height
		index := height - startHeight

		go func(ch chan helper.Result, i int64) {
			defer helper.HandleOutOfNilResponse(idx.Entry)
			defer wg.Done()

			requester := idx.APIClient.R().SetContext(context.Background())
		RETRY:
			resp, err := requester.Get(BlockTxsQueryPath(height))
			if err != nil {
				idx.Errorf("failed to call block results api, %s", err)
				time.Sleep(time.Second * 1)
				goto RETRY
			}

			blockHeight, blockTimestamp, covenantSigs, err := ExtractBabylonCovenantSignature(resp.Body())
			if err != nil {
				idx.Errorln(err, height)
				goto RETRY
			}

			if height == endHeight-1 {
				endBlockTimestamp = blockTimestamp
			}

			var newBcsList = make([]model.BabylonCovenantSignature, 0)
			for _, sig := range covenantSigs {

				// It's not yet clear if Committee members can change dynamically, we've added some temporary code to prevent panic
				pkID, exists := idx.covenantCommitteeMap[sig.Pk]
				if !exists {
					idx.Errorf("Missing covenant committee entry for PK: %s", sig.Pk)
					continue
				}

				newCovenantSignature := model.BabylonCovenantSignature{
					ChainInfoID:      idx.ChainInfoID,
					Height:           blockHeight,
					CovenantBtcPkID:  pkID,
					BTCStakingTxHash: sig.StakingTxHash,
					Timestamp:        blockTimestamp,
				}

				newBcsList = append(newBcsList, newCovenantSignature)
			}

			ch <- helper.Result{
				Item:    newBcsList,
				Success: true,
				Index:   i,
			}
		}(ch, index)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for r := range ch {
		if r.Success {
			item := r.Item.([]model.BabylonCovenantSignature)
			covenantSignatureList[r.Index] = item
			continue
		}
	}

	for _, css := range covenantSignatureList {
		err := idx.repo.InsertBabylonCovenantSignatureList(idx.ChainInfoID, endHeight, css)
		if err != nil {
			return lastIndexPointerHeight, err
		}
	}

	idx.updatePrometheusMetrics(covenantSignatureList, endBlockTimestamp)
	idx.Debugf("updated babylon covenant signature in %v block", endBlockTimestamp)
	return endHeight, nil
}

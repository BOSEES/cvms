package model

import (
	"time"

	"github.com/uptrace/bun"
)

type BabylonCovenantSignature struct {
	bun.BaseModel    `bun:"table:babylon_covenant_signature"`
	ID               int64     `bun:"id,pk,autoincrement"`
	ChainInfoID      int64     `bun:"chain_info_id,pk,notnull"`
	Height           int64     `bun:"height,notnull"`
	CovenantBtcPkID  int64     `bun:"covenant_btc_pk_id,notnull"`
	BTCStakingTxHash string    `bun:"btc_staking_tx_hash,notnull"`
	Timestamp        time.Time `bun:"timestamp,notnull"`
}

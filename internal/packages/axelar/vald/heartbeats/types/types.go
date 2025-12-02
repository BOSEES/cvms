package types

var (
	SupportedChains = []string{"axelar"}
)

const (
	// common
	CommonValidatorQueryPath = "/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED&pagination.count_total=true&pagination.limit=500"

	// axelar
	AxelarChainMaintainersQueryPath = "/axelar/nexus/v1beta1/chain_maintainers/{chain}"
	AxelarProxyResisterQueryPath    = "/axelar/snapshot/v1beta1/proxy?operator_address={validator_operator_address}"
)

type CommonAxelarHeartbeats struct {
	Validators             []BroadcastorStatus
	LatestHeartBeatsHeight int64
}

type BroadcastorStatus struct {
	Moniker                  string  `json:"moniker"`
	ValidatorOperatorAddress string  `json:"validator_operator_address"`
	BroadcastorAddress       string  `json:"broadcastor_address"`
	Status                   string  `json:"status"`
	LatestHeartBeat          float64 `json:"latest_heartbeat"`
}

type CommonValidatorsQueryResponse struct {
	Validators []struct {
		OperatorAddress string `json:"operator_address"`
		Description     struct {
			Moniker string `json:"moniker"`
		} `json:"description"`
	} `json:"validators"`
	Pagination struct {
		NextKey interface{} `json:"-"`
		Total   string      `json:"-"`
	} `json:"-"`
}

type AxelarProxyResisterStatus struct {
	ProxyAddress string `json:"proxy_address"`
	Status       string `json:"status"`
}

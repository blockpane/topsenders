package top_senders

// Message is a combination of the various message types we're interested in
// since each message might have different fields for the sender, we add
// them all and just concatenate them together in CountMessageTypeBySender
type Message struct {
	Type string `json:"@type"`

	Sender           string `json:"sender"`
	Voter            string `json:"voter"`
	Delegator        string `json:"delegator_address"`
	Grantee          string `json:"grantee"`
	FromAddress      string `json:"from_address"`
	Signer           string `json:"signer"`
	Creator          string `json:"creator"`
	ValidatorAddress string `json:"validator_address"`
	CosmosReceiver   string `json:"cosmos_receiver"`
	From             string `json:"from"`
	Requester        string `json:"requester"`
	Initiator        string `json:"initiator"`
	Depositor        string `json:"depositor"`
	Orchestrator     string `json:"orchestrator"`
	Trader           string `json:"trader"`
}

type Tx struct {
	Body struct {
		Messages []Message `json:"messages"`
	} `json:"body"`
}

type Root struct {
	Txs []Tx `json:"txs"`
}

type Status struct {
	Result struct {
		NodeInfo struct {
			Network string `json:"network"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHeight string `json:"latest_block_height"`
		} `json:"sync_info"`
	} `json:"result"`
}

func (s Status) Height() string {
	return s.Result.SyncInfo.LatestBlockHeight
}

package top_senders

import (
	"fmt"
	"github.com/fatih/camelcase"
	"regexp"
	"strings"
)

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

type Graph struct {
	Nodes       []Node `json:"nodes"`
	Edges       []Edge `json:"edges"`
	Description string `json:"description"`
}

func NewGraph(desc string) Graph {
	return Graph{
		Description: desc,
		Nodes:       make([]Node, 0),
		Edges:       make([]Edge, 0),
	}
}

type Node struct {
	ID      string      `json:"id"`
	Label   string      `json:"label"`
	Title   string      `json:"title"`
	Scaling NodeScaling `json:"scaling"`
	Value   int         `json:"value"`
	Shape   string      `json:"shape"`
	Shadow  bool        `json:"shadow"`
}

type NodeScaling struct {
	Min   int  `json:"min"`
	Max   int  `json:"max"`
	Label bool `json:"label"`
}

var accountRex = regexp.MustCompile(`^\w+1\w+$`)

func NewNode(id string, value int) Node {
	var label, shape string
	if accountRex.MatchString(id) {
		shape = "ellipse"
		label = fmt.Sprintf("%s...%s (%d)", id[:8], id[len(id)-4:], value)
	} else {
		shape = "box"
		label = fmt.Sprintf("%s (%d)", strings.Join(camelcase.Split(id), " "), value)
	}
	return Node{
		ID:    id,
		Label: label,
		Title: id,
		Value: value,
		Shape: shape,
		Scaling: NodeScaling{
			Min:   1,
			Max:   15,
			Label: true,
		},
		Shadow: true,
	}
}

type Edge struct {
	From    string      `json:"from"`
	To      string      `json:"to"`
	Value   int         `json:"value"`
	Scaling EdgeScaling `json:"scaling"`
	Label   string      `json:"label"`
	// Title   string      `json:"title"`
	Shadow bool `json:"shadow"`
}

type EdgeScaling struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func NewEdge(from, to string, value int) Edge {
	return Edge{
		From:  from,
		To:    to,
		Value: value,
		Label: fmt.Sprintf("%d", value),
		Scaling: EdgeScaling{
			Min: 1,
			Max: 10,
		},
		Shadow: true,
	}
}

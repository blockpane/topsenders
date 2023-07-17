package top_senders

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fatih/camelcase"
)

var (
	ApiURL   string
	RpcUrl   string
	Interval time.Duration
	LookBack int
	Top      int
)

func Txs() {
	var height, lastheight, chain string
	var count int
	var err error

	combined := make([]map[string]map[string]int, LookBack)

	for {
		height, chain, err = GetBlockHeight(RpcUrl)
		if height == lastheight {
			time.Sleep(Interval)
			continue
		}
		lastheight = height
		if err != nil {
			log.Println(err)
			time.Sleep(Interval)
			continue
		}
		root, err := GetTxs(ApiURL, height)
		if err != nil {
			log.Println(err)
			time.Sleep(Interval)
			continue
		}
		result := CountMessageTypeBySender(root, height)
		combined = append(combined[1:], result)
		result = CombineCountTables(combined)

		resultSorted := make([]map[string]int, 0)
		for sender, msgTypes := range result {
			for msgType, count := range msgTypes {
				row := make(map[string]int)
				row[sender+" "+msgType] = count
				resultSorted = append(resultSorted, row)
			}
		}

		if count < LookBack {
			count += 1
		}

		if len(resultSorted) == 0 {
			fmt.Print("\033[2J")
			fmt.Printf("Top %d accounts by TX type Last %d blocks on %s (%s):\n\n", len(resultSorted), count, chain, height)
			time.Sleep(Interval)
			continue
		}
		sort.Slice(resultSorted, func(i, j int) bool {
			for k := range resultSorted[i] {
				for l := range resultSorted[j] {
					if resultSorted[i][k] == resultSorted[j][l] {
						return sort.StringsAreSorted([]string{k, l})
					}
					return resultSorted[i][k] > resultSorted[j][l]
				}
			}
			return false
		})

		if len(resultSorted) > Top {
			resultSorted = resultSorted[:Top]
		}

		// TODO: this is a hack to clear the screen using an ANSI reset, need to find a better way
		fmt.Print("\033[2J")

		fmt.Printf("Top %d accounts by TX type Last %d blocks on %s (%s):\n\n", len(resultSorted), count, chain, height)
		for _, row := range resultSorted {
			for k, v := range row {
				ks := strings.Split(k, " ")
				if len(ks) == 1 {
					ks = append(ks, "")
				}
				fmt.Printf("%50s %40s     %5d\n", ks[0], strings.Join(camelcase.Split(ks[1]), " "), v)
			}
		}
		fmt.Println()
		time.Sleep(Interval)
	}

}

// CountMessageTypeBySender will return a map of sender addresses to a map of message types to counts
func CountMessageTypeBySender(root Root, height string) (countTable map[string]map[string]int) {
	rex := regexp.MustCompile(`^Aggregate.+ote`)
	countTable = make(map[string]map[string]int)
	for _, tx := range root.Txs {
		for _, message := range tx.Body.Messages {
			sender := message.Sender + message.Delegator + message.Voter +
				message.Grantee + message.FromAddress + message.Signer +
				message.Creator + message.CosmosReceiver + message.From +
				message.Requester + message.Initiator + message.Depositor +
				message.Trader

			t := strings.Split(message.Type, ".")
			var isEvm bool
			if len(t) > 1 {
				if t[1] == "evm" {
					isEvm = true
				}
			}
			msgType := strings.TrimLeft(t[len(t)-1], "Msg")

			// prevent duplicates
			switch {
			case isEvm:
				sender = "unknown_eth_address"
			case sender == "" && message.Orchestrator != "":
				sender = message.Orchestrator
			case sender == "" && message.ValidatorAddress != "":
				sender = message.ValidatorAddress
			case sender == "" && rex.MatchString(msgType):
				sender = "<oracles>"
			case sender == "" && message.ValidatorAddress == "":
				sender = height
			}

			if countTable[sender] == nil {
				countTable[sender] = make(map[string]int)
			}
			countTable[sender][msgType]++
		}
	}
	return countTable
}

// CombineCountTables will combine multiple count tables into a single count table
func CombineCountTables(countTables []map[string]map[string]int) map[string]map[string]int {
	combined := make(map[string]map[string]int)
	for _, countTable := range countTables {
		for sender, msgTypes := range countTable {
			if combined[sender] == nil {
				combined[sender] = make(map[string]int)
			}
			for msgType, count := range msgTypes {
				combined[sender][msgType] += count
			}
		}
	}
	return combined
}

// GetBlockHeight will retrieve the height from the /status endpoint
func GetBlockHeight(url string) (height, network string, err error) {
	url = fmt.Sprintf("%s/status", url)
	client := http.DefaultClient
	client.Timeout = 5 * time.Second
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var status Status
	err = json.Unmarshal(body, &status)
	if err != nil {
		return
	}
	network = status.Result.NodeInfo.Network
	height = status.Height()
	return
}

// GetTxs will retrieve the txs from the /cosmos/tx/v1beta1/txs/block/<height> endpoint
// and return a Root struct
func GetTxs(url, height string) (root Root, err error) {
	url = fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/block/%s", url, height)
	client := http.DefaultClient
	client.Timeout = 5 * time.Second
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusNotImplemented {
		log.Fatalf("Error: required endpoint is not present! Is chain running 0.45.2 or later?\n\n%s returned status code %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &root)
	return
}

package domain

import (
	"regexp"
	"strings"
	"time"

	cmtypes "github.com/cometbft/cometbft/abci/types"
	"github.com/shopspring/decimal"
)

type Tx struct {
	Jsonrpc string  `json:"jsonrpc,omitempty"`
	Id      string  `json:"id,omitempty"`
	Result  *Result `json:"result,omitempty"`
}

type Result struct {
	Query string `json:"query,omitempty"`
	Data  *Data  `json:"data,omitempty"`
}

type Data struct {
	Type  string `json:"type,omitempty"`
	Value *Value `json:"value,omitempty"`
}

type Value struct {
	TxResult            *TxResult            `json:"TxResult,omitempty"`
	Block               *Block               `json:"block"`
	BlockId             *BlockId             `json:"block_id"`
	ResultFinalizeBlock *ResultFinalizeBlock `json:"result_finalize_block"`
}

type TxResult struct {
	Height string    `json:"height,omitempty"`
	Tx     string    `json:"tx,omitempty"`
	Result *ResultTX `json:"result,omitempty"`
}

type Block struct {
	Header struct {
		Version struct {
			Block string `json:"block"`
		} `json:"version"`
		ChainId     string    `json:"chain_id"`
		Height      string    `json:"height"`
		Time        time.Time `json:"time"`
		LastBlockId struct {
			Hash  string `json:"hash"`
			Parts struct {
				Total int    `json:"total"`
				Hash  string `json:"hash"`
			} `json:"parts"`
		} `json:"last_block_id"`
		LastCommitHash     string `json:"last_commit_hash"`
		DataHash           string `json:"data_hash"`
		ValidatorsHash     string `json:"validators_hash"`
		NextValidatorsHash string `json:"next_validators_hash"`
		ConsensusHash      string `json:"consensus_hash"`
		AppHash            string `json:"app_hash"`
		LastResultsHash    string `json:"last_results_hash"`
		EvidenceHash       string `json:"evidence_hash"`
		ProposerAddress    string `json:"proposer_address"`
	} `json:"header"`
	Data struct {
		Txs []string `json:"txs"`
	} `json:"data"`
	Evidence struct {
		Evidence []interface{} `json:"evidence"`
	} `json:"evidence"`
	LastCommit struct {
		Height  string `json:"height"`
		Round   int    `json:"round"`
		BlockId struct {
			Hash  string `json:"hash"`
			Parts struct {
				Total int    `json:"total"`
				Hash  string `json:"hash"`
			} `json:"parts"`
		} `json:"block_id"`
		Signatures []struct {
			BlockIdFlag      int       `json:"block_id_flag"`
			ValidatorAddress string    `json:"validator_address"`
			Timestamp        time.Time `json:"timestamp"`
			Signature        string    `json:"signature"`
		} `json:"signatures"`
	} `json:"last_commit"`
}

type BlockId struct {
	Hash  string `json:"hash"`
	Parts struct {
		Total int    `json:"total"`
		Hash  string `json:"hash"`
	} `json:"parts"`
}

type ResultFinalizeBlock struct {
	Events                []cmtypes.Event `json:"events"`
	ValidatorUpdates      []interface{}   `json:"validator_updates"`
	ConsensusParamUpdates struct {
		Block struct {
			MaxBytes string `json:"max_bytes"`
			MaxGas   string `json:"max_gas"`
		} `json:"block"`
		Evidence struct {
			MaxAgeNumBlocks string `json:"max_age_num_blocks"`
			MaxAgeDuration  string `json:"max_age_duration"`
			MaxBytes        string `json:"max_bytes"`
		} `json:"evidence"`
		Validator struct {
			PubKeyTypes []string `json:"pub_key_types"`
		} `json:"validator"`
		Version struct {
		} `json:"version"`
		Abci struct {
		} `json:"abci"`
	} `json:"consensus_param_updates"`
	AppHash string `json:"app_hash"`
}

type ResultTX struct {
	Data      string          `json:"data,omitempty"`
	Log       string          `json:"log,omitempty"`
	GasWanted string          `json:"gas_wanted,omitempty"`
	GasUsed   string          `json:"gas_used,omitempty"`
	Events    []cmtypes.Event `json:"events,omitempty"`
}

type Event struct {
	Type       string       `json:"type,omitempty"`
	Attributes []*Attribute `json:"attributes,omitempty"`
}

type Attribute struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	Index bool   `json:"index,omitempty"`
}

/*
Decode the type:transfer to find the ammAccount:
Use recipient as ammAccount
*/
func (tx *Result) ammAccount() string {
	if tx.Data != nil &&
		tx.Data.Value != nil &&
		tx.Data.Value.TxResult.Result != nil &&
		tx.Data.Value.TxResult.Result.Events != nil {
		return ammAccountFromEvents(tx.Data.Value.TxResult.Result.Events)
	}
	return ""
}

func ammAccountFromEvents(events []cmtypes.Event) string {
	for _, event := range events {
		if event.Type == "transfer" {
			for _, attribute := range event.Attributes {
				if attribute.Key == "recipient" {
					return attribute.Value
				}
			}
		}
	}
	return ""
}

// Amount is in the format of int32string
func parseAmount(amount string) (string, string, decimal.Decimal) {
	re := regexp.MustCompile(`(\d+)(.*)`)

	matches := re.FindStringSubmatch(amount)
	if len(matches) != 3 {
		return "", "", decimal.Zero
	}
	amountInt, err := decimal.NewFromString(matches[1])
	if err != nil {
		return "", "", decimal.Zero
	}
	// The currency+issuer is index 2 in the array
	// Split the currency and issuer by -:
	currency, issuer := parseCurrency(matches[2])
	return currency, issuer, amountInt
}

func ParseAmount(amount string) (string, string, decimal.Decimal) {
	return parseAmount(amount)
}

func parsePair(pair string) (string, string, string, string) {
	// Split the pair into currencies:
	currencies := strings.Split(pair, ",")

	if len(currencies) == 2 {
		c, i := parseCurrency(currencies[0])
		c2, i2 := parseCurrency(currencies[1])
		return c, i, c2, i2
	}
	// If there are 3 currencies, then the first and third are the currencies, the second is the issuer
	// Determine which currency is the build in currency:
	buildIn := -1
	for i, currency := range currencies {
		if currency == "udevcore" || currency == "utestcore" || currency == "ucore" {
			// This is the build in currency:
			buildIn = i
			break
		}
	}
	// If the build in currency is the first, then the second is the currency with the issuer
	if buildIn == 0 {
		return currencies[0], "", currencies[1], currencies[2]
	}
	// If the build in currency is the second, then the first is the currency with the issuer
	return currencies[0], currencies[1], currencies[2], ""
}

// A currency consists out of:
// currency-issuer
// Or
// currency (in case of udevcore/utestcore/ucore)
// Or in the case of ibc (ibc/HASH) return the issuer only
func parseCurrency(currency string) (string, string) {
	// Split the currency into currency and issuer:
	ci := strings.Split(currency, "-")
	if len(ci) == 1 {
		// This is a build in currency or an IBC
		if strings.HasPrefix(currency, "ibc/") {
			// This is an IBC
			ci = strings.Split(currency, "/")
			return "", ci[1]
		}
		return currency, ""
	}
	return ci[0], ci[1]
}

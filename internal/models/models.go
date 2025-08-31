package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type GenerateRequest struct {
	Code   string         `json:"code"`
	Format string         `json:"format,omitempty"`
	Data   map[string]any `json:"data"`
}

func FlattenToStringMap(v any, prefix string) (map[string]string, error) {
	out := make(map[string]string)
	switch data := v.(type) {
	case map[string]any:
		for k, vv := range data {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			m, err := FlattenToStringMap(vv, key)
			if err != nil {
				return nil, err
			}
			for kk, vv := range m {
				out[kk] = vv
			}
		}
	case []any:
		b, _ := json.Marshal(data)
		out[prefix] = string(b)
	case string:
		out[prefix] = data
	case float64, int, int64, bool, float32:
		out[prefix] = fmt.Sprintf("%v", data)
	case nil:
		out[prefix] = ""
	default:
		return nil, errors.New("unsupported data type in FlattenToStringMap")
	}
	return out, nil
}

//// Transaction model
//type Transaction struct {
//	CreationTime    string
//	ProcessingTime  string
//	Description     string
//	OperationAmount string
//	AccountAmount   string
//	Commission      string
//}
//
//// Context for rendering the template
//type StatementContext struct {
//	Date              string
//	AccountNumber     string
//	AccountCurrency   string
//	CardName          string
//	ClientName        string
//	ClientTaxCode     string
//	StatementDateFrom string
//	StatementDateTo   string
//	CurrentTime       string
//	BlockedSum        string
//
//	InitialBalance string
//	Income         string
//	Expenses       string
//	FinalBalance   string
//
//	Transactions     []Transaction
//	WaitTransactions []Transaction
//}

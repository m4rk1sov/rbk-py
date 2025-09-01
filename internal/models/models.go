package models

import (
	"fmt"
	"reflect"
)

type GenerateRequest struct {
	Code   string         `json:"code"`
	Format string         `json:"format,omitempty"`
	Data   map[string]any `json:"data"`
}

//func FlattenToStringMap(v any, prefix string) (map[string]string, error) {
//	out := make(map[string]string)
//	switch data := v.(type) {
//	case map[string]any:
//		for k, vv := range data {
//			key := k
//			if prefix != "" {
//				key = prefix + "." + k
//			}
//			m, err := FlattenToStringMap(vv, key)
//			if err != nil {
//				return nil, err
//			}
//			for kk, vv := range m {
//				out[kk] = vv
//			}
//		}
//	case []any:
//		b, _ := json.Marshal(data)
//		out[prefix] = string(b)
//	case string:
//		out[prefix] = data
//	case float64, int, int64, bool, float32:
//		out[prefix] = fmt.Sprintf("%v", data)
//	case nil:
//		out[prefix] = ""
//	default:
//		return nil, errors.New("unsupported data type in FlattenToStringMap")
//	}
//	return out, nil
//}

// FlattenToStringMap converts nested structs/maps into a flat map[string]string
func FlattenToStringMap(data interface{}, prefix string) (map[string]string, error) {
	result := make(map[string]string)
	err := flattenHelper(data, prefix, result)
	if err != nil {
		return nil, err
	}

	// Expand keys so that we support {key}, {{key}}, and {{ key }}
	expanded := make(map[string]string, len(result)*3)
	for k, v := range result {
		expanded[k] = v
		expanded["{{"+k+"}}"] = v
		expanded["{{ "+k+" }}"] = v
	}

	return expanded, nil
}

func flattenHelper(data interface{}, prefix string, out map[string]string) error {
	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.Map:
		for _, key := range val.MapKeys() {
			strKey := fmt.Sprintf("%v", key.Interface())
			fullKey := strKey
			if prefix != "" {
				fullKey = prefix + "." + strKey
			}
			if err := flattenHelper(val.MapIndex(key).Interface(), fullKey, out); err != nil {
				return err
			}
		}
	case reflect.Struct:
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			fieldName := field.Name
			fullKey := fieldName
			if prefix != "" {
				fullKey = prefix + "." + fieldName
			}
			if err := flattenHelper(val.Field(i).Interface(), fullKey, out); err != nil {
				return err
			}
		}
	default:
		// Primitive value, just stringify
		out[prefix] = fmt.Sprintf("%v", val.Interface())
	}

	return nil
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

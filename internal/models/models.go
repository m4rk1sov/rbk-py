package models

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
		
	}
}

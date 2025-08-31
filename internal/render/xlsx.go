package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ExcelizeFill(template []byte, data map[string]any) ([]byte, error) {
	f, err := excelize.OpenReader(bytes.NewReader(template))
	if err != nil {
		return nil, err
	}
	defer func(f *excelize.File) {
		if closeErr := f.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}(f)
	
	// replace by placeholders
	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, err
		}
		for rowIdx, row := range rows {
			for cellIdx, cell := range row {
				if hasBraces(cell) {
					newVal := replacePlaceholders(cell, data)
					cellName, _ := excelize.CoordinatesToCellName(cellIdx+1, rowIdx+1)
					if err := f.SetCellStr(sheet, cellName, newVal); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	
	// table expansion on sheet name
	if tableAny, ok := data["card_statement"]; ok {
		if tableRows, ok := tableAny.([]any); ok && len(tableRows) > 0 {
			const sheet = "Card_statement"
			if idx, err := f.GetSheetIndex(sheet); idx > 0 && err != nil {
				//TODO adding the row choice mechanic
				srcRow := 3
				for i, rowDataAny := range tableRows {
					rowData, ok := rowDataAny.(map[string]any)
					if ok {
						dstRow := srcRow + i
						cols, err := f.Cols(sheet)
						if err != nil {
							return nil, err
						}
						// TODO choice for column
						colIdx := 2
						for cols.Next() {
							cellName, err := excelize.CoordinatesToCellName(colIdx, srcRow)
							if err != nil {
								return nil, err
							}
							val, err := f.GetCellValue(sheet, cellName)
							if err != nil {
								return nil, err
							}
							newVal := replacePlaceholders(val, rowData)
							dstCell, err := excelize.CoordinatesToCellName(colIdx, dstRow)
							if err != nil {
								return nil, err
							}
							err = f.SetCellStr(sheet, dstCell, newVal)
							if err != nil {
								return nil, err
							}
							colIdx++
						}
					} else {
						return nil, errors.New("no row data")
					}
				}
			} else {
				return nil, err
			}
		}
	}
	
	var out bytes.Buffer
	if err := f.Write(&out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func hasBraces(s string) bool {
	return len(s) >= 4 && (contains(s, "{{") && contains(s, "}}"))
}

func contains(s, sub string) bool {
	return bytes.Contains([]byte(s), []byte(sub))
}

func replacePlaceholders(input string, data map[string]any) string {
	out := input
	for k, v := range data {
		ph1 := "{{ " + k + " }}"
		ph2 := "{{" + k + "}}"
		out = bytesReplaceAll(out, ph1, toString(v))
		out = bytesReplaceAll(out, ph2, toString(v))
	}
	return out
}

func bytesReplaceAll(s, old, new string) string {
	return string(bytes.ReplaceAll([]byte(s), []byte(old), []byte(new)))
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

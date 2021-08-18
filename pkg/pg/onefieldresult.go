package pg

import (
	"fmt"
	"strings"
	"time"
)

func ResultValueToString(value interface{}) (s string, err error){
	switch v := value.(type) {
	case string:
		return v, nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case time.Duration:
		return v.String(), nil
	case time.Time:
		return v.String(), nil
	case int:
		return fmt.Sprintf("%d", v), nil
	case int8:
		return fmt.Sprintf("%d", v), nil
	case int16:
		return fmt.Sprintf("%d", v), nil
	case int32:
		return fmt.Sprintf("%d", v), nil
	case int64:
		return fmt.Sprintf("%d", v), nil
	case uint:
		return fmt.Sprintf("%d", v), nil
	case uint8:
		return fmt.Sprintf("%d", v), nil
	case uint16:
		return fmt.Sprintf("%d", v), nil
	case uint32:
		return fmt.Sprintf("%d", v), nil
	case uint64:
		return fmt.Sprintf("%d", v), nil
	case []byte:
		return fmt.Sprintf("%d", v), nil
	default:
		return "", fmt.Errorf("unhandled datatype %e", value)
	}
}

type OneFieldResult map[string]string
type OneFieldResults []OneFieldResult

func NewOneFieldResultFromByteArrayArray(cols []string, values []interface{}) (ofr OneFieldResult, err error) {
	ofr = make(OneFieldResult)
	if len(cols) != len(values) {
		return ofr, fmt.Errorf("number of cols different then number of values")
	}
	for i, col := range cols {
		value, err := ResultValueToString(values[i])
		if err != nil {
			return ofr, err
		}
		ofr[col] = value
	}
	return ofr, nil
}

func (ofr OneFieldResult) String() (s string) {
	var results []string
	for key, value := range ofr {
		key = strings.Replace(key, "'", "\\'", -1)
		value = strings.Replace(value, "'", "\\'", -1)
		results = append(results, fmt.Sprintf("'%s': '%s'", key, value))
	}
	return fmt.Sprintf("[ '%s' ]", strings.Join(results, "', '"))
}

func (ofr OneFieldResult) Columns() (result []string) {
	for key := range ofr {
		result = append(result, key)
	}
	return result
}
func (ofr OneFieldResult) Compare(other OneFieldResult) (err error) {
	if len(ofr) != len(other) {
		return fmt.Errorf("number of columns different between row ('%v') and compared row ('%v')", ofr.Columns(), other.Columns())
	}
	for key, value := range ofr {
		otherValue, exists := other[key]
		if ! exists {
			return fmt.Errorf("column row ('%s') not in compared row", key)
		}
		if value != otherValue {
			return fmt.Errorf("column '%s' differs between row ('%s'), and comparedrow ('%s')", key, value, otherValue)
		}
	}
	return nil
}

func (results OneFieldResults) String() (s string) {
	var arr []string
	for _, result := range results {
		arr = append(arr, result.String())
	}
	return fmt.Sprintf("[ { %s } ]", strings.Join(arr, "}, {"))
}

func (results OneFieldResults) Compare(other OneFieldResults) (err error) {
	if len(results) != len(other) {
		return fmt.Errorf("different result (%s) then expected (%s)", results.String(),
			other.String())
	}
	for i, result := range results {
		err = result.Compare(other[i])
		if err != nil {
			return fmt.Errorf("different %d'th result: %e", i, err)
		}
	}
	return nil
}


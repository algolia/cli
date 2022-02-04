package export

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strconv"

	"github.com/itchyny/gojq"
)

func FilterJSON(w io.Writer, input io.Reader, queryStr string) error {
	query, err := gojq.Parse(queryStr)
	if err != nil {
		return err
	}

	jsonData, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}

	var responseData interface{}
	err = json.Unmarshal(jsonData, &responseData)
	if err != nil {
		return err
	}

	iter := query.Run(responseData)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := v.(error); isErr {
			return err
		}
		if text, e := jsonScalarToString(v); e == nil {
			_, err := fmt.Fprintln(w, text)
			if err != nil {
				return err
			}
		} else {
			var jsonFragment []byte
			jsonFragment, err = json.Marshal(v)
			if err != nil {
				return err
			}
			_, err = w.Write(jsonFragment)
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(w, "\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func jsonScalarToString(input interface{}) (string, error) {
	switch tt := input.(type) {
	case string:
		return tt, nil
	case float64:
		if math.Trunc(tt) == tt {
			return strconv.FormatFloat(tt, 'f', 0, 64), nil
		} else {
			return strconv.FormatFloat(tt, 'f', 2, 64), nil
		}
	case nil:
		return "", nil
	case bool:
		return fmt.Sprintf("%v", tt), nil
	default:
		return "", fmt.Errorf("cannot convert type to string: %v", tt)
	}
}

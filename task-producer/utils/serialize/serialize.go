package serialize

import (
	"encoding/json"
)

func Serialize(v interface{}) (string, error) {
	//Convert job struct into json
	jsonString, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}

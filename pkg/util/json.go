package util

import (
	"encoding/json"
	"io"
	"log"
)

func JsonDecode(jsonByte []byte) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal(jsonByte, &data)

	return data, err
}

func JsonDecodeFile(file io.Reader) {
	decoder := json.NewDecoder(file)
	var data map[string]interface{}

	for {
		if err := decoder.Decode(&data); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	/*for key, value := range data{


	}*/

}

func setValues(data *interface{}, key *interface{}, value *interface{}) {

}

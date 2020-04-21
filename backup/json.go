package backup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func LoadSchemaTableMap(jsonFileName string) map[string][]string {
	jsonFile, err := os.Open(jsonFileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string][]string
	json.Unmarshal([]byte(byteValue), &result)
	var res map[string][]string = make(map[string][]string)
	if result != nil {
		for k, v := range result {
			for i := range v {
				v[i] = strings.ToUpper(v[i])
			}
			res[strings.ToUpper(k)] = v
		}
	}
	return res
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func logJson(data map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return nil, err
	}
	fmt.Println(string(jsonData))
	return jsonData, nil
}

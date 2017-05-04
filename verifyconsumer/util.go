package main

import (
	"encoding/json"
	_ "errors"
)

func response2Map(response []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	return result, nil
}

/*
func response2Map(response []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	if status, ok := result["status"]; ok {
		if status.(string) == "OK" {
			return result, nil
		} else {
			return nil, errors.New(result["status"].(string))
		}
	}

	return nil, errors.New("Failed")
}


func response2Map(response []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	if _, ok := result["error_message"]; ok {
		return nil, errors.New(result["error_message"].(string))
	}

	return result, nil
}*/

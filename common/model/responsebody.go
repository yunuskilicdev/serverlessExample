package model

import (
	"encoding/json"
)

type ResponseBody struct {
	Message        string
	ResponseObject interface{}
}

func (respBody ResponseBody) ConvertToJson() string {
	var jsonData []byte
	jsonData, responseJsonErr := json.Marshal(respBody)
	if responseJsonErr != nil {
		return ""
	}
	return string(jsonData)
}

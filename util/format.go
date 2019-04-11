package util

import (
	"delonakc.com/api/config"
	"encoding/json"
	"fmt"
)

// 响应数据格式
//
// code 状态码 0: success -1: error
type ResponseData struct {
	Code int `json:"code"`
	Err string `json:"err"`
	Data interface{} `json:"data"`
}

func ResponseSuccess(val interface{}) string {
	res := ResponseData{
		Code: config.Sucess,
		Err: "",
		Data: val,
	}

	value, _ := json.Marshal(&res)

	return fmt.Sprintf("%s", value)
}

func ResponseError(msg string) string {
	res := ResponseData{
		Code: config.CommonError,
		Err: msg,
		Data: nil,
	}

	value, _ := json.Marshal(&res)

	return fmt.Sprintf("%s", value)
}

func ResponseErrorWithCode(msg string, code int) string {
	res := ResponseData{
		Code: code,
		Err: msg,
		Data: nil,
	}

	value, _ := json.Marshal(&res)

	return fmt.Sprintf("%s", value)
}
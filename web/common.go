package main

import (
	"fmt"
	"encoding/json"
)

type ApiResponse struct {
	Ok bool
	Result interface{}
}

func apiError(ctx * web.Context, message string) {
	response := ApiResponse {
		Ok: false,
		Result: message,
	}
	ctx.Write(toJson(response));
}

func apiOk(result interface{}) (ApiResponse) {
	return ApiResponse {
		Ok: true,
		Result: result,
	}
}

func toJson(item interface{}) []byte {
	b, err := json.Marshal(item)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b;
}

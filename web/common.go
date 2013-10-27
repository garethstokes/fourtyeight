package main

import (
	"fmt"
  "os"
	"encoding/json"
  "github.com/hoisie/web"
  "github.com/garethstokes/fourtyeight/apns"
)

type ApiResponse struct {
	Ok bool `json:"ok"`
	Result interface{} `json:"result"`
}

func apiError(ctx * web.Context, message interface{}) {
	response := ApiResponse {
		Ok: false,
		Result: message,
	}
	ctx.Write(toJson(response));
}

func ok(ctx * web.Context, result interface{}) {
  ctx.Write(toJson(apiOk( result )))
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

func sendPushNotificationTo(token string, from string) {
  fmt.Println("sending push notification: ( " + token + " )")
  payload := apns.NewPayload()
  payload.Alert = from + " has just dropped a new message."

  pn := apns.NewPushNotification()
  pn.DeviceToken = token
  pn.AddPayload(payload)

  var wd, _ = os.Getwd()

  client := apns.NewClient(
    "gateway.sandbox.push.apple.com:2195",
    wd + "/keys/apns-dev-cert.pem",
    wd + "/keys/apns-dev-key-noenc.pem",
  )
  resp := client.Send(pn)

  alert, _ := pn.PayloadString()

  fmt.Println("  Alert:", alert)
  fmt.Println("Success:", resp.Success)
  fmt.Println("  Error:", resp.Error)
}

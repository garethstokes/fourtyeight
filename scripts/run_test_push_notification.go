package main

import (
  "fmt"
  "os"
  "crypto/tls"
  apns "github.com/garethstokes/fourtyeight/apns"
)

func main() {
  payload := apns.NewPayload()
  payload.Alert = "duh, computors"

  pn := apns.NewPushNotification()
  pn.DeviceToken = "b8172707ae9f253b07b722b846649f526e091200cdb8f07a2594eed4f8af16d6"
  pn.AddPayload(payload)

  var wd, _ = os.Getwd()

  cert, err := tls.LoadX509KeyPair(
    wd + "/apns-dev-cert.pem",
    wd + "/apns-dev-key-noenc.pem",
  )

  fmt.Println(cert)
  fmt.Println(err)

  if err != nil {
    fmt.Println("fucking errors man")
    return
  }

  client := apns.NewClient(
    "gateway.sandbox.push.apple.com:2195",
    wd + "/apns-dev-cert.pem",
    wd + "/apns-dev-key-noenc.pem",
  )
  resp := client.Send(pn)

  alert, _ := pn.PayloadString()
  fmt.Println("  Alert:", alert)
  fmt.Println("Success:", resp.Success)
  fmt.Println("  Error:", resp.Error)
}

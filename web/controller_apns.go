package main

import (
  "fmt"
	"encoding/json"
	"github.com/hoisie/web"
  "github.com/garethstokes/fourtyeight/cache"
  "github.com/garethstokes/fourtyeight/apns_android"
	"github.com/garethstokes/fourtyeight/personal"
)

func WarmApnCache(){
    // check if user is logged in
    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()
    
    p.FillCacheWithNotificationTokens()

}

func SendPushNotificationTo(users []string, message string){
    iosDeviceTokens := make([]string, len(users))
    androidDeviceTokens := make([]string, len(users))
 
    //gather the tokens for each user and each platform
    for _, user := range users{
      //ios
      iosToken := cache.Get("apns", user)
      if(iosToken!=nil){
        iosDeviceTokens = append(iosDeviceTokens, iosToken.(string))
      }
      //android
      androidToken := cache.Get("apns_android", user)
      if(androidToken!=nil){
        androidDeviceTokens = append(androidDeviceTokens, androidToken.(string))
      }
    }
 
    //ios
    if(len(iosDeviceTokens) > 0){
     // sendPushNotificationTo(deviceToken.(string), person.Username)
    }

    //android
    if(len(androidDeviceTokens) > 0){
      apns_android.SendNotification(0, message, androidDeviceTokens)
    }
}

func ApnsController() {

  // register a client for push notifications
  type ApnsRegisterParams struct {
    Token string
    DeviceToken string 
  }

  web.Post("/apns/register", func(ctx * web.Context) {
    ctx.SetHeader("Content-Type", "application/json", true)

    // get params
    params := new(ApnsRegisterParams)
    err := json.NewDecoder(ctx.Request.Body).Decode(&params)
    if err != nil {
      apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
      return
    }

    if params.DeviceToken == ""{
      apiError( ctx, "INVALID_NOTIFY_TOKEN" )
      fmt.Println("INVALID_NOTIFY_TOKEN")
      return
    }

    if params.Token == "" {
      apiError( ctx, "INVALID_TOKEN" )
      fmt.Println("INVALID_USER_TOKEN")
      return
    }

    // check if user is logged in
    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()
    
    user, error := p.GetLoggedInUser(params.Token)

    if error != nil {
      apiError( ctx, "INVALID_TOKEN" )
      fmt.Println("INVALID_USER_TOKEN")
      return
    }

    //woo hoo it worked, we found them

    //put it in the cache i guess
    cache.Set("apns", user.Username, params.DeviceToken)
    
    //put it in the database is better
    error = p.RegisteriOSDevice(user.Username, params.DeviceToken)

    if error != nil {
      apiError( ctx, "INVALID_TOKEN" )
      fmt.Println("INVALID_USER_TOKEN - failed to store the device token for apns")

      return
    }

    ok( ctx, params.DeviceToken )
  })


  web.Post("/apns/register/android", func(ctx * web.Context) {
    ctx.SetHeader("Content-Type", "application/json", true)

    // get params
    params := new(ApnsRegisterParams)
    err := json.NewDecoder(ctx.Request.Body).Decode(&params)
    if err != nil {
      apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
      return
    }

    // check if user is logged in
    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()
    
    user, error := p.GetLoggedInUser(params.Token)

    if error != nil {
      apiError( ctx, "INVALID_TOKEN" )
      fmt.Println("INVALID_USER_TOKEN")
      return
    }

    //woo hoo it worked, we found them

    //put it in the cache i guess
    cache.Set("apns_android", user.Username, params.DeviceToken)
    
    //put it in the database is better
    error = p.RegisterAndroidDevice(user.Username, params.DeviceToken)

    if error != nil {
      apiError( ctx, "INVALID_TOKEN" )
      fmt.Println("INVALID_USER_TOKEN - failed to store the device token for apns")

      return
    }

    ok( ctx, params.DeviceToken )
  })
}

package main

import (
  "fmt"
	"encoding/json"
	"github.com/hoisie/web"
  "github.com/garethstokes/fourtyeight/cache"
	"github.com/garethstokes/fourtyeight/personal"
)

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

    // check if user is logged in
    user := cache.Get("users", params.Token).(* personal.Person)
    if user == nil {
      apiError( ctx, "INVALID_TOKEN" )
      return
    }

    // set the cache
    //
    // we use the users username instead of token
    // so that this will perist over time
    //
    // TODO: store the array of items with some sort of timeout
    cache.Set("apns", user.Username, params.DeviceToken)

    ok( ctx, params.DeviceToken )
  })
}

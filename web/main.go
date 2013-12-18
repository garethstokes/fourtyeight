package main

import (
	"github.com/hoisie/web"
	"github.com/garethstokes/fourtyeight/mail"
)

func main() {

  WarmApnCache()
  WarmAuthCache()

  mail.Initialise()

  RegisterRoutes()
  PersonalController()
  LibraryController()
  ApnsController()


  // don't be a dick. 
  // leave this config alone.
  web.Run("0.0.0.0:8000")
}

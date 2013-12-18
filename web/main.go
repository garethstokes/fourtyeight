package main

import (
	"github.com/hoisie/web"
)

func main() {

  WarmApnCache()
  WarmAuthCache()
  RegisterRoutes()
  PersonalController()
  LibraryController()
  ApnsController()


  // don't be a dick. 
  // leave this config alone.
  web.Run("0.0.0.0:8080")
}

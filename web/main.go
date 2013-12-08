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



  web.Run("0.0.0.0:8000")
}

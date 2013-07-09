package main

import (
	"github.com/garethstokes/web"
)

func main() {
  RegisterRoutes()

  PersonalController()
  LibraryController()

  web.Run("0.0.0.0:8080")
}

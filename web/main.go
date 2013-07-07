package main

import (
	"github.com/garethstokes/web"
	"github.com/garethstokes/fourtyeight/personal"
)

func main() {
  RegisterRoutes();

  personal := personal.Store()
  personal.OpenSession()
  //defer personal.CloseSession()
  personal.InitialiseSchema()

  web.Run("0.0.0.0:8080")
}

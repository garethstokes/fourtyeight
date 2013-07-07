package main

import (
	"github.com/garethstokes/web"
)

func main() {
  RegisterRoutes();

  web.Run("0.0.0.0:8080")
}

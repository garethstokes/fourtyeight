package main

import (
	"github.com/garethstokes/fourtyeight/personal"
  	"fmt"
)

func main() {
  schema := "fourtyeight_development"
  fmt.Printf( "Dropping database :: %s\n", schema )

  p := personal.Store()
  p.OpenSession()
  defer p.CloseSession()

  p.InitialiseSchema()

}

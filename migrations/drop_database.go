package main

import (
	"github.com/garethstokes/fourtyeight/personal"
	"github.com/garethstokes/fourtyeight/library"
  "fmt"
)

func main() {
  schema := "fourtyeight_development"
  fmt.Printf( "Dropping database :: %s\n", schema )

  p := personal.Store()
  p.OpenSession()
  defer p.CloseSession()

  p.DropSchema()

  l := library.Store()
  l.OpenSession()
  l.DestroyCollectionAndCloseSession()
}

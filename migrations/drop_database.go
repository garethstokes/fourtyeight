package main

import (
	"github.com/garethstokes/fourtyeight/personal"
	"github.com/garethstokes/fourtyeight/library"
  "fmt"
)

func main() {
  schema := "fourtyeight_development"
  fmt.Printf( "Dropping database :: %s\n", schema )

  l := library.Store()
  l.OpenSession()
  l.DestroyCollectionAndCloseSession()

  p := personal.Store()
  p.OpenSession()
  p.DestroyCollectionAndCloseSession()
}

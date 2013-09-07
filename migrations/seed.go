package main

import (
	"github.com/garethstokes/fourtyeight/personal"
	"github.com/garethstokes/fourtyeight/library"
  "fmt"
)

func main() {
  schema := "fourtyeight_development"

  // PERSONAL
  p := personal.Store()
  p.OpenSession()
  defer p.CloseSession()

  p.Seed()

  // LIBRARY
  fmt.Printf( "Seeding library :: %s\n", schema )

  l := library.Store()
  l.Schema = schema

  l.OpenSession()
  defer l.CloseSession()

  l.CreateFrom(& library.Post { "@garrydanger", "http://i.imgur.com/FudYBky.jpg", "Took me a while to figure out that hand-situation." }, 60 * 60 * 48)
  l.CreateFrom(& library.Post { "@shredder", "", "guys, i think that i might need a shave." }, 60 * 60 * 48)
  l.CreateFrom(& library.Post { "@shredder", "", "guys, i think that i might need a shave." }, 1)
}

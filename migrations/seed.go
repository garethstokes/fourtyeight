package main

import (
	"github.com/garethstokes/fourtyeight/personal"
	"github.com/garethstokes/fourtyeight/library"
  "fmt"
  "time"
)

func main() {
  schema := "fourtyeight_development"
  fmt.Printf( "Seeding personal database :: %s\n", schema )

  p := personal.Store()
  p.OpenSession()
  defer p.CloseSession()

  garry := new( personal.Person )
  garry.Username = "@garrydanger"
  garry.Email = "garrydanger@gmail.com"
  garry.AvatarUrl = "https://si0.twimg.com/profile_images/2083020030/Photo_on_2012-03-16_at_15.47__2.jpg"

  _, error := p.Create( garry, "bobafett" )
  if error != nil {
    fmt.Printf( "Creating user, garry... ERROR\n" )
    fmt.Printf( "%s\n", error )
    return
  }

  fmt.Printf( "Creating user, garry... SUCCESS\n" )

  fmt.Printf( "Seeding library :: %s\n", schema )

  l := library.Store()
  l.Schema = schema

  l.OpenSession()
  defer l.CloseSession()

  l.CreateFrom(& library.Post { "@garrydanger", "http://i.imgur.com/FudYBky.jpg", "Took me a while to figure out that hand-situation.", time.Now().UTC() })
  l.CreateFrom(& library.Post { "@shredder", "", "guys, i think that i might need a shave.", time.Now().UTC() })
}

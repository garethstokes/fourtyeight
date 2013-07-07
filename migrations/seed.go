package main

import (
	"github.com/garethstokes/fourtyeight/personal"
  "fmt"
)

func main() {
  schema := "fourtyeight_development"
  fmt.Printf( "Seeding database :: %s\n", schema )

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

}

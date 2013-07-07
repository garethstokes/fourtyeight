package personal

import (
  "testing"
  "fmt"
)

func TestInsertAndFind(t * testing.T) {
  var store = new(Personal)
  store.Schema = "fourtyeight_test"

  store.OpenSession()
  defer store.CloseSession()

  store.InitialiseSchema()
  defer store.DropSchema()

  person := new( Person )
  person.Username = "@garrydanger"
  person.Email = "garrydanger@gmail.com"
  person.AvatarUrl = "https://si0.twimg.com/profile_images/2083020030/Photo_on_2012-03-16_at_15.47__2.jpg"

  password := "this is a password"

  _, error := store.Create( person, password )
  if error != nil {
    t.Fail()
  }

  person, error = store.FindByUsername( "@garrydanger" )
  if error != nil {
    t.Fail()
  }

  person, error = store.Validate( person.Username, password )
  if error != nil {
    fmt.Printf(
      "Incorrect Password ( %s, %s )\n",
      person.Username,
      password)

    t.Fail()
  }

  //fmt.Printf( "%@\n", person )
}

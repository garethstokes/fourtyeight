package personal

import (
  "testing"
  "fmt"
)

func TestInsertAndFind(t * testing.T) {
  fmt.Print( "\n\nTest Insert And Find\n" )

  var store = new(Personal)
  store.Schema = "fourtyeight_test"

  store.OpenSession()
  defer store.CloseSession()

  store.InitialiseSchema()
  defer store.DropSchema()

  person := garrydanger()
  password := "this is a password"

  _, error := store.Create( person, password )
  if error != nil {
    t.Fail()
  }

  person, error = store.FindByName( "@garrydanger" )
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

func TestGetFollowers(t * testing.T) {
  fmt.Print( "\n\nTest Followers \n" )

  var store = new( Personal )
  store.Schema = "fourtyeight_test"

  store.OpenSession()
  defer store.CloseSession()

  store.InitialiseSchema()
  defer store.DropSchema()

  // create the test user accounts
  store.Seed()

  garrydanger, _ := store.FindByName( "@garrydanger" )
  shredder, _ := store.FindByName( "@shredder" )

  followers, error := store.AddFollowerTo( garrydanger, shredder )
  if error != nil {
    fmt.Print( error.Error() )
    t.Fail()
    return
  }

  if len(followers) != 1 {
    fmt.Printf( "Incorrect number of followers, expecting 1 but found %d\n", len(followers) )
    t.Fail()
    return
  }

  if followers[0].Username != shredder.Username {
    fmt.Print( "Incorrect follower\n" )
    t.Fail()
  }
}

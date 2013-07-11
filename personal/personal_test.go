package personal

import (
  "testing"
  "fmt"
)

func BeginTest( name string ) (store * Personal) {
  fmt.Printf( "\n%s\n", name )

  store = new(Personal)
  store.Schema = "fourtyeight_test"

  store.OpenSession()

  store.InitialiseSchema()

  store.Seed()
  return store
}

func CleanUp(store * Personal) {
  store.DropSchema()
  store.CloseSession()
}

func TestFindByName(t * testing.T) {
  s := BeginTest( "Find by name" )
  defer CleanUp(s)

  person, error := s.FindByName( "@garrydanger" )
  if error != nil {
    t.Fail()
  }

  if person.Username != "@garrydanger" {
    fmt.Print( "Incorrect user returned." )
    t.Fail()
  }
}

func TestValidate(t * testing.T) {
  s := BeginTest( "Validate user" )
  defer CleanUp(s)

  username := "@garrydanger"
  password := "bobafett"

  person, error := s.Validate( username, password )
  if error != nil {
    fmt.Printf(
      "Incorrect Password ( %s, %s )\n",
      username,
      password)

    t.Fail()
  }

  if person.Username != username {
    fmt.Print( "wrong username returned." )
    t.Fail()
  }
}

func TestGetFollowers(t * testing.T) {
  s := BeginTest( "Get Followers" )
  defer CleanUp(s)

  garrydanger, _ := s.FindByName( "@garrydanger" )
  shredder, _ := s.FindByName( "@shredder" )

  followers, _ := s.FollowersFor( garrydanger )

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

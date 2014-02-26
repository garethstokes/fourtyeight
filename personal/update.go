package personal

import (
  "labix.org/v2/mgo/bson"
  "github.com/garethstokes/fourtyeight/cache"
  "fmt"
)

func (s * Personal) Update( person * Person ) error {
  fmt.Printf( "Personal.Update :: %@\n", person )

  err := s.collection.Update(bson.M{"key": person.Key}, person)

  cache.Set("users", person.LoginToken, person)

  return err;
}


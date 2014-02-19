package personal

import (
  "github.com/garethstokes/fourtyeight/passwords"
  "time"
  "labix.org/v2/mgo/bson"
  "fmt"
)

func (s * Personal) Create( person * Person, password string ) (* Person, error) {
 
  fmt.Printf( "Personal.Create :: %@\n", person )
  
  key := bson.NewObjectId()

  person.Key = key

  pass := passwords.Compute( password )
  
  person.Password = pass.Hash
  person.Salt = pass.Salt
  person.Iterations = pass.Iterations

  person.DateCreated = time.Now().UTC().Unix()

  err := s.collection.Insert(person)
  if err != nil {
    panic(err)
  }

  return person, nil
}

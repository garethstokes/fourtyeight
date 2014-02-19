package personal

import (
  "labix.org/v2/mgo/bson"
  "fmt"
)

func (s * Personal) Update( person * Person ) error {
  fmt.Printf( "Personal.Update :: %@\n", person )

  err := s.collection.Update(bson.M{"key": person.Key}, person)
  return err;
}


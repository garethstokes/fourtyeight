package personal

import (
  "github.com/garethstokes/fourtyeight/passwords"
  "time"
  "fmt"
)

type PersonValidateError struct {
  When time.Time
  What string
}

func (e PersonValidateError) Error() string {
  return fmt.Sprintf("%v: %v", e.When, e.What)
}

func (s * Personal) Validate( username string, password string ) (* Person, error) {
  s.logf( "Personal.Validate( '%s', '%s' )", username, password )

  person, error := s.FindByName( username )
  if error != nil {
    return nil, error
  }

  auth := passwords.ComputeWithSalt( password, person.Iterations, person.Salt )
  if auth.Hash != person.Password {
    e := PersonValidateError {
		  time.Now(),
		  "Incorrect password",
	  }
    return nil, e
  }

  return person, nil
}

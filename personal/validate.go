package personal

import (
  "github.com/garethstokes/fourtyeight/passwords"
  "time"
  "fmt"
)

var (
  validate_db_columns = "password, salt, iterations"
  validate_db_predicate = fmt.Sprintf( "user_id = ? limit 1" )
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
  var iterations int
  var hash, salt string

  person, error := s.FindByUsername( username )
  if error != nil {
    return nil, error
  }

  sql := fmt.Sprintf(
    "SELECT %s FROM user WHERE %s",
    validate_db_columns,
    validate_db_predicate)

  row := s.db.QueryRow( sql, person.id )
  error = row.Scan( &hash, &salt, &iterations )

  if error != nil {
    s.error( error.Error() )
    return nil, error
  }

  auth := passwords.ComputeWithSalt( password, iterations, salt )
  if auth.Hash != hash {
    e := PersonValidateError {
		  time.Now(),
		  "Incorrect password",
	  }
    return nil, e
  }

  return person, nil
}

package personal

import (
  "github.com/garethstokes/fourtyeight/passwords"
  "time"
)

func (s * Personal) Create( person * Person, password string ) (* Person, error) {
  sql := "INSERT INTO user ( username, email, avatar_url, password, salt, iterations, date_created ) VALUES ( ?, ?, ?, ?, ?, ?, ? );"
  statement, error := s.db.Prepare( sql )
  if error != nil {
    panic( error )
  }
  defer statement.Close()

  s.log(sql)

  pass := passwords.Compute( password )

  auth := new( PersonAuthorisation )
  auth.Password = pass.Hash
  auth.Salt = pass.Salt
  auth.Iterations = pass.Iterations

  _, error = statement.Exec(
    person.Username,
    person.Email,
    person.AvatarUrl,
    auth.Password,
    auth.Salt,
    auth.Iterations, 
    time.Now())

  if error != nil {
    s.error( error.Error() )
    return nil, error
  }

  return person, nil
}

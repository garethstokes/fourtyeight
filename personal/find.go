package personal

import (
  "fmt"
)

var (
  db_columns = "user_id, username, email, avatar_url"
  db_predicate = fmt.Sprintf( "username = ? limit 1" )
)

func (s * Personal) FindByUsername( name string ) (* Person, error) {
  var user_id int
  var username, email, avatar_url string

  sql := fmt.Sprintf(
    "SELECT %s FROM user WHERE %s",
    db_columns,
    db_predicate)

  row := s.db.QueryRow( sql, name )
  error := row.Scan( &user_id, &username, &email, &avatar_url )

  if error != nil {
    s.error( error.Error() )
    return nil, error
  }

  person := new( Person )
  person.id = user_id
  person.Username = username
  person.Email = email
  person.AvatarUrl = avatar_url

  return person, nil
}

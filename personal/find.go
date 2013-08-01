package personal

import (
  "fmt"
  "database/sql"
)

var (
  db_columns = "user_id, username, email, avatar_url, date_created"
  db_predicate = fmt.Sprintf( "username = ? limit 1" )
)

func (s * Personal) FindByName( name string ) (* Person, error) {
  return s.txFindByName( name, nil )
}

func (s * Personal) txFindByName( name string, tx * sql.Tx ) (* Person, error) {
  var user_id int
  var username, email, avatar_url string
  var date_created int64
  var row * sql.Row

  s.logf( "personal.FindByName :: %s", name )

  sql := fmt.Sprintf(
    "SELECT %s FROM user WHERE %s",
    db_columns,
    db_predicate)

  if tx == nil {
    row = s.db.QueryRow( sql, name )
  } else {
    row = tx.QueryRow( sql, name )
  }

  error := row.Scan( &user_id, &username, &email, &avatar_url, &date_created )

  if error != nil {
    s.error( error.Error() )
    return nil, error
  }

  person := new( Person )
  person.id = user_id
  person.Username = username
  person.Email = email
  person.AvatarUrl = avatar_url
  person.DateCreated = date_created

  return person, nil
}

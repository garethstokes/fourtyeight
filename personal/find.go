package personal

import (
  "fmt"
  "errors"
  "database/sql"
  "github.com/garethstokes/fourtyeight/cache"
)

var (
  db_columns = "user_id, username, email, avatar_url, loginToken, date_created"
  db_predicate_token = fmt.Sprintf( "loginToken = ? limit 1" )
  db_predicate = fmt.Sprintf( "username = ? limit 1" )
)


func (s * Personal) FillCacheWithLoginTokens(){

  s.logf( "personal.FillCacheWithLoginTokens" )

  sql := fmt.Sprintf(
    "SELECT %s FROM user WHERE loginToken <> '' ",
    db_columns)

  rows, err := s.db.Query(sql)

  if err != nil {
    s.log(err.Error())
  }

  var count int
  count = 0 
  
  var user_id int
  var username, email, avatar_url, loginToken string
  var date_created int64

  for rows.Next() {
    count++

    error := rows.Scan( &user_id, &username, &email, &avatar_url, &loginToken, &date_created )

    if error != nil {
      s.error( error.Error() )
      return 
    }

    person := new( Person )
    person.id = user_id
    person.Username = username
    person.Email = email
    person.AvatarUrl = avatar_url
    person.DateCreated = date_created
    person.LoginToken = loginToken
    
    cache.Set("users", loginToken, person) 
    
  }
  
  fmt.Printf("Populated cache with %d login tokens \n", count)
  
  if err := rows.Err(); err != nil {
    s.log(err.Error())
  }

}



 


func (s * Personal) FindByToken(token string) (* Person, error) {
  
  s.logf( "personal.FindByToken :: %s", token )

  sql := fmt.Sprintf(
    "SELECT %s FROM user WHERE %s",
    db_columns,
    db_predicate_token)

  return s.txFindBySql(sql, token, nil)
}





func (s * Personal) FindByName( name string ) (* Person, error) {
  return s.txFindByName( name, nil )
}

func (s * Personal) txFindByName( name string, tx * sql.Tx ) (* Person, error) {

  s.logf( "personal.FindByName :: %s", name )

  sql := fmt.Sprintf(
    "SELECT %s FROM user WHERE %s",
    db_columns,
    db_predicate)

  return s.txFindBySql(sql, name, tx)

}





func (s * Personal) txFindBySql( sqlStr string, sqlParam string, tx * sql.Tx ) (* Person, error) {
  var user_id int
  var username, email, avatar_url, loginToken string
  var date_created int64
  var row * sql.Row

  if tx == nil {
    row = s.db.QueryRow( sqlStr, sqlParam )
  } else {
    row = tx.QueryRow( sqlStr, sqlParam )
  }

  error := row.Scan( &user_id, &username, &email, &avatar_url, &loginToken, &date_created )

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
  person.LoginToken = loginToken

  return person, nil
}


func (s * Personal) GetLoggedInUser(loginToken string)(* Person, error) {
    
    // check if user is logged in to the cache

    var user * Person
  
    person := cache.Get("users", loginToken)
    if person == nil {

      //not in cache lets try the DB
      //GArry danger might have restarted the server
      //LOlcats!

      user, error := s.FindByToken(loginToken)
      
      if error != nil {
        s.error( error.Error() )
        return nil, error
      }
      
      if user == nil {
        err := errors.New("User with token not found")
        return nil, err
      }
    }else{

      user = person.(* Person)

    }

    return user, nil
}

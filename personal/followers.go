package personal

import (
  "fmt"
  "database/sql"
)

var (
  followers_db_columns = "u.user_id, u.username, u.email, u.avatar_url, u.date_created"
  followers_db_predicate = fmt.Sprintf( "f.user_id = ?" )
)

func (s * Personal) FollowersFor( user * Person ) ([]Person, error) {
  return s.txFollowersFor( user, nil )
}


func (s * Personal) txFollowersFor( user * Person, tx * sql.Tx ) ([]Person, error) {
  var user_id int
  var username, email, avatar_url string
  var date_created int64
  var rows * sql.Rows
  var error error

  s.logf( "personal.FollowersFor :: %s", user.Username )

  sql := fmt.Sprintf(
    "SELECT %s FROM follower f INNER JOIN user u on u.user_id = f.follower_id WHERE %s",
    followers_db_columns,
    followers_db_predicate)

  if tx == nil {
    rows, error = s.db.Query( sql, user.id )
  } else {
    rows, error = tx.Query( sql, user.id )
  }

  s.logf( "SQL: %s\n", sql )

  if error != nil {
    s.error( error.Error() )
    return nil, error
  }

  followers := make( []Person, 0 )

  for rows.Next() {
    if error = rows.Scan( &user_id, &username, &email, &avatar_url, &date_created ); error != nil {
      s.error( error.Error() )
      return nil, error
    }

    person := new( Person )
    person.id = user_id
    person.Username = username
    person.Email = email
    person.AvatarUrl = avatar_url
    person.DateCreated = date_created

    s.logf( "adding %s to result", person.Username )
    followers = append( followers, * person )
  }

  return followers, nil
}

func (s * Personal) Following( user * Person ) ([]Person, error) {
  return s.txFollowing( user, nil )
}



func (s * Personal) txFollowing( user * Person, tx * sql.Tx ) ([]Person, error) {
  var user_id int
  var username, email, avatar_url string
  var date_created int64
  var rows * sql.Rows
  var error error

  s.logf( "personal.Following :: %s", user.Username )

  sql := fmt.Sprintf(
    "SELECT %s FROM user u INNER JOIN follower f ON f.user_id = u.user_id WHERE f.follower_id = ?;",
    followers_db_columns)

  if tx == nil {
    rows, error = s.db.Query( sql, user.id )
  } else {
    rows, error = tx.Query( sql, user.id )
  }

  s.logf( "SQL: %s\n", sql )

  if error != nil {
    s.error( error.Error() )
    return nil, error
  }

  followers := make( []Person, 0 )

  for rows.Next() {
    if error = rows.Scan( &user_id, &username, &email, &avatar_url, &date_created ); error != nil {
      s.error( error.Error() )
      return nil, error
    }

    person := new( Person )
    person.id = user_id
    person.Username = username
    person.Email = email
    person.AvatarUrl = avatar_url
    person.DateCreated = date_created

    s.logf( "adding %s to result", person.Username )
    followers = append( followers, * person )
  }

  return followers, nil
}

func (s * Personal) RemoveFollowerFrom( user * Person, follower * Person ) ([]Person, error) {
   s.logf( "persona.RemoveFollowerFrom :: ( %s, %s )\n", user.Username, follower.Username )

  tx, error := s.db.Begin()
  if error != nil {
    s.error( error.Error() )
    return nil, error
  }

  var followers []Person
  followers, error = s.txFollowersFor( user, tx )
  if error != nil {
    s.error( error.Error() )
    return nil, error
  }

  // look up the user if we don't have any ids
  if user.id == 0 {
    user, _ = s.txFindByName( user.Username, tx )
  }

  // look up the follower if we don't have any ids
  if follower.id == 0 {
    follower, _ = s.txFindByName( follower.Username, tx )
  }

  // if we reach here then it is safe to add
  // the user to the list of followers.

  sql := "DELETE FROM follower WHERE user_id =? AND follower_id=?;"
  
  statement, error := tx.Prepare( sql )
  if error != nil {
    s.error( error.Error() )
    return nil, error
  }
  defer statement.Close()

  s.log(sql)

  _, error = statement.Exec(user.id, follower.id)

  if error != nil {
    s.error( error.Error() )
    return nil, error
  }
 
  followers, _ = s.txFollowersFor( user, tx )
  if error == nil {
    err := tx.Commit()
    s.logf( "tx_result: %@\n", err )
  } else {
    s.error( error.Error() )
    tx.Rollback()
  }

  return followers, nil
}

func (s * Personal) AddFollowerTo( user * Person, follower * Person ) ([]Person, error) {
  s.logf( "persona.AddFollowerTo :: ( %s, %s )\n", user.Username, follower.Username )

  tx, error := s.db.Begin()
  if error != nil {
    s.error( error.Error() )
    return nil, error
  }

  var followers []Person
  followers, error = s.txFollowersFor( user, tx )
  if error != nil {
    s.error( error.Error() )
    return nil, error
  }

  // look for the follower in the list and return
  // if already there. 
  for i := 0; i < len(followers); i++ {
    f := followers[i]
    if f.Username == follower.Username {
      s.logf( "%s is already following %s", follower.Username, user.Username )
      tx.Rollback()
      return followers, nil
    }
  }

  // look up the user if we don't have any ids
  if user.id == 0 {
    user, _ = s.txFindByName( user.Username, tx )
  }

  // look up the follower if we don't have any ids
  if follower.id == 0 {
    follower, _ = s.txFindByName( follower.Username, tx )
  }

  // if we reach here then it is safe to add
  // the user to the list of followers.
  sql := fmt.Sprintf( "INSERT INTO follower (user_id, follower_id) VALUES ( %d, %d );",
    user.id,
    follower.id)

  s.log( sql )

  _, error = tx.Exec(sql)
  followers, _ = s.txFollowersFor( user, tx )
  if error == nil {
    err := tx.Commit()
    s.logf( "tx_result: %@\n", err )
  } else {
    s.error( error.Error() )
    tx.Rollback()
  }

  return followers, nil
}

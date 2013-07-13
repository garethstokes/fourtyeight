package personal

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "time"
import "fmt"

// the backticks are field tags that let
// the json parser know what to name the fields
type Person struct {
  id int
  Username string `json:"name"`
  Email string `json:"email"`
  AvatarUrl string `json:"avatarUrl"`
  DateCreated time.Time `json:"dateCreated"`
}

type PersonAuthorisation struct {
  Password string
  Salt string
  Iterations int
}

type Personal struct {
  db * sql.DB
  Schema string
}

func Store() * Personal {
  store := new( Personal )
  return store
}

func (s * Personal) OpenSession() {
  //db, err := sql.Open("mysql", "user:password@/dbname")
  if len( s.Schema ) == 0 {
    s.Schema = "fourtyeight_development"
  }

  connectionString := fmt.Sprintf( "root@/%s?parseTime=true", s.Schema )
  db, err := sql.Open("mysql", connectionString )
  if err != nil {
    s.error( err.Error() )
  }

  s.db = db

  s.logf( "Personal Schema :: %s", s.Schema )
}

func (s * Personal) CloseSession() {
  s.log("Personal.CloseSession")
  s.db.Close()
}

func (s * Personal) run(sql string) {
  statement, error := s.db.Prepare( sql )
  if error != nil {
    panic( error )
  }
  defer statement.Close()

  _, error = statement.Exec()
  if error != nil {
    s.error( error.Error() )
  }
}

func (s * Personal) log(message string) {
  fmt.Printf("%s\n", message)
}

func (s * Personal) logf(message string, args ...interface{}) {
  text := fmt.Sprintf( message, args... )
  fmt.Printf( "%s\n", text )
}

func (s * Personal) error(message string) {
  fmt.Printf("ERROR: %s\n", message)
}

func (s * Personal) createTable(name string, sql string) {
  s.logf( "Initialising Schema :: creating table %s", name )
  s.run( sql )
  s.run( fmt.Sprintf("INSERT INTO db_schema ( name, date_created ) VALUES ( '%s', NOW() );", name ))
}

func (s * Personal) InitialiseSchema() {
  sql := "SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = 'schema';"

  result, err := s.db.Exec( sql, s.Schema )
  if err != nil {
    panic( err.Error() )
  }

  rows, _ := result.RowsAffected()
  fmt.Printf( "number of rows: %d\n", rows )

  if rows == 0 {
    s.logf( "Initialising Schema :: %s", s.Schema )
    s.run( "CREATE TABLE db_schema ( table_id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255), date_created TIMESTAMP );" )

    var sql string
    sql = "CREATE TABLE user ( user_id INT(11) PRIMARY KEY AUTO_INCREMENT, username VARCHAR(255) NOT NULL UNIQUE, email VARCHAR(255) NOT NULL UNIQUE, avatar_url VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL, salt VARCHAR(255) NOT NULL, iterations INT, date_created DATETIME NOT NULL);"
    s.createTable( "user", sql );

    sql = "CREATE TABLE follower (user_id INT(11) PRIMARY KEY AUTO_INCREMENT, follower_id INT(11) NOT NULL, CONSTRAINT `fk_follower_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`), CONSTRAINT `fk_follower_follower_id` FOREIGN KEY (`follower_id`) REFERENCES `user` (`user_id`), UNIQUE KEY `unique_followers` (`user_id`, `follower_id`));"
    s.createTable( "follower", sql )

    sql = "CREATE TABLE waitinglist ( id INT(11) PRIMARY KEY AUTO_INCREMENT, email VARCHAR(255) NOT NULL, date_created TIMESTAMP );"
    s.createTable( "waitinglist", sql )
  }
}

func (s * Personal) DropSchema() {
  s.log( "Personal.DropSchema" )

  s.run( fmt.Sprintf("DROP DATABASE %s;", s.Schema) );
  s.run( fmt.Sprintf("CREATE DATABASE %s;", s.Schema) );
}

package personal

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "time"
import "fmt"

type Person struct {
  id int
  Username string
  Email string
  AvatarUrl string
  DateCreated time.Time
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

func (s * Personal) OpenSession() {
  //db, err := sql.Open("mysql", "user:password@/dbname")
  if len( s.Schema ) == 0 {
    s.Schema = "fourtyeight_development"
  }

  connectionString := fmt.Sprintf( "root@/%s", s.Schema )
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

func (s * Personal) InitialiseSchema() {
  s.run("CREATE TABLE user ( user_id INT PRIMARY KEY , username VARCHAR(255), email VARCHAR(255), avatar_url VARCHAR(255), password VARCHAR(255), salt VARCHAR(255), iterations INT);")
}

func (s * Personal) DropSchema() {
  s.log( "Personal.DropSchema" )

  s.run( fmt.Sprintf("DROP DATABASE %s;", s.Schema) );
  s.run( fmt.Sprintf("CREATE DATABASE %s;", s.Schema) );
}

package personal

import "fmt"
import "labix.org/v2/mgo"

// the backticks are field tags that let
// the json parser know what to name the fields

type NotificationToken struct {
  Platform int
  Token string
}

type Person struct {
  PersonAuthorisation
  Key interface{} `json:"key"`
  Username string `json:"name"`
  Email string `json:"email"`
  AvatarUrl string `json:"avatarUrl"`
  LoginToken string `json:"-"`
  DateCreated int64 `json:"dateCreated"`
  Followers []string `json:"followers"`
  Following []string `json:"following"`
  NotificationTokens []NotificationToken `json:"-"`
}

func (p * Person) Validate() []string {
  err := make([]string, 0)
  if len(p.Username) == 1 {
    err = append(err, "username: missing")
  }

  if len(p.Email) == 0 {
    err = append(err, "email: missing")
  }

  return err
}

type PersonAuthorisation struct {
  Password string `json:"-"`
  Salt string `json:"-"`
  Iterations int `json:"-"`
}

type Personal struct {
  session * mgo.Session
  collection * mgo.Collection
  Schema string
}

func Store() * Personal {
  store := new( Personal )
  return store
}

func (s * Personal) OpenSession() {
  fmt.Print( "Personal.OpenSession\n" )

    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    if len( s.Schema ) == 0 {
      s.Schema = "fourtyeight_development"
    }

    s.session = session
    s.collection = session.DB(s.Schema).C("personal")
}

func (s * Personal) CloseSession() {
  s.log("Personal.CloseSession")
  s.session.Close()
}


func (s * Personal) DestroyCollectionAndCloseSession() {
  fmt.Print( "Personal.DestroyCollectionAndCloseSession\n" )
  s.collection.DropCollection()
  s.session.Close()
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

/*
func (s * Personal) SaveToWaitingList(email string) error {
  sql := "INSERT INTO waitinglist ( email, date_created ) VALUES ( ?, NOW() );"
  statement, error := s.db.Prepare( sql )
  if error != nil {
    s.error( error.Error() )
    return error
  }
  defer statement.Close()

  s.log(sql)

  _, error = statement.Exec( email )

  if error != nil {
    s.error( error.Error() )
    return error
  }

  return nil
}
 
*/

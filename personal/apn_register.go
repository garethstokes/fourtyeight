package personal

// import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "fmt"
import "github.com/garethstokes/fourtyeight/cache"

func (s * Personal) FindAndroidTokens(username string) ([]string) {
  return s.FindTokens(username, 0)
}

func (s * Personal) FindiOSTokens(username string)  ([]string){
 return s.FindTokens(username, 1) 
}

func (s * Personal) FillCacheWithNotificationTokens() {

  rows, err := s.db.Query("SELECT username, token, deviceType FROM pushNotificationRegister")

  if err != nil {
    s.log(err.Error())
  }
  
  var deviceType, count int
  count = 0
  var token, username string
  
  for rows.Next() {
    count++

    if err := rows.Scan(&username, &token, &deviceType); err != nil {
      s.log(err.Error())
    }
    if(deviceType == 1){
      cache.Set("apns", username, token) 
    }else{
      cache.Set("apns_android", username, token)
    }
  }
  
  fmt.Printf("Populated cache with %d deviceTokens \n", count)
  
  if err := rows.Err(); err != nil {
    s.log(err.Error())
  }

  return 
}

func (s * Personal) FindTokens(username string, deviceType int) ([]string) {
  
  tokens := make([]string, 2)
  
  rows, err := s.db.Query("SELECT token FROM pushNotificationRegister WHERE username=? AND deviceType=?", username, deviceType)

  if err != nil {
    s.log(err.Error())
  }

  for rows.Next() {
    var token string
    if err := rows.Scan(&token); err != nil {
      s.log(err.Error())
    }

    tokens = append( tokens, token )
   
    fmt.Printf("%s token found for %d\n", token, username)
  }

  if err := rows.Err(); err != nil {
    s.log(err.Error())
  }
  
  return tokens
}

func (s * Personal) RegisterAndroidDevice(username string, token string) error {
  return s.RegisterDevice(username, token, 0)
}

func (s * Personal) RegisteriOSDevice(username string, token string) error {
  return s.RegisterDevice(username, token, 1)
}

func (s * Personal) RegisterDevice(username string, token string, deviceType int) error {
  sql := "INSERT INTO pushNotificationRegister ( username, token, deviceType, date_created ) VALUES ( ?, ?, ?, NOW() );"
  statement, error := s.db.Prepare( sql )
  if error != nil {
    s.error( error.Error() )
    return error
  }
  defer statement.Close()

  s.log(sql)

  _, error = statement.Exec( username, token, deviceType )

  if error != nil {
    s.error( error.Error() )
    return error
  }

  return nil
}
package personal

import (
  "fmt"
  "labix.org/v2/mgo/bson"
  "github.com/garethstokes/fourtyeight/cache"
)

var (
  ANDROID   = 0
  IOS       = 1
)

func (s * Personal) FillCacheWithNotificationTokens() {
  var result = make([]Person, 0)
  var count = 0

  s.collection.Find(bson.M{}).All(&result)

  for _, p := range result {
    for _, apn := range p.NotificationTokens {
      if(apn.Platform == IOS){
        cache.Set("apns", p.Username, apn.Token)
      }else{
        cache.Set("apns_android", p.Username, apn.Token)
      }
      count++
    }
  }

  fmt.Printf("Populated cache with %d deviceTokens \n", count)
}

// func (s * Personal) RegisterAndroidDevice(username string, token string) error {
//   return s.RegisterDevice(username, token, 0)
// }

// func (s * Personal) RegisteriOSDevice(username string, token string) error {
//   return s.RegisterDevice(username, token, 1)
// }

// func (s * Personal) RegisterDevice(username string, token string, deviceType int) error {

//   s.UnRegisterDevice(token)

//   sql := "INSERT INTO pushNotificationRegister ( username, token, deviceType, date_created ) VALUES ( ?, ?, ?, NOW() );"
  
//   statement, error := s.db.Prepare( sql )
//   if error != nil {
//     s.error( error.Error() )
//     return error
//   }
//   defer statement.Close()

//   s.log(sql)

//   _, error = statement.Exec( username, token, deviceType )

//   if error != nil {
//     s.error( error.Error() )
//     return error
//   }

//   return nil
// }

// func (s * Personal) UnRegisterDevice(token string) error{
//   //remove from cache first
//   //WARNING assuming that there is never going to be an android token that is the same as the ios token... safe assumption right? LIke the song...
//   cache.Remove("apns_android", token)
//   cache.Remove("apns", token)

//   sql := "DELETE FROM pushNotificationRegister where token =?;"
  
//   statement, error := s.db.Prepare( sql )
//   if error != nil {
//     s.error( error.Error() )
//     return error
//   }
//   defer statement.Close()

//   s.log(sql)

//   _, error = statement.Exec(token)

//   if error != nil {
//     s.error( error.Error() )
//     return error
//   }

//   return nil
// }

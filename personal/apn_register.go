package personal

import (
  "fmt"
  "github.com/garethstokes/fourtyeight/cache"
)

var (
  ANDROID   = 0
  IOS       = 1
)

func (s * Personal) FillCacheWithNotificationTokens() {
  var count = 0

  result := s.FindAll()

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
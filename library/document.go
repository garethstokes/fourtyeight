package library

import (
  "time"
  "fmt"
)

type Document struct {
  Key interface{} `json:"key"`
  MainPost * Post `json:"mainPost"`
  Comments []Post `json:"comments"`
  DateCreated int64 `json:"dateCreated"`
  ExpirationDelta int64 `json:"expirationDelta"`
}

func (d * Document) expired() bool {
  now := time.Now().UTC().Unix()
  delta := d.ExpirationDelta / 1000000000
  expiry := d.DateCreated + int64(delta)

  fmt.Printf("now: %d, expiry: %d, created: %d\n", now, expiry, d.DateCreated)
  return now > expiry
}

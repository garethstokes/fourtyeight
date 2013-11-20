package library

import (
  "time"
  "fmt"
)

/*
  TODO: create a lastUpdated field and
        update it on:
            - adding a comment
            - liking a post
 */
type Document struct {
  Key interface{} `json:"key"`
  MainPost * Post `json:"mainPost"`
  Comments []Post `json:"comments"`
  DateCreated int64 `json:"dateCreated"`
  ExpirationDelta int64 `json:"expirationDelta"`
  LastModified int64 `json:"lastModified"`
}

func (d * Document) expired() bool {
  now := time.Now().UTC().Unix()
  delta := d.ExpirationDelta
  expiry := d.DateCreated + int64(delta)

  fmt.Printf("now: %d, expiry: %d, created: %d\n", now, expiry, d.DateCreated)
  return now > expiry
}

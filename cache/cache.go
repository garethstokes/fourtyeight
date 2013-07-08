package cache

import "fmt"

var db = make( map[string] interface{} )

func Get(key string) interface{} {
  fmt.Printf( "cache.Get :: %s\n", key )
  return db[key]
}

func Set(key string, value interface{}) {
  fmt.Printf( "cache.Set :: %s\n", key )
  db[key] = value
}

func SetWithTimeout(key string, value interface{}, timeout int64) {
  db[key] = value
}

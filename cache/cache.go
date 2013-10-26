package cache

import "fmt"

var db = make( map[string] interface{} )

func Get(namespace string, key string) interface{} {
  key = namespace + "." + key
  fmt.Printf( "cache.Get :: %s\n", key )
  return db[key]
}

func Set(namespace string, key string, value interface{}) {
  key = namespace + "." + key
  fmt.Printf( "cache.Set :: %s\n", key )
  db[key] = value
}

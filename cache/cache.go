package cache

var db = make( map[string] interface{} )

func Get(key string) interface{} {
  return db[key]
}

func Set(key string, value interface{}) {
  db[key] = value
}

func SetWithTimeout(key string, value interface{}, timeout int64) {
  db[key] = value
}

package main

import (
	"github.com/garethstokes/fourtyeight/library"
  "fmt"
)

func main() {

  l := library.Store()
  l.OpenSession()
  defer l.CloseSession()

  l.DeleteExpiredPosts()

}

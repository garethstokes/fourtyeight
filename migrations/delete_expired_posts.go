package main

import (
	"github.com/garethstokes/fourtyeight/library"
)

func main() {

  l := library.Store()
  l.OpenSession()
  defer l.CloseSession()

  l.DeleteExpiredPosts()

}

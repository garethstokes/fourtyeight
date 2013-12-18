package main

import (
  "fmt"
  "github.com/garethstokes/fourtyeight/mail"
)

func main() {
  mail.Initialise()

  err := mail.SendWelcomeEmail("test@betechnology.com.au")
  if err != nil {
    fmt.Println("error")
    fmt.Println(err)
    return
  }

  fmt.Println("email sent successfully")
}

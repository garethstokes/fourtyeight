package mail

import (
  "fmt"
  "net/smtp"
)

var (
  auth smtp.Auth
)

func Initialise() {
  // Set up authentication information.
  auth = smtp.PlainAuth(
    "",
    "magic_elves@shortfuse.io",
    "i8the4kin$",
    "smtp.gmail.com",
  )
}

func send(to []string, message []byte) error {

  // Connect to the server, authenticate, set the sender 
  // and recipient, and send the email all in one step.
  err := smtp.SendMail(
    "smtp.gmail.com:25",
    auth,
    "magic_elves@shortfuse.io",
    to,
    message,
  )

  return err
}

func SendWelcomeEmail(toAddress string) error {
  err := send([]string { toAddress }, []byte("welcome to the drop!"))
  if err != nil {
    return err
  }

  admin := []string {
    "gareth@betechnology.com.au",
    "nick.g.skelton@gmail.com",
    "david_cave@hotmail.com",
  }

  message := fmt.Sprintf("a new user has signed up - %s", toAddress)
  fmt.Println(message)
  return send(admin, []byte(message))
}

package library

import (
  "testing"
  "time"
)

func post() * Post {
  post := new(Post)
  post.OwnerId = "@garrydanger"
  post.Image = "http://i.imgur.com/FudYBky.jpg"
  post.Text = "Took me a while to figure out that hand-situation."

  return post
}

func TestInsertAndFind(t * testing.T) {
  var library = new(Library)
  library.Schema = "fourtyeight_test"

  library.OpenSession()
  defer library.DestroyCollectionAndCloseSession()

  post := post()

  // 48 hours
  library.CreateFrom(post, 60 * 60 * 48)
  documents := library.FindDocumentsFor("@garrydanger")

  if len(documents) != 1 {
    t.Fatal("document not found.")
  }
}

func TestExpiration(t * testing.T) {
  var library = new(Library)
  library.Schema = "fourtyeight_test"

  library.OpenSession()
  defer library.DestroyCollectionAndCloseSession()

  post := post()
  library.CreateFrom(post, 1)
  documents := library.FindDocumentsFor("@garrydanger")
  if len(documents) != 1 {
    t.Fatal("document did not save.")
  }
  if documents[0].ExpirationDelta != 1 {
    t.Fatal("document did not save ExpirationDelta in correct form.")
  }

  time.Sleep(2 * time.Second)
  documents = library.FindDocumentsFor("@garrydanger")

  if len(documents) != 0 {
    t.Fatal("document did not expire.")
  }
}

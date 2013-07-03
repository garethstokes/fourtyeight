package library

import (
  "testing"
)

func TestInsertAndFind(t * testing.T) {
  var store = new(Store)

  store.OpenSession()
  defer store.DestroyCollectionAndCloseSession()

  post := new(Post)
  post.OwnerId = 69
  post.Image = "http://i.imgur.com/FudYBky.jpg"
  post.Text = "Took me a while to figure out that hand-situation."

  document := store.CreateFrom(post)
  t.Log(document)

  documents := store.FindDocumentsFor(69)
  t.Log(documents)
}

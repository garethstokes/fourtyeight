package library

import (
  "fmt"
  "testing"
  "time"
  "labix.org/v2/mgo/bson"
  "github.com/garethstokes/fourtyeight/personal"
)

func post() * Post {
  post := new(Post)
  post.OwnerId = "garrydanger"
  post.Image = "http://i.imgur.com/FudYBky.jpg"
  post.Text = "Took me a while to figure out that hand-situation."
  post.DateCreated = 23712938
  return post
}

func garrydanger() []personal.Person {
  var result = make([]personal.Person,1)
  gd := new(personal.Person)
  gd.Username = "garrydanger"

  result = append(result, * gd)
  return result
}

func TestInsertAndFind(t * testing.T) {
  var library = new(Library)
  library.Schema = "fourtyeight_test"

  library.OpenSession()
  defer library.DestroyCollectionAndCloseSession()

  post := post()

  // 48 hours
  library.CreateFrom(post, 60 * 60 * 48)
  documents := library.FindDocumentsFor(garrydanger(),0)

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
  documents := library.FindDocumentsFor(garrydanger(),0)
  if len(documents) != 1 {
    t.Fatal("document did not save.")
  }
  if documents[0].ExpirationDelta != 1 {
    t.Fatal("document did not save ExpirationDelta in correct form.")
  }

  time.Sleep(2 * time.Second)
  documents = library.FindDocumentsFor(garrydanger(),0)

  if len(documents) != 0 {
    t.Fatal("document did not expire.")
  }
}

func TestLikeEvenThoughImaCowboy(t * testing.T) {
 var library = new(Library)
  library.Schema = "fourtyeight_test"

  library.OpenSession()
  defer library.DestroyCollectionAndCloseSession()
  
  
  
  p := post()

  // 48 hours
  library.CreateFrom(p, 60 * 60 * 48)
  documents := library.FindDocumentsFor(garrydanger(), 0)

  if len(documents) != 1 {
    t.Fatal("document not found.")
  }

  var firstDoco = documents[0]
  
  comment := post()

  //fmt.Printf("bleeme: %s\n", (firstDoco.Key.(bson.ObjectId)).Hex())

  firstDocoNew := library.AddPost(comment, (firstDoco.Key.(bson.ObjectId)).Hex())
 
  fmt.Printf("addpost: %s\n", firstDocoNew)
  
  key := (firstDocoNew.Key.(bson.ObjectId)).Hex()

  fmt.Printf("new key: %s\n", key)

  result := library.LikePost(key, 0, "shredder")
  

  if len(result.Comments) != 1 {
    t.Fatal("comennt not found.")
  }

  fmt.Println(result.Comments[0].LikedBy)
  if len(result.Comments[0].LikedBy) != 1 {
    t.Fatal("liker not found.")
  }

}

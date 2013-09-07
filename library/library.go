package library

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
  "fmt"
  "github.com/garethstokes/fourtyeight/personal"
)

type Post struct {
  OwnerId string `json:"ownerId"`
  Image string `json:"imageUrl"`
  Text string `json:"text"`
}

type Library struct {
  session * mgo.Session
  collection * mgo.Collection
  Schema string
}

func Store() * Library {
  store := new( Library )
  return store
}

func (s * Library) OpenSession() {
  fmt.Print( "Library.OpenSession\n" )

    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    if len( s.Schema ) == 0 {
      s.Schema = "fourtyeight_development"
    }

    s.session = session
    s.collection = session.DB(s.Schema).C("documents")
}

func (s * Library) CloseSession() {
  fmt.Print( "Library.CloseSession\n" )
  s.session.Close()
}

func (s * Library) DestroyCollectionAndCloseSession() {
  fmt.Print( "Library.DestroyCollectionAndCloseSession\n" )
  s.collection.DropCollection()
  s.session.Close()
}

func (s * Library) CreateFrom(post * Post, expiry int64) * Document {
  fmt.Printf( "Library.CreateFrom :: %@\n", post )
  fmt.Printf( "Duration: %d\n", expiry)

  key := bson.NewObjectId()

  document := new(Document)
  document.Key = key.String()
  document.MainPost = post
  document.Comments = make([]Post, 0)
  document.ExpirationDelta = expiry
  document.DateCreated = time.Now().UTC().Unix()

  err := s.collection.Insert(document)
  if err != nil {
    panic(err)
  }

  return document
}

func (s * Library) AddPost(post * Post, documentKey string) * Document {
  fmt.Printf( "Library.AddPost :: to document, %s\n", documentKey )

  if bson.IsObjectIdHex(documentKey) == false {
    return nil
  }

  key := bson.ObjectIdHex( documentKey )

  document := new( Document )
  err :=s.collection.Find(bson.M{"key": key}).One( &document )
  if err == nil {
    document.Comments = append( document.Comments, *post )
    s.collection.Update( bson.M{"key": key}, document )
  } else {
    fmt.Printf( "ERROR: %s\n", err.Error() )
  }

  return document
}

func (s * Library) FindOne( id string ) * Document {
  fmt.Printf( "Library.FindOne :: %s\n", id )

  if bson.IsObjectIdHex(id) == false {
    return nil
  }

  document := new( Document )
  err :=s.collection.Find(bson.M{"key": bson.ObjectIdHex(id)}).One( &document )
  if err != nil {
    fmt.Printf( "ERROR: %s\n", err.Error() )
    return nil
  }

  return document
}

func (s * Library) FindDocumentsFor(users []personal.Person) []Document {

    var result = make([]Document, 100)

    var queries = make([]bson.M, len(users))
    for i := range users {
      queries[i] = bson.M{"mainpost.ownerid": users[i].Username}
    }

    var query = bson.M{"$or": queries}
    err := s.collection.Find(query).All(&result)
    if err != nil {
        panic(err)
    }

    documents := make([]Document, 0)
    for _, d := range result {
      if d.expired() {
        continue
      }

      documents = append(documents, d)
    }

    fmt.Printf("documents: %d\n", len(documents))
    return documents
}

func (s * Library) FindAllFor(userId string) []Document {

    var result = make([]Document, 100)

    err := s.collection.Find(bson.M{}).All(&result)
    if err != nil {
        panic(err)
    }

    return result
}

func (s * Library) DeleteOne(id string) error {
  return s.collection.Remove(bson.M{"key": bson.ObjectIdHex(id)})
}

func (s * Library) DeleteExpiredPosts() error {
  job := &mgo.MapReduce{
    Map:        "function() { emit(this.key, [7, this.expirationDelta]); }",
    Reduce:     "function(key, values) { return Array.sum(values[0]); }",
  }

  var result []struct { Id string; Value int64 }
  _, err := s.collection.Find(nil).MapReduce(job, &result)
  if err != nil {
    return err
  }

  for _, item := range result {
    fmt.Println(item)
  }

  return nil
}

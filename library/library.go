package library

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
  "fmt"
)

type Post struct {
  OwnerId string `json:"ownerId"`
  Image string `json:"imageUrl"`
  Text string `json:"text"`
  DateCreated time.Time `json:"dateCreated"`
}

type Document struct {
  LastUpdated time.Time `json:"lastUpdated"`
  MainPost * Post `json:"mainPost"`
  Comments []Post `json:"comments"`
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

func (s * Library) CreateFrom(post * Post) * Document {
  fmt.Printf( "Library.CreateFrom :: %@\n", post )

  document := new(Document)
  document.LastUpdated = time.Now().UTC()
  document.MainPost = post
  document.Comments = make([]Post, 0)

  err := s.collection.Insert(document)
  if err != nil {
    panic(err)
  }

  return document
}

func (s * Library) FindDocumentsFor(userId string) []Document {

    var result = make([]Document, 100)

    err := s.collection.Find(bson.M{"mainpost.ownerid": userId}).All(&result)
    if err != nil {
        panic(err)
    }

    return result
}

func (s * Library) FindAllFor(userId string) []Document {

    var result = make([]Document, 100)

    err := s.collection.Find(bson.M{}).All(&result)
    if err != nil {
        panic(err)
    }

    return result
}

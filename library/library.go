package library

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
)

type Post struct {
  OwnerId string
  Image string
  Text string
  Timestamp time.Time
}

type Document struct {
  LastUpdated time.Time
  MainPost * Post
  Comments []Post
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
    s.session.Close()
}

func (s * Library) DestroyCollectionAndCloseSession() {
    s.collection.DropCollection()
    s.session.Close()
}

func (s * Library) CreateFrom(post * Post) * Document {
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

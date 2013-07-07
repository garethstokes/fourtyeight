package library

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
)

type Post struct {
  OwnerId uint64
  Image string
  Text string
  Timestamp time.Time
}

type Document struct {
  LastUpdated time.Time
  MainPost * Post
  Comments []Post
}

type Store struct {
  session * mgo.Session
  collection * mgo.Collection
  Schema string
}

func (s * Store) OpenSession() {
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

func (s * Store) CloseSession() {
    s.session.Close()
}

func (s * Store) DestroyCollectionAndCloseSession() {
    s.collection.DropCollection()
    s.session.Close()
}

func (s * Store) CreateFrom(post * Post) * Document {
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

func (s * Store) FindDocumentsFor(userId uint64) []Document {

    var result = make([]Document, 100)

    err := s.collection.Find(bson.M{"mainpost.ownerid": userId}).All(&result)
    if err != nil {
        panic(err)
    }

    return result
}

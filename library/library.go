package library

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
  "fmt"
)

/*
  TODO: like api changes
      - include the mongo key with all posts, similar to a document
      - include likedby and followers for a post
*/
type Post struct {
  OwnerId string `json:"ownerId"`
  Image string `json:"imageUrl"`
  Text string `json:"text"`
  DateCreated int64 `json:"dateCreated"` 
  LikedBy []string `json:"likedBy"`
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
  document.Key = key
  document.MainPost = post
  document.Comments = make([]Post, 0)
  document.ExpirationDelta = expiry
  document.DateCreated = time.Now().UTC().Unix()
  document.LastModified = time.Now().UTC().Unix()

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
    document.LastModified = time.Now().UTC().Unix()
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

func (s * Library) LikePost(id string, position int, username string) * Document{
  fmt.Printf("LikePost :: Hi fucking llo\n")
  var document = s.FindOne(id)
  if document == nil{
    fmt.Printf("LikePost :: Returning nil\n")
    return nil
  }
  
  if(position <= 0){
    fmt.Printf("liking the main post\n")
    var mainPost = document.MainPost
    mainPost.LikedBy = append( mainPost.LikedBy, username )
    document.MainPost = mainPost;
    //update last modified
    document.LastModified = time.Now().UTC().Unix()
  
    s.collection.Update( bson.M{"key": bson.ObjectIdHex(id)}, document )
  
  }else{
    position--;
    fmt.Printf("liking %dth comment as %s\n", position, username)
    if document.Comments == nil{
      fmt.Printf("LikePost :: Comments is nil, cancelling the like\n")
      return nil
    }

    fmt.Printf("LikePost :: Comments is %d\n", len(document.Comments))
    
    if len(document.Comments) > position {
      fmt.Printf("if so\n")
      var post = document.Comments[position]

      post.LikedBy = append( post.LikedBy, username )
      fmt.Printf("LikePost :: LikedBy is %d\n", len(post.LikedBy))
      document.Comments[position] = post;
      //update last modified
      document.LastModified = time.Now().UTC().Unix()
  
      s.collection.Update( bson.M{"key": bson.ObjectIdHex(id)}, document )
      fmt.Printf("if so done\n")
    }
  }
  return document
}


func (s * Library) UnlikePost(id string, position int, username string) * Document{
  fmt.Printf("UnlikePost :: Starting\n")
  var document = s.FindOne(id)
  if document == nil{
    fmt.Printf("UnlikePost :: Returning nil couldnt find the document\n")
    return nil
  }
  
  if(position <= 0){
    fmt.Printf("unliking the main post\n")
    var mainPost = document.MainPost

    var filteredLikedBy = make([]string, 0)

    for _, value := range mainPost.LikedBy{
      if(value != username){
        filteredLikedBy = append(filteredLikedBy, value)
      }else{
        fmt.Printf("filtered one username out woo hoo\n")
      }

    }
    
    mainPost.LikedBy = filteredLikedBy
    document.MainPost = mainPost
    //update last modified
    document.LastModified = time.Now().UTC().Unix()
  
    s.collection.Update( bson.M{"key": bson.ObjectIdHex(id)}, document )
  
  }else{
    position--;
    fmt.Printf("liking %dth comment as %s\n", position, username)
    if document.Comments == nil{
      fmt.Printf("UnlikePost :: Comments is nil, cancelling the like\n")
      return nil
    }

    fmt.Printf("UnlikePost :: Comments is %d\n", len(document.Comments))
    
    if len(document.Comments) > position {
      fmt.Printf("if so\n")
      var post = document.Comments[position]

      var filteredLikedBy = make([]string, 0)

      for _, value := range post.LikedBy{
        if(value != username){
          filteredLikedBy = append(filteredLikedBy, value)
        }else{
          fmt.Printf("filtered one username out woo hoo\n")
        }

      } 

      post.LikedBy = filteredLikedBy
      document.Comments[position] = post
      //update last modified
      document.LastModified = time.Now().UTC().Unix()
  
      s.collection.Update( bson.M{"key": bson.ObjectIdHex(id)}, document ) 
    }
  }
  return document
}

func (s * Library) FindDocumentsFor(users []string, timestamp int) []Document {

    var result = make([]Document, 100)

    fmt.Printf("documents: %d\n", timestamp)

    var queries = make([]bson.M, len(users))
    for i := range users {
      queries[i] = bson.M{ "$and": []bson.M{
                      bson.M{"mainpost.ownerid": users[i]},
                      bson.M{"lastmodified": bson.M{ "$gt": timestamp }},
                   }}
    }

    var query = bson.M{"$or": queries}

    fmt.Printf("query: %s\n", query)

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

func (s * Library) FindPublicDocuments(timestamp int) []Document {

    var result = make([]Document, 100)

    fmt.Printf("public documents: %d\n", timestamp)
    
    //public documents are any that we three created
    var users = make([]string, 3)
    users = append(users, "shredder")
    users = append(users, "garrydanger")
    users = append(users, "caveman")

    var queries = make([]bson.M, len(users))
    for i := range users {
      queries[i] = bson.M{ "$and": []bson.M{
                      bson.M{"mainpost.ownerid": users[i]},
                      bson.M{"lastmodified": bson.M{ "$gt": timestamp }},
                   }}
    }

    var query = bson.M{"$or": queries}

    fmt.Printf("query: %s\n", query)

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

func (s * Library) DeleteOne(id string) error {
  return s.collection.Remove(bson.M{"key": bson.ObjectIdHex(id)})
}

func (s * Library) DeleteOneWithObjectId(id bson.ObjectId) error {
  return s.collection.Remove(bson.M{"key": id})
}

func (s * Library) DeleteExpiredPosts() error {
  var m = `
  function() {
    var d = new Date();
    var now = d.getTime() / 1000;
    var expiration = this.datecreated + this.expirationdelta;

    //emit(this.key, [now, expiration, expiration < now]);
    emit(this.key, expiration < now);
  }
  `

  job := &mgo.MapReduce{
    Map:        m,
    Reduce:     "function(k, v) { return 7; }",
  }

  var result []struct { Id interface{} "_id"; IsExpired bool "value" }
  //var result []map[string]interface{}
  _, err := s.collection.Find(nil).MapReduce(job, &result)
  if err != nil {
    return err
  }

  counter := 0

  for _, post := range result {
    if post.IsExpired {
      counter++;
      s.DeleteOneWithObjectId(post.Id.(bson.ObjectId))
    }
  }

  fmt.Printf("Deleted %d items\n", counter)

  return nil
}

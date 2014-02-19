package main

import (
  "fmt"
  "time"
  "strconv"
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/garethstokes/fourtyeight/library"
	"github.com/garethstokes/fourtyeight/personal"
	"github.com/garethstokes/fourtyeight/cache"
  "github.com/garethstokes/fourtyeight/push_notifications"
  "labix.org/v2/mgo/bson"
)

func LibraryController() {

  // GET Documents from the library for a user
  // Example: /library/{USER_TOKEN}/{TIMESTAMP}
  //
  // will find all relevent documents for a user, ordered
  // by the create_created
  //
  web.Get("/library/([0-9A-Z]{25})(/[0-9]+)?", func(ctx * web.Context, token string, timestamp string) {
    ctx.SetHeader("Content-Type", "application/json", true);

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    user := cache.Get("users", token )
    if user == nil {
      apiError( ctx, "INVALID_TOKEN" )
      return
    }

    // timestamp
    if len(timestamp) == 0 {
      var daysAgo int64 = 60 * 60 * 24 * 7 // 7 days ago
      var nowPlusDaysAgo int = int(time.Now().UTC().Unix() - daysAgo)
      timestamp = "/" + strconv.Itoa(nowPlusDaysAgo)
    }

    var ts, _ = strconv.Atoi(timestamp[1:])

    fmt.Println(time.Unix(int64(ts), 0))
    fmt.Println(time.Now())

    fmt.Println(ts)

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    person := user.(* personal.Person)
    posts := l.FindDocumentsFor(append(person.Followers, person.Username), ts)

    ok( ctx, posts )
  })

  // The public endpoint for the above action
  web.Get("/library/?", func(ctx * web.Context) {
    ctx.SetHeader("Content-Type", "application/json", true);

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    posts := l.FindPublicDocuments(0)

    ok( ctx, posts )
  })

  web.Get("/library/(/[0-9]+)?", func(ctx * web.Context, timestamp string) {
    ctx.SetHeader("Content-Type", "application/json", true);

    // timestamp
    if len(timestamp) == 0 {
      var daysAgo int64 = 60 * 60 * 24 * 7 // 7 days ago
      var nowPlusDaysAgo int = int(time.Now().UTC().Unix() - daysAgo)
      timestamp = "/" + strconv.Itoa(nowPlusDaysAgo)
    }

    var ts, _ = strconv.Atoi(timestamp[1:])

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    posts := l.FindPublicDocuments(ts)

    ok( ctx, posts )
  })

  // POST Create a document attached to a user
  // Example: localhost:8080/library/e51n4EZvN8KL7uoQUtmbWw==/document 
  //
  // will attach a expiry delta that needs to be specified
  //
  type PostWithExpiry struct {
    library.Post
    Expiry int64
  }
  web.Post("/library/(.+)/document", func(ctx * web.Context, token string) {
    ctx.SetHeader("Content-Type", "application/json", true)
  fmt.Println("POST2 :: Document")
    user := cache.Get("users", token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    person := user.(* personal.Person)

    post := new(PostWithExpiry)
    err := json.NewDecoder(ctx.Request.Body).Decode(&post)
    if err != nil {
			apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
			return
    }

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    p := new(library.Post)
    p.OwnerId = person.Username
    p.Image = post.Image
    p.Text = post.Text
    p.DateCreated = time.Now().UTC().Unix()

    document := l.CreateFrom( p, post.Expiry )

    // loop through followers and send them 
    // all a notifiation
    personalStore := personal.Store()
    personalStore.OpenSession()
    defer personalStore.CloseSession()

    documentId  := (document.Key.(bson.ObjectId)).Hex()
    go push_notifications.SendPushNotificationAboutAPost(person.Followers, person.Username + " just dropped", documentId)

    ok( ctx, document )
  })

  web.Get("/document/(.+)", func(ctx * web.Context, documentId string) {
    ctx.SetHeader("Content-Type", "application/json", true)

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    document := l.FindOne( documentId )
    if document == nil {
      apiError( ctx, "Incorrect document id" )
      return
    }

    ok( ctx, document )
  })

  web.Post("/library/(.+)/document/(.+)/post", func(ctx * web.Context, token string, documentId string) {
    ctx.SetHeader("Content-Type", "application/json", true)
    fmt.Println("POST :: Document")
    user := cache.Get("users", token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    post := new( library.Post )
    err := json.NewDecoder(ctx.Request.Body).Decode(&post)
    if err != nil {
      apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
      return
    }

    post.OwnerId = user.(* personal.Person).Username
    post.DateCreated = time.Now().UTC().Unix()
    
    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    //Create the post
    document := l.AddPost( post, documentId )

    /////// BUILD A LIST OF THE PEOPLE WHO SHOULD BE NOTIFIED
    /////// THIS SHOULD BE REFACTORED OUT TO SOMEWHERE NICER SOON COS ITS PROBABLY GONNA BE DUPLICATED
    person := user.(* personal.Person)

    //Add the owner of the post, if they are commenting on their own post, leave them out
    if(person.Username != document.MainPost.OwnerId){ 
      go push_notifications.SendPushNotificationAboutAPostToOne(document.MainPost.OwnerId, person.Username + " just commented on your post", documentId)
    }

    usersWhoShouldBeNotified := make([]string, 1)
    //go through all the commenters and add them, unless they are the owner of this new comment, 
    //or the owner of the post (who we have already added above)
    for _, p := range document.Comments{
      if(person.Username != p.OwnerId && document.MainPost.OwnerId != p.OwnerId){
        usersWhoShouldBeNotified = append (usersWhoShouldBeNotified , p.OwnerId)
      }
    }

    //notify the people above
    go push_notifications.SendPushNotificationAboutAPost(usersWhoShouldBeNotified, person.Username + " also commented on "+document.MainPost.OwnerId+"'s post", documentId)

    /////// END NOTIFICATION LOGIX


    ok( ctx, document )
  })

  web.Post("/library/(.+)/document/(.+)/like(/[0-9]+)?", func(ctx * web.Context, token string, documentId string, position string) {
    fmt.Printf("LIKE :: Document \n")
   
    ctx.SetHeader("Content-Type", "application/json", true)
    
    isLikingMainPost := false

    if len(position) == 0 {
      position = "/0"
      isLikingMainPost = true
    }

    fmt.Printf("\nLIKE :: Document position %s " , position)

    var posi, _ = strconv.Atoi(position[1:]) 
     
    user := cache.Get("users", token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }
 
    var usrName = user.(* personal.Person).Username

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()
  
    fmt.Printf("\nLIKE :: Document  " +documentId)
    // fmt.Printf("LIKE :: Position %d\n", position)
    fmt.Printf("\nLIKE :: usrName " +usrName)
  
    document := l.LikePost( documentId , posi, usrName)

    /////// BUILD A LIST OF THE PEOPLE WHO SHOULD BE NOTIFIED
    /////// THIS SHOULD BE REFACTORED OUT TO SOMEWHERE NICER SOON COS ITS PROBABLY GONNA BE DUPLICATED
    person := user.(* personal.Person)
    
    //Add the owner of the post, if they are liking on their own post, leave them out
    
    if (isLikingMainPost){
      //someone is liking the main post
      if(person.Username != document.MainPost.OwnerId){
        //the person who is liking the post is not the owner of the post
        go push_notifications.SendPushNotificationAboutAPostToOne(document.MainPost.OwnerId, person.Username + " just liked your post", documentId)
      }
      //notify the other people who like it already
      go push_notifications.SendPushNotificationAboutAPost(document.MainPost.LikedBy, person.Username + " also liked " + document.MainPost.OwnerId + "'s post", documentId)
      
    }else{ 
      //go through all the comments, find the one being liked
      for i, p :=range document.Comments{
        if(i==posi){
          //found the comment in question
          if(person.Username != p.OwnerId){
            //notify the original comment maker
            go push_notifications.SendPushNotificationAboutAPostToOne(p.OwnerId, person.Username + " just liked your comment", documentId)
          }
          //notify the other people who like it already
          go push_notifications.SendPushNotificationAboutAPost(p.LikedBy, person.Username + " also liked " + p.OwnerId + "'s comment", documentId)
        }
      }
    }
   
    /////// END NOTIFICATION LOGIX

    ok( ctx, document )
  })

  web.Post("/library/(.+)/document/(.+)/unlike(/[0-9]+)?", func(ctx * web.Context, token string, documentId string, position string) {
    fmt.Printf("UNLIKE :: Document \n")
   
    ctx.SetHeader("Content-Type", "application/json", true)
    
    if len(position) == 0 {
      position = "/0"
    }
    fmt.Printf("\nUNLIKE :: Document position %s " , position)

    var posi, _ = strconv.Atoi(position[1:]) 
    
    fmt.Printf("\nUNLIKE :: Document posi %d " , posi)

    user := cache.Get("users", token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }
 
    var usrName = user.(* personal.Person).Username

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()
  
    fmt.Printf("\nUNLIKE :: Document  " +documentId)
    // fmt.Printf("LIKE :: Position %d\n", position)
    fmt.Printf("\nUNLIKE :: usrName " +usrName)
  

    document := l.UnlikePost( documentId , posi, usrName)

    ok( ctx, document )
  })

  web.Post("/library/(.+)/delete/(.+)", func(ctx * web.Context, token string, documentId string) {
    ctx.SetHeader("Content-Type", "application/json", true)

    user := cache.Get("users", token)
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    err := l.DeleteOne(documentId)
    if err != nil {
      apiError(ctx, "There was a problem contacting library service.")
      fmt.Printf("ERROR: %s\n", err.Error())
      return
    }

    ok( ctx, true )
  })
}

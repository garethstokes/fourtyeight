package personal

import (
  "labix.org/v2/mgo/bson"
  "fmt"
  "errors"
  "github.com/garethstokes/fourtyeight/cache"
)

var (

  //   var queries = make([]bson.M, len(users))
  //   for i := range users {
  //     queries[i] = bson.M{ "$or": []bson.M{
  //                     bson.M{"username": bson.M{ "$regex": val }},
  //                     bson.M{"email": bson.M{ "$regex": val }},
  //                  }}
  //   }
  // db_columns = "user_id, username, email, avatar_url, loginToken, date_created"
  // db_predicate_token = fmt.Sprintf( "loginToken = ? limit 1" )
  // db_predicate = fmt.Sprintf( "username = ? limit 1" )
)

func (s * Personal) FillCacheWithLoginTokens(){

  s.logf( "personal.FillCacheWithLoginTokens" )
  var result = make([]Person, 0)
  query := bson.M{"logintoken": bson.M{ "$ne": "" }}

  err := s.collection.Find(query).All(&result)
  if err != nil {
      panic(err)
  }

  fmt.Printf("Warming cache with this many tokens: %d\n", len(result))
  for _, p := range result {
    cache.Set("users", p.LoginToken, & p)
  }
}

func (s * Personal) FindByString(val string) ([]Person){
  //This query could be better, 
  // + email search shoud only search chars BEFORE the @
  // + limit results to 100
  // + favour results that start with 'val' - this might involve doing 3 searchs?
  //TODO get a regex dungeon master to write a sick regex bra

  query := bson.M { "$or": []bson.M {
                                      bson.M{"username": bson.M{ "$regex": val }},
                                      bson.M{"email": bson.M{ "$regex": val }},
                                    }}

  s.logf( "Personal.FindByString :: value: %s", val)

  persons := make([]Person, 0)
  err :=s.collection.Find(query).All( &persons )
  if err != nil {
    fmt.Printf( "ERROR: %s\n", err.Error())
  }

  return persons
}

func (s * Personal) FindAll() ([]Person) {
  var result = make([]Person, 0)
  s.collection.Find(bson.M{}).All(&result)
  return result
}

func (s * Personal) FindByToken(token string) (* Person, error) {
  return s.findBy("logintoken", token)
}

func (s * Personal) FindByName( name string ) (* Person, error) {
  return s.findBy("username", name)
}

func (s * Personal) findBy(key string, val string) (* Person, error) {
  s.logf( "Personal.FindBy :: key: %s, value: %s", key, val)

  person := new( Person )
  err :=s.collection.Find(bson.M{key: val}).One( &person )
  if err != nil {
    fmt.Printf( "ERROR: %s\n", err.Error())
    return nil, err
  }

  return person, nil
}


func (s * Personal) GetLoggedInUser(loginToken string)(* Person, error) {
    
    // check if user is logged in to the cache

    var user * Person
  
    person := cache.Get("users", loginToken)
    if person == nil {

      //not in cache lets try the DB
      //GArry danger might have restarted the server
      //LOlcats!

      user, error := s.FindByToken(loginToken)
      
      if error != nil {
        s.error( error.Error() )
        return nil, error
      }
      
      if user == nil {
        err := errors.New("User with token not found")
        return nil, err
      }
    }else{

      user = person.(* Person)

    }

    return user, nil
}

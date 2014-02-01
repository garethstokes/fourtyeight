package push_notifications

import (
  "github.com/garethstokes/fourtyeight/cache"
)

type PushNotificationContent struct {
	Type string `json:"type"`
	Message string `json:"message"`
  PostIdentifier string `json:"postId"`
  Username string `json:"username"`
	//TODO put in some other data if you want, could even send new posts this way
	// eg. newPost Document `json:"newPost"`
}

func SendPushNotificationAboutAPostToOne(user string, message string, postid string){
   quickWrapper := make([]string, 1)
   quickWrapper = append(quickWrapper, user)
   SendPushNotificationAboutAPost(quickWrapper, message, postid)
}

func SendPushNotificationAboutAUserToOne(user string, message string, username string){
   quickWrapper := make([]string, 1)
   quickWrapper = append(quickWrapper, user)
   SendPushNotificationAboutAUser(quickWrapper, message, username)
}

func SendPushNotificationAboutAPost(users []string, message string, postid string){

    pn := new(PushNotificationContent)
    pn.Type = "USER"
    pn.Message = message
    pn.PostIdentifier = postid

    SendPushNotificationTo(users, pn)
}

func SendPushNotificationAboutAUser(users []string, message string, username string){

    pn := new(PushNotificationContent)
    pn.Type = "POST"
    pn.Message = message
    pn.Username = username

    SendPushNotificationTo(users, pn)
}

func SendPushNotificationTo(users []string, pn * PushNotificationContent){
    iosDeviceTokens := make([]string, 0)
    androidDeviceTokens := make([]string, 0)

    
    //gather the tokens for each user and each platform
    for _, user := range users{
      //ios
      iosToken := cache.Get("apns", user)
      if(iosToken!=nil){
        // TODO batch ios notifications same as android
        // iosDeviceTokens = append(iosDeviceTokens, iosToken.(string))
        // TODO SEND THE POSTID AS PART OF THE PAYLOAD
        go SendPushNotificationIOS(iosToken.(string), pn)
      }
      //android
      androidToken := cache.Get("apns_android", user)
      if(androidToken!=nil){
        androidDeviceTokens = append(androidDeviceTokens, androidToken.(string))
      }
    }
 
    //ios
    if(len(iosDeviceTokens) > 0){
      // TODO batch ios notifications same as android
      // sendPushNotificationTo(deviceToken.(string), person.Username)
    }

    //android
    if(len(androidDeviceTokens) > 0){
      go SendPushNotificationAndroid(0, pn, androidDeviceTokens)
    }
}


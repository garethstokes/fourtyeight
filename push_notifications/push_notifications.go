package push_notifications

import (
  "github.com/garethstokes/fourtyeight/cache"
)

type PushNotificationContent struct {
	Type string `json:"type"`
	Message string `json:"message"`
	PostIdentifier string `json:"postId"`
	//TODO put in some other data if you want, could even send new posts this way
	// eg. newPost Document `json:"newPost"`
}

func SendPushNotificationToOne(user string, message string, postid string){
   quickWrapper := make([]string, 1)
   quickWrapper = append(quickWrapper, user)
   SendPushNotificationTo(quickWrapper, message, postid)
}

func SendPushNotificationTo(users []string, message string, postid string){
    iosDeviceTokens := make([]string, 0)
    androidDeviceTokens := make([]string, 0)

    pn := new(PushNotificationContent)
    pn.Type = "POST"
    pn.Message = message
    pn.PostIdentifier = postid

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


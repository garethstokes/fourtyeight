package apns_android

import (
	"net/http"
	"io/ioutil"
	"fmt" 
	"bytes"
	"encoding/json"
)

const GOOGLE_URL = "https://android.googleapis.com/gcm/send"

//Google supports 4 at most, so you can use more but it will only keep 4
const COLLAPSE_KEY_NEW_POSTS = "New posts"
const COLLAPSE_KEY_LIKES = "Someone has liked your post"
const COLLAPSE_KEY_NEW_COMMENTS = "Someone has commented on your post"
// const COLLAPSE_KEY_NEW_POSTS = "New posts"

type AndroidPushNotification struct{
	//TODO time_to_live int64 `json:"time_to_live"`
	RegistrationIds []string `json:"registration_ids"`
	Data * PushNotificationContent `json:"data"`
}

type PushNotificationContent struct {
	Message string `json:"message"`
	PostIdentifier string `json:"post"`
	//TODO put in some other data if you want, could even send new posts this way
	// eg. newPost Document `json:"newPost"`
}
//time to live is the time between now and when the post dies in seconds
func SendNotification( timeToLive int64, message string, post string, recipients []string ){

	client := &http.Client{
	}

	content := new(PushNotificationContent)
	content.Message = message

	bodyPost := new(AndroidPushNotification)
	bodyPost.RegistrationIds = recipients
	bodyPost.Data = content

	b, err := json.Marshal(bodyPost)
	if err != nil {
		fmt.Printf( "Problem marshalling json: %s", err.Error())
	}

	rdr := bytes.NewReader(b)

	req, err := http.NewRequest("POST", GOOGLE_URL, rdr)
	if err != nil {
		fmt.Printf( "Problem creating request: %s", err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "key=AIzaSyDXIbZpR-lrGgBEsAqQCUGpB3oMR6E_Ysk")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf( "Problem doing post: %s", err.Error())
	}

	//body is a []byte
	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil {
		fmt.Printf( "Problem reading post body: %s", err.Error())
	}
	
	//bodyBuf *Buffer
	bodyBuf := bytes.NewBuffer(body)
	
	//bodyStr string
	bodyStr, err := bodyBuf.ReadString('\n')
	if err != nil {
		fmt.Printf( "Problem reading post bodyStr: %s", err.Error())
	}
	
	fmt.Printf( "Response from google push notifications:\n %s\n", bodyStr )


}

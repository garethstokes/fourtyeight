package push_notifications

import (
	"testing"
)

// Create a new Payload that specifies simple text,
// a badge counter, and a custom notification sound.
func mockPayload() (payload *Payload) {
	payload = NewPayload()
	payload.Alert = "You have mail!"
	payload.Badge = 42
	payload.Sound = "bingbong.aiff"
	return
}

// See the commentary in push_notification.go for information
// on why we're testing a badge of value 0.
func mockZeroBadgePayload() (payload *Payload) {
	payload = mockPayload()
	payload.Badge = 0
	return
}

// Create a new AlertDictionary. Apple recommends you not use
// the more complex alert style unless absolutely necessary.
func mockAlertDictionary() (dict *AlertDictionary) {
	args := make([]string, 1)
	args[0] = "localized args"

	dict = NewAlertDictionary()
	dict.Body = "Complex Message"
	dict.ActionLocKey = "Play a Game!"
	dict.LocKey = "localized key"
	dict.LocArgs = args
	dict.LaunchImage = "image.jpg"
	return
}

func TestBasicAlert(t *testing.T) {
	payload := mockPayload()
	pn := NewPushNotification()

	pn.AddPayload(payload)

	bytes, _ := pn.ToBytes()
	json, _ := pn.PayloadJSON()
	if len(bytes) != 82 {
		t.Error("expected 82 bytes; got", len(bytes))
	}
	if len(json) != 69 {
		t.Error("expected 69 bytes; got", len(json))
	}
}

func TestAlertDictionary(t *testing.T) {
	dict := mockAlertDictionary()
	payload := mockPayload()
	payload.Alert = dict

	pn := NewPushNotification()
	pn.AddPayload(payload)

	bytes, _ := pn.ToBytes()
	json, _ := pn.PayloadJSON()
	if len(bytes) != 207 {
		t.Error("expected 207 bytes; got", len(bytes))
	}
	if len(json) != 194 {
		t.Error("expected 194 bytes; got", len(bytes))
	}
}

func TestCustomParameters(t *testing.T) {
	payload := mockPayload()
	pn := NewPushNotification()

	pn.AddPayload(payload)
	pn.Set("foo", "bar")

	if pn.Get("foo") != "bar" {
		t.Error("unable to set a custom property")
	}
	if pn.Get("not_set") != nil {
		t.Error("expected a missing key to return nil")
	}

	bytes, _ := pn.ToBytes()
	json, _ := pn.PayloadJSON()
	if len(bytes) != 94 {
		t.Error("expected 94 bytes; got", len(bytes))
	}
	if len(json) != 81 {
		t.Error("expected 81 bytes; got", len(json))
	}
}

func TestZeroBadgeChangesToNegativeOne(t *testing.T) {
	payload := mockZeroBadgePayload()
	pn := NewPushNotification()
	pn.AddPayload(payload)

	if payload.Badge != -1 {
		t.Error("expected 0 badge value to be converted to -1; got", payload.Badge)
	}
}

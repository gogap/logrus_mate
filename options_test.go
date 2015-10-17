package logrus_mate

import (
	"testing"
)

func TestGetStringValue(t *testing.T) {
	opts := Options{"key": "string value"}

	val, err := opts.String("key")

	if err != nil {
		t.Error("get string value of key failed")
		return
	}

	if val != "string value" {
		t.Error("the value of key is not 'string value'")
	}
}

func TestGetIntValue(t *testing.T) {
	opts := Options{"key": 1986}

	val, err := opts.Int("key")

	if err != nil {
		t.Error("get int value of key failed")
		return
	}

	if val != 1986 {
		t.Error("the value of key is not 1986")
	}
}

func TestGetFloatValue(t *testing.T) {
	opts := Options{"key": 88.88}

	val, err := opts.Float64("key")

	if err != nil {
		t.Error("get float value of key failed")
		return
	}

	if val != 88.88 {
		t.Error("the value of key is not 88.88")
	}
}

type user struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestGetObjectValue(t *testing.T) {

	var userA = user{
		Name:  "zeal",
		Email: "xujinzheng@gmail.com",
	}

	userB := user{}

	opts := Options{"key": userA}

	err := opts.Object("key", &userB)

	if err != nil {
		t.Error("get object value of key failed")
		return
	}

	if userB.Name != userA.Name || userB.Email != userA.Email {
		t.Error("the value of key's object's value is not correct")
	}
}

func TestToObjectValue(t *testing.T) {

	userA := user{}

	opts := Options{"name": "zeal", "email": "xujinzheng@gmail.com"}

	err := opts.ToObject(&userA)

	if err != nil {
		t.Error("get object value of key failed")
		return
	}

	if userA.Name != "zeal" || userA.Email != "xujinzheng@gmail.com" {
		t.Error("the options could not convert into object")
	}
}

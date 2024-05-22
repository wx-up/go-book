package grpc

import "testing"

func Test_User(t *testing.T) {
	u := &User{}
	if _, ok := u.Contact.(*User_Email); ok {
		t.Log(u.GetEmail())
	}
}

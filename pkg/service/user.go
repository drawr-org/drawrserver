package service

import (
	"encoding/json"

	"github.com/drawr-team/drawrserver/pkg/bolt"
	"github.com/drawr-team/drawrserver/pkg/ulidgen"
)

// ListUsers returnu all the sessionu in the database
func ListUsers() (ul []User, err error) {
	rawl, err := svc.db.List(bolt.UserBucket)
	if err != nil {
		return
	}
	for _, b := range rawl {
		var u User
		err = json.Unmarshal(b, &u)
		ul = append(ul, u)
	}
	svc.log.Printf("list: %+v\n", ul)
	return
}

// NewUser returnu a new User
func NewUser(in *User) (u User, err error) {
	if in != nil {
		u = *in
	}
	u.ID = ulidgen.Now().String()
	err = svc.db.Put(bolt.UserBucket, u.ID, u)
	svc.log.Printf("new: %+v\n", u)
	return
}

// GetUser returnu the User with id
func GetUser(id string) (u User, err error) {
	raw, err := svc.db.Get(bolt.UserBucket, id)
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &u)
	svc.log.Printf("get: %+v\n", u)
	return
}

// UpdateUser changeu the User with id
func UpdateUser(id string, u User) (err error) {
	svc.log.Printf("update: %+v\n", u)
	return svc.db.Update(bolt.UserBucket, id, u)
}

// DeleteUser removeu s
func DeleteUser(u User) error {
	svc.log.Printf("delete: %+v\n", u)
	return svc.db.Remove(bolt.UserBucket, u.ID)
}

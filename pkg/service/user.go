package service

import (
	"github.com/boltdb/bolt"
	"github.com/drawr-team/drawrserver/pkg/model"
	"github.com/drawr-team/drawrserver/pkg/ulidgen"

	log "github.com/golang/glog"
	"github.com/simia-tech/boltx"
)

// ListUsers returnu all the sessionu in the database
func ListUsers() (ul []model.User, err error) {
	log.Info("list users")

	err = svc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucketKey)
		_, _, err := boltx.ForEach(b, &model.User{}, func(k []byte, v interface{}) (boltx.Action, error) {
			ul = append(ul, *v.(*model.User))
			return boltx.ActionReturn, nil
		})
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	})
	log.V(2).Infof("payload: %+v", ul)
	return
}

// NewUser returnu a new User
func NewUser(in *model.User) (u model.User, err error) {
	log.Info("new user")
	if in != nil {
		u = *in
	}
	u.ID = ulidgen.Now().String()

	err = svc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(sessionBucketKey)
		return boltx.PutModel(b, []byte(u.ID), &u)
	})

	log.V(2).Info("payload: %+v\n", u)
	return
}

// GetUser returnu the User with id
func GetUser(id string) (u model.User, err error) {
	log.Info("get user")

	err = svc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucketKey)
		ok, err := boltx.GetModel(b, []byte(id), &u)
		if !ok {
			err = ErrNotFound
		}
		return err
	})

	log.V(2).Info("payload: %+v\n", u)
	return
}

// UpdateUser changeu the User with id
func UpdateUser(id string, u model.User) (err error) {
	log.Info("update user")
	defer log.V(2).Info("payload: %+v\n", u)

	return svc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucketKey)
		return boltx.PutModel(b, []byte(id), &u)
	})
}

// DeleteUser removeu s
func DeleteUser(u model.User) error {
	log.Info("delete user")
	log.V(2).Info("payload: %+v\n", u)
	return boltx.DeleteFromBucket(svc.db, userBucketKey, []byte(u.ID))
}

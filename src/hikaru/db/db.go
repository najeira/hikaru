package db

import (
	"appengine"
	"appengine/datastore"
)

func KeyZero(c appengine.Context, kind string, parent *datastore.Key) *datastore.Key {
	return datastore.NewKey(c, kind, "", 0, parent)
}

func KeyIncomplete(c appengine.Context, kind string, parent *datastore.Key) *datastore.Key {
	return datastore.NewKey(c, kind, "", 0, parent)
}

func KeyStr(c appengine.Context, kind string, stringID string, parent *datastore.Key) *datastore.Key {
	return datastore.NewKey(c, kind, stringID, 0, parent)
}

func KeyInt(c appengine.Context, kind string, intID int64, parent *datastore.Key) *datastore.Key {
	return datastore.NewKey(c, kind, "", intID, parent)
}

func Put(c appengine.Context, key *datastore.Key, data interface{}) (*datastore.Key, error) {
	return datastore.Put(c, key, data)
}

func PutAsync(c appengine.Context, key *datastore.Key, data interface{}) (<-chan *datastore.Key, <-chan error) {
	key_ch := make(chan *datastore.Key, 1)
	err_ch := make(chan error, 1)
	go func() {
		new_key, err := datastore.Put(c, key, data)
		err_ch <- err
		key_ch <- new_key
	}()
	return key_ch, err_ch
}

// func Get(c appengine.Context, key *Key, dst interface{}) error {
// func GetMulti(c appengine.Context, key []*Key, dst interface{}) error {
// func Put(c appengine.Context, key *Key, src interface{}) (*Key, error) {
// func PutMulti(c appengine.Context, key []*Key, src interface{}) ([]*Key, error) {
// func Delete(c appengine.Context, key *Key) error {
// func DeleteMulti(c appengine.Context, key []*Key) error {
// func NewIncompleteKey(c appengine.Context, kind string, parent *Key) *Key {
// func NewKey(c appengine.Context, kind, stringID string, intID int64, parent *Key) *Key {
// // func AllocateIDs(c appengine.Context, kind string, parent *Key, n int) (low, high int64, err error) {




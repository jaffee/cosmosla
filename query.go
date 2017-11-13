package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (m *Main) Read(collection *mgo.Collection) error {
	start := time.Now()
	n, err := collection.Find(nil).Count()
	if err != nil {
		return errors.Wrap(err, "counting all records")
	}
	fmt.Printf("%v total records: %v\n", time.Since(start), n)

	res := Person{}
	start = time.Now()
	err = collection.Find(nil).One(&res)
	if err != nil {
		return errors.Wrap(err, "finding first record")
	}
	fmt.Printf("%v first result:\n%v\n", time.Since(start), res)

	start = time.Now()
	res = Person{}
	err = collection.Find(bson.M{"tiles.os": true, "tiles.du": true}).One(&res)
	if err != nil {
		return errors.Wrap(err, "finding first segment record")
	}
	fmt.Printf("%v first segment record:\n%v\n", time.Since(start), res)

	res = Person{}
	start = time.Now()
	// idbytes, err := hex.DecodeString("5a09a414f21bc91a5c7e7669")
	// if err != nil {
	// 	return errors.Wrap(err, "decoding hex")
	// }
	err = collection.Find(bson.M{"_id": bson.ObjectIdHex("5a09a414f21bc91a5c7e7669")}).One(&res)
	if err != nil {
		return errors.Wrap(err, "finding single record")
	}
	fmt.Printf("%v single record %v\n", time.Since(start), res)

	return nil
}

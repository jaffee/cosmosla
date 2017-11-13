package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (m *Main) Read(collection *mgo.Collection) error {
	n, err := collection.Find(nil).Count()
	if err != nil {
		return errors.Wrap(err, "counting all records")
	}
	fmt.Printf("total records: %v\n", n)

	start := time.Now()
	res := Person{}
	err = collection.Find(bson.M{"tiles.os": true, "tiles.du": true}).One(&res)
	duration := time.Since(start)
	if err != nil {
		return errors.Wrap(err, "finding record")
	}
	fmt.Printf("%v results in %v\n", res, duration)
	return nil
}

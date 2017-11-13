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
	n, err := collection.Find(bson.M{"tiles.os": true, "tiles.du": true}).Count()
	duration := time.Since(start)
	if err != nil {
		return errors.Wrap(err, "finding record")
	}
	fmt.Printf("%v results in %v\n", n, duration)
	return nil
}

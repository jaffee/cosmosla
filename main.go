package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/Azure/go-autorest/autorest/utils"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	database string
	password string
)

func init() {
	database = utils.GetEnvVarOrExit("AZURE_DATABASE")
	password = utils.GetEnvVarOrExit("AZURE_DATABASE_PASSWORD")
}

// Person represents a document in the collection
type Person struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Tiles map[string]bool
}

type Main struct {
	Num    int
	Insert bool
	Query  bool
	Seed   int64
}

func main() {
	m := Main{}
	flag.IntVar(&m.Num, "num", 10000, "number of docs to insert")
	flag.BoolVar(&m.Insert, "insert", false, "do insertions")
	flag.BoolVar(&m.Query, "query", false, "do queries")
	flag.Int64Var(&m.Seed, "seed", 1, "seed for rng")
	flag.Parse()
	err := m.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Main) Run() error {
	// DialInfo holds options for establishing a session with a MongoDB cluster.
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s.documents.azure.com:10255", database)}, // Get HOST + PORT
		Timeout:  60 * time.Second,
		Database: database, // It can be anything
		Username: database, // Username
		Password: password, // PASSWORD
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		},
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return errors.Errorf("Can't connect to mongo, go error %v\n", err)
	}

	defer session.Close()

	// SetSafe changes the session safety mode.
	// If the safe parameter is nil, the session is put in unsafe mode, and writes become fire-and-forget,
	// without error checking. The unsafe mode is faster since operations won't hold on waiting for a confirmation.
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode.
	session.SetSafe(&mgo.Safe{})

	// get collection
	collection := session.DB(database).C("people")

	if m.Insert {
		err := m.Write(collection)
		if err != nil {
			return errors.Wrap(err, "inserting")
		}
	}

	if m.Query {
		err := m.Read(collection)
		if err != nil {
			return errors.Wrap(err, "querying")
		}
	}
	return nil
}

func (m *Main) Write(collection *mgo.Collection) error {
	g := NewGenerator(m.Seed)
	// insert documents into collection
	for i := 0; i < m.Num; i++ {
		err := collection.Insert(g.Person())
		if err != nil {
			return errors.Wrap(err, "inserting person")
		}
	}
	return nil
}

type Generator struct {
	r *rand.Rand
}

func NewGenerator(seed int64) Generator {
	src := rand.NewSource(seed)
	return Generator{
		r: rand.New(src),
	}
}

var letters = "abcdefghijklmnopqrstuvwxyz12345678"

func (t Generator) Tile() string {
	length := t.r.Intn(4) + 2
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		ret[i] = letters[t.r.Intn(len(letters))]
	}
	return string(ret)
}

func (t Generator) Person() *Person {
	tiles := make(map[string]bool)
	num := t.r.Intn(930) + 70
	for i := 0; i < num; i++ {
		tiles[t.Tile()] = true
	}
	return &Person{
		Tiles: tiles,
	}
}

// // update document
// updateQuery := bson.M{"_id": result.Id}
// change := bson.M{"$set": bson.M{"fullname": "react-native"}}
// err = collection.Update(updateQuery, change)
// if err != nil {
// 	log.Fatal("Error updating record: ", err)
// 	return errors.Wrap(err, "updating")

// }

// // delete document
// err = collection.Remove(updateQuery)
// if err != nil {
// 	return errors.Wrap(err, "deleting record")
// }

// // Get Document from collection
// result := Person{}
// err = collection.Find(bson.M{"fullname": "react"}).One(&result)
// if err != nil {
// 	log.Fatal("Error finding record: ", err)
// 	return
// }

// fmt.Println("Description:", result.Description)

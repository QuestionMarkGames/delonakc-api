package database

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
)

func New(host string, port int, db string) *MongoDB {
	url := fmt.Sprintf("%s:%d/%s", host, port, db)
	session, err := mgo.Dial(url)

	fmt.Println(url)
	
	if err != nil {
		log.Println("Mongodb create failed:")
		log.Fatal(err)
	}

	if err := session.Ping(); err != nil {
		log.Println("Mongodb connect failed:")
		log.Fatal(err)
	}

	fmt.Println("====================== ")
	fmt.Println("mongodb connect sucess ")
	fmt.Println("====================== ")

	return &MongoDB{ DB: session, DatabaseName: db }
}

type MongoDB struct {
	DB *mgo.Session
	C *mgo.Collection
	DatabaseName string
}

func (m *MongoDB) Close() {
	m.DB.Close()
}

func (m *MongoDB) SetCollection(c string) {
	collection := m.DB.DB(m.DatabaseName).C(c)
	m.C = collection
}
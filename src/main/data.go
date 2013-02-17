package main

import (
    "labix.org/v2/mgo"
    "fmt"
)

func NewMongoEmailDataWriter() (medw *MongoEmailDataWriter){
    medw = new(MongoEmailDataWriter)

    session, err := mgo.Dial("localhost")

    if err != nil {
        panic(err)
    }

    medw.email = session.DB("inbox").C("email")

    index := mgo.Index{
        Key: []string{"user_id", "uid"},
        Unique: true,
        DropDups: true,
    }

    err = medw.email.EnsureIndex(index)

    if err != nil {
        panic(err)
    }

    return medw
}

type MongoEmailDataWriter struct {
    email *mgo.Collection
}

func (medw *MongoEmailDataWriter) WriteEmailData(ed EmailData) {
    fmt.Println(ed.Uid)
    err := medw.email.Insert(&ed)

    if err != nil {
        fmt.Println("Unable to insert, skipping: ", err)
    }
}

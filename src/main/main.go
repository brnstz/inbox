package main

import (
    "flag"
    "fmt"
    "time"
)

type Inbox struct {
    Senders map[string]*Sender
}

type Sender struct {
    Count int
    Ids []uint32
}

type EmailData struct {
    Domain string
    Address string
    Name string
    Uid uint32
    Date time.Time
    Size uint32
    User_Id int
}

type EmailDataWriter interface {
    WriteEmailData(ed EmailData)
}

type TerminalEmailDataWriter struct {
    test bool
}

func (tedw *TerminalEmailDataWriter) WriteEmailData(ed EmailData) {
    fmt.Printf("%s,%s,%s,%d,%s,%d\n", ed.Domain, ed.Address, ed.Name, ed.Uid, ed.Date.Format("2006-01-02"), ed.Size)
}

func main() {
    var (
        server string
        user string
        pw string
        user_id int
    )

    flag.StringVar(&server, "server", "", "IMAP server hostname")
    flag.StringVar(&user, "user", "", "IMAP username")
    flag.StringVar(&pw, "pw", "", "IMAP pw")
    flag.IntVar(&user_id, "user_id", 0, "mongo user_id")
    flag.Parse()

    e := NewEmail(server, user, pw, user_id)
    //edw := new(TerminalEmailDataWriter)
    edw := NewMongoEmailDataWriter()
    e.GetCounts(1, edw)
}

package main

import (
    //"github.com/sbinet/go-imap/go1/imap"
    //"code.google.com/p/go-imap/imap"
    "imap"
    "crypto/tls"
    "flag"
    "fmt"
    "reflect"
    //"os"
)

type Inbox struct {
    Senders map[string]Sender
}

type Sender struct {
    Count int
    Ids []int
}



func login(server string, user string, pw string) (c *imap.Client, err error) {
    c, err = imap.DialTLS(server, new(tls.Config))

    if (err != nil) {
        return nil, err
    }

    _, err = c.Login(user, pw)

    if (err != nil) {
        return nil, err
    }

    return c, nil
}

func fetchIds(c *imap.Client) (ids []uint32) {
    seq_set, err := imap.NewSeqSet("1:*")

    if (err != nil) {
        panic(err)
    }

    cmd, err := imap.Wait(c.Fetch(seq_set, "FLAGS"))

    if (err != nil) {
        panic(err)
    }

    for _, row := range(cmd.Data) {

        v, ok := row.Fields[0].(uint32)
        if (ok) {
            ids = append(ids, v)
        }
    }

    return ids
}

//func parseMessage

func fetchMessage(c *imap.Client, id uint32) {
    seq_set, err := imap.NewSeqSet(fmt.Sprintf("%d", id))

    if err != nil {
        panic(err)
    }

    cmd, err := imap.Wait(c.Fetch(seq_set, "ALL"))

    if err != nil {
        panic(err)
    }

    /*
    res, err := cmd.Result(1)

    if err != nil {
        panic(err)
    }
    */

    //fmt.Println(reflect.TypeOf(cmd.Data[0].Fields[2]))
    //fmt.Println(cmd.Data[0].Fields[2])
    //fmt.Println(cmd.Data[0])
    //fmt.Println(reflect.TypeOf(cmd.Data[0]))
    //fmt.Println(reflect.TypeOf(imap.AsFieldMap(cmd.Data[0].Fields[2])["ENVELOPE"]))
    //fmt.Println(imap.AsFieldMap(cmd.Data[0].Fields[2])["ENVELOPE"])
    /*
    var (
        emailName string
        emailAddr string
        msgId int
    )
    */
    /*
    for k, v := range imap.AsList(imap.AsList(cmd.Data[0].Fields[2])[1]) {
        fmt.Println(k, " ", v)
    }
    */
    emailAddrInfo := imap.AsList(imap.AsList(imap.AsList(imap.AsList(cmd.Data[0].Fields[2])[1])[2])[0])
    emailName := imap.AsString(emailAddrInfo[0])
    emailAddress := fmt.Sprintf("%s@%s", imap.AsString(emailAddrInfo[2]), imap.AsString(emailAddrInfo[3]))
    fmt.Println(emailName)
    fmt.Println(emailAddress)

    //fmt.Println(res)
    //fmt.Println(res.Decoded)
    //fmt.Println(cmd.String())

}

func blah() {

    fmt.Println(reflect.TypeOf(1))
    _ = new(Inbox)
}

func main() {
    var (
        server string
        user string
        pw string
    )

    flag.StringVar(&server, "server", "", "IMAP server hostname")
    flag.StringVar(&user, "user", "", "IMAP username")
    flag.StringVar(&pw, "pw", "", "IMAP pw")
    flag.Parse()

    c, err := login(server, user, pw)

    if (err != nil) {
        panic(err)
    }

    _, err = imap.Wait(c.Select("INBOX", true))

    if (err != nil) {
        panic(err)
    }

    ids := fetchIds(c)

    for _, id := range(ids) {
        fmt.Printf("ID SAYS: %d\n", id)
        fetchMessage(c, id)
    }

    /*
    res, err := cmd.Result(0)
    fmt.Println(res.String())
    */
     /*
    res, err := cmd.Result(0)
    fmt.Println(res.String())

    for row := range(res)
    */

}

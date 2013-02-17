package main

import (
    //"github.com/sbinet/go-imap/go1/imap"
    //"code.google.com/p/go-imap/imap"
    "imap"
    "crypto/tls"
    "flag"
    "fmt"
    "reflect"
    "encoding/json"
    "os"
    "io/ioutil"
)

type Inbox struct {
    Senders map[string]*Sender
}

type Sender struct {
    Count int
    Ids []uint32
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

func fetchIds(c *imap.Client) (uids []uint32) {
    seq_set, err := imap.NewSeqSet("1:*")

    if (err != nil) {
        panic(err)
    }

    cmd, err := imap.Wait(c.UIDFetch(seq_set, "FLAGS"))

    if (err != nil) {
        panic(err)
    }

    for _, row := range(cmd.Data) {
        fmt.Println(row)
        fmt.Println(row.MessageInfo().UID)

        /*
        v, ok := row.Fields[0].(uint32)
        if (ok) {
            ids = append(ids, v)
        }
        */
        uids = append(uids, row.MessageInfo().UID)
    }

    return uids
}

//func parseMessage

func fetchMessage(c *imap.Client, uid uint32) (emailName string, emailAddr string) {
    seq_set, err := imap.NewSeqSet(fmt.Sprintf("%d", uid))

    if err != nil {
        panic(err)
    }

    cmd, err := imap.Wait(c.UIDFetch(seq_set, "ALL"))
    rsp := cmd.Data[0]
    fmt.Println(rsp.MessageInfo().Attrs["ENVELOPE"])
    emailAddrInfo := imap.AsList(imap.AsList(imap.AsList(imap.AsList(rsp.MessageInfo().Attrs["ENVELOPE"])[2]))[0])
    emailName = imap.AsString(emailAddrInfo[0])
    emailAddr = fmt.Sprintf("%s@%s", imap.AsString(emailAddrInfo[2]), imap.AsString(emailAddrInfo[3]))

    return emailName, emailAddr
}

func blah() {

    fmt.Println(reflect.TypeOf(1))
    _ = new(Inbox)
    os.Exit(1)
}

func addToInbox(inbox *Inbox, uid uint32, emailName string, emailAddr string) {

    if inbox.Senders == nil {
        inbox.Senders = map[string]*Sender {}
    }

    full := fmt.Sprintf("%s %s", emailName, emailAddr)

    sender, found := inbox.Senders[full]
    if found {
        sender.Count++
        sender.Ids = append(sender.Ids, uid)
    } else {
        sender = new(Sender)
        sender.Count = 1
        sender.Ids = append(sender.Ids, uid)
        inbox.Senders[full] = sender
    }
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

    uids := fetchIds(c)

    var emailAddr, emailName string
    inbox := new(Inbox)
    for _, uid := range(uids) {
        fmt.Printf("UID SAYS: %d\n", uid)
        emailName, emailAddr = fetchMessage(c, uid)
        addToInbox(inbox, uid, emailName, emailAddr)
    }

    b, err := json.Marshal(inbox)

    if (err != nil) {
        panic(err)
    }

    ioutil.WriteFile("inbox.json", b, 0644)

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

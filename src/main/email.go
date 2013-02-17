package main

import (
    "github.com/sbinet/go-imap/go1/imap"
    "crypto/tls"
    "fmt"
    //"os"
)

type Email struct {
    conn *imap.Client
}

func NewEmail(server string, user string, pw string) (e *Email) {
    e = new(Email)

    var err error

    e.conn, err = e.login(server, user, pw)
    if err != nil {
        panic(err)
    }

    _, err = imap.Wait(e.conn.Select("INBOX", true))
    if err != nil {
        panic(err)
    }

    return e
}

func (e *Email) login(server string, user string, pw string) (c *imap.Client, err error) {
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

func (e *Email) ParseFetchResp(resp *imap.Response) (ed EmailData, err error) {
    emailInfo := imap.AsList(imap.AsList(imap.AsList(imap.AsList(resp.MessageInfo().Attrs["ENVELOPE"])[2]))[0])

    alias := imap.AsString(emailInfo[2])

    ed.Name = imap.AsString(emailInfo[0])
    ed.Domain = imap.AsString(emailInfo[3])
    ed.Address = fmt.Sprintf("%s@%s", alias, ed.Domain)
    ed.Uid = resp.MessageInfo().UID
    ed.Date = resp.MessageInfo().InternalDate
    ed.Size = resp.MessageInfo().Size
    // FIXME
    ed.User_Id = 1234

    return ed, nil
}

func (e *Email) GetCounts(minUid int, edw EmailDataWriter) (err error){
    seqSet, err := imap.NewSeqSet(fmt.Sprintf("%d:*", minUid))

    if err != nil {
        return err
    }

    cmd, err := imap.Wait(e.conn.UIDFetch(seqSet, "ALL"))

    if err != nil {
        return err
    }
    for _, resp := range cmd.Data {
        ed, _ := e.ParseFetchResp(resp)
        edw.WriteEmailData(ed)
        /* FIXME test
        if i > 10 {
            os.Exit(1)
        }
        */
    }

    return nil
}

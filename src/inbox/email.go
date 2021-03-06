package inbox

import (
	"crypto/tls"
	"fmt"
	"github.com/sbinet/go-imap/go1/imap"
)

type Email struct {
	conn    *imap.Client
	User_Id int
}

func NewEmail(server string, user string, pw string, user_id int) (e *Email) {
	e = new(Email)

	e.User_Id = user_id

	var err error

	e.conn, err = e.login(server, user, pw)
	if err != nil {
		panic(err)
	}

	_, err = imap.Wait(e.conn.Select("INBOX", false))
	if err != nil {
		panic(err)
	}

	return e
}

func (e *Email) login(server string, user string, pw string) (c *imap.Client, err error) {
	c, err = imap.DialTLS(server, new(tls.Config))

	if err != nil {
		return nil, err
	}

	_, err = c.Login(user, pw)

	if err != nil {
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
	ed.User_Id = e.User_Id

	return ed, nil
}

func (e *Email) getCountsOneLoop(minUid uint32, edw EmailDataWriter) (lastUid uint32, err error) {
	seqSet, err := imap.NewSeqSet(fmt.Sprintf("%d:*", minUid))

	if err != nil {
		return 0, err
	}

	cmd, err := imap.Wait(e.conn.UIDFetch(seqSet, "ALL"))

	if err != nil {
		return 0, err
	}
	for _, resp := range cmd.Data {
		ed, _ := e.ParseFetchResp(resp)
		edw.WriteEmailData(ed)
		lastUid = ed.Uid
	}

	return lastUid, nil
}

func (e *Email) GetCounts(minUid uint32, edw EmailDataWriter) (err error) {
	var lastUid uint32
	for {
		lastUid, err = e.getCountsOneLoop(minUid, edw)

		if err != nil {
			return err
		}

		if minUid == lastUid {
			break
		}

		minUid = lastUid
	}

	return nil
}

func (e *Email) DeleteMany(uids []uint32) (err error) {
	seqSet := new(imap.SeqSet)
	for i := 0; i < len(uids); i++ {
		seqSet.AddNum(uids[i])
	}

	if err != nil {
		return err
	}
	flags := imap.NewFlagSet(`\Deleted`, `\Seen`)

	cmd, err := imap.Wait(e.conn.UIDStore(seqSet, "+FLAGS", flags))
	fmt.Println(cmd)

	if err != nil {
		return err
	}

	cmd, err = imap.Wait(e.conn.Expunge(nil))
	fmt.Println(cmd)

	if err != nil {
		return err
	}

	cmd, err = imap.Wait(e.conn.Close(true))
	fmt.Println(cmd)

	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"flag"
	//"fmt"
	"inbox"
)

func main() {
	var (
		server  string
		user    string
		pw      string
		user_id int
	)

	flag.StringVar(&server, "server", "", "IMAP server hostname")
	flag.StringVar(&user, "user", "", "IMAP username")
	flag.StringVar(&pw, "pw", "", "IMAP pw")
	flag.IntVar(&user_id, "user_id", 0, "mongo user_id")
	flag.Parse()

	//e := inbox.NewEmail(server, user, pw, user_id)
	//e.GetCounts(1, ed)
	//edw := inbox.TerminalEmailDataWriter(true)
	ed := inbox.NewMongoEmailData("localhost", "inbox")

	ed.GetEmailData(user_id, "domain")
	/*results, err := ed.GetEmailData(user_id, "domain")
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		fmt.Println(result)
	}
	*/
}

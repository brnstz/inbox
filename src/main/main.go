package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"inbox"
	"io/ioutil"
	"net/http"
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

	mongoData := inbox.NewMongoEmailData("localhost", "inbox")

	http.HandleFunc("/load.json", func(w http.ResponseWriter, r *http.Request) {
		imapCon := inbox.NewEmail(server, user, pw, user_id)
		mongoData.DeleteOldData(user_id)
		imapCon.GetCounts(1, mongoData)
		fmt.Fprint(w, `{"success": 1}`)
	})

	http.HandleFunc("/readCounts.json", func(w http.ResponseWriter, r *http.Request) {
		results, _ := mongoData.GetEmailData(user_id, "domain")
		jsonBytes, _ := json.MarshalIndent(results, "", " ")
		fmt.Fprintf(w, "%s", jsonBytes)
	})

	http.HandleFunc("/readAddressCounts.json", func(w http.ResponseWriter, r *http.Request) {
		results, _ := mongoData.GetEmailData(user_id, "address")
		jsonBytes, _ := json.MarshalIndent(results, "", " ")
		fmt.Fprintf(w, "%s", jsonBytes)
	})

	http.HandleFunc("/delete.json", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("Body: %s\n", body)
		r.Body.Close()
		var gev inbox.GetEmailValue
		json.Unmarshal(body, &gev)

		imapCon := inbox.NewEmail(server, user, pw, user_id)
		err := imapCon.DeleteMany(gev.Uids)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		} else {
			fmt.Fprint(w, "it worked\n")
		}
	})
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":80", nil)
}

Work in progress
================
Use MongoDB to map-reduce an IMAP inbox into groupings by sender domain, 
address, etc. Show who is using the most space in your inbox. Provides a 
web UI to bulk archive these emails.


Depends On
==========
- http://golang.org/
- https://github.com/sbinet/go-imap
- http://labix.org/mgo
- http://www.mongodb.org/

Getting Started
==============
- go get github.com/sbinet/go-imap/go1/imap
- go get labix.org/v2/mgo
- go get labix.org/v2/mgo/bson
- Have a mongod running locally (TODO: make configurable)
- Edit run.sh and insert your IMAP server, user, and pw.
- Have port 80 open and know your sudo pw (TODO: make configurable)
- ./run.sh
- Load your emails into MongoDB: 
    - Hit http://localhost/load.json and wait for success 
- Take a look at your emails:
    - Go to http://localhost/static/
- Hit archive button

TODO
====
- Create login screen
- Creating loading overlay
- Integrate with gmail oauth token
- Fix simultaneous connections bug
- Make buttons more response
- Add "friends" list (by default assume emails not junk)
- Dynamic create user name based on login
- Save data over time rather than re-creating

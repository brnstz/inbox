#!/bin/sh

# Login details for IMAP server
IMAP_USERNAME='FIXME'
IMAP_PW='FIXME'
IMAP_SERVER='FIXME'

# Any integer. This is used to associate your IMAP data in the mongodb
USER_ID=1234

export GOPATH=`pwd`

# Runs sudo to run on port 80
go build main && sudo ./main --server $IMAP_SERVER --user $IMAP_USERNAME --pw $IMAP_PW --user_id $USER_ID

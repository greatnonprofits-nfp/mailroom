package ticket

import (
	"testing"

	_ "github.com/nyaruka/mailroom/services/tickets/mailgun"
	_ "github.com/nyaruka/mailroom/services/tickets/zendesk"
	"github.com/nyaruka/mailroom/testsuite"
	"github.com/nyaruka/mailroom/testsuite/testdata"
	"github.com/nyaruka/mailroom/web"
)

func TestTicketClose(t *testing.T) {
	testsuite.Reset()
	db := testsuite.DB()

	// create 2 open tickets and 1 closed one for Cathy across two different ticketers
	testdata.InsertOpenTicket(db, testdata.Org1, testdata.Cathy, testdata.Mailgun, "Need help", "Have you seen my cookies?", "17", testdata.Admin)
	testdata.InsertOpenTicket(db, testdata.Org1, testdata.Cathy, testdata.Zendesk, "More help", "Have you seen my cookies?", "21", nil)
	testdata.InsertClosedTicket(db, testdata.Org1, testdata.Cathy, testdata.Zendesk, "Old question", "Have you seen my cookies?", "34", testdata.Editor)

	web.RunWebTests(t, "testdata/close.json", nil)
}

func TestTicketReopen(t *testing.T) {
	testsuite.Reset()
	db := testsuite.DB()

	// create 2 closed tickets and 1 open one for Cathy
	testdata.InsertClosedTicket(db, testdata.Org1, testdata.Cathy, testdata.Mailgun, "Need help", "Have you seen my cookies?", "17", testdata.Admin)
	testdata.InsertClosedTicket(db, testdata.Org1, testdata.Cathy, testdata.Zendesk, "More help", "Have you seen my cookies?", "21", nil)
	testdata.InsertOpenTicket(db, testdata.Org1, testdata.Cathy, testdata.Zendesk, "Old question", "Have you seen my cookies?", "34", testdata.Editor)

	web.RunWebTests(t, "testdata/reopen.json", nil)
}

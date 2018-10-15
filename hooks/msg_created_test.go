package hooks

import (
	"testing"

	"github.com/nyaruka/mailroom"
	"github.com/nyaruka/mailroom/testsuite"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
)

func TestMsgCreated(t *testing.T) {
	mailroom.Config = mailroom.NewMailroomConfig()
	mailroom.Config.AttachmentDomain = "foo.bar.com"
	defer func() { mailroom.Config = nil }()

	// add a URN for cathy so we can test all urn sends
	testsuite.DB().MustExec(
		`INSERT INTO contacts_contacturn(identity, path, scheme, priority, contact_id, org_id) 
		                          VALUES('tel:12065551212', '12065551212', 'tel', 10, $1, 1)`,
		Cathy)

	// TODO: test setting reply_to_id

	tcs := []HookTestCase{
		HookTestCase{
			Actions: ContactActionMap{
				Cathy: []flows.Action{
					actions.NewSendMsgAction(newActionUUID(), "Hello World", nil, []string{"yes", "no"}, true),
				},
				Evan: []flows.Action{
					actions.NewSendMsgAction(newActionUUID(), "Hello Attachments", []string{"image/png:/images/image1.png"}, nil, true),
				},
			},
			Assertions: []SQLAssertion{
				SQLAssertion{
					SQL:   "select count(*) from msgs_msg where text='Hello World' and contact_id = $1 and metadata = $2",
					Args:  []interface{}{Cathy, `{"quick_replies":["yes","no"]}`},
					Count: 2,
				},
				SQLAssertion{
					SQL:   "select count(*) from msgs_msg where text='Hello Attachments' and contact_id = $1 and attachments[1] = $2",
					Args:  []interface{}{Evan, "image/png:https://foo.bar.com/images/image1.png"},
					Count: 1,
				},
			},
		},
	}

	RunActionTestCases(t, tcs)
}

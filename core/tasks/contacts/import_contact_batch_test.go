package contacts_test

import (
	"testing"

	_ "github.com/nyaruka/mailroom/core/handlers"
	"github.com/nyaruka/mailroom/core/tasks/contacts"
	"github.com/nyaruka/mailroom/testsuite"
	"github.com/nyaruka/mailroom/testsuite/testdata"

	"github.com/stretchr/testify/require"
)

func TestImportContactBatch(t *testing.T) {
	ctx := testsuite.CTX()
	rt := testsuite.RT()
	db := rt.DB

	importID := testdata.InsertContactImport(db, testdata.Org1)
	batchID := testdata.InsertContactImportBatch(db, importID, []byte(`[
		{"name": "Norbert", "language": "eng", "urns": ["tel:+16055740001"]},
		{"name": "Leah", "urns": ["tel:+16055740002"]}
	]`))

	task := &contacts.ImportContactBatchTask{ContactImportBatchID: batchID}

	err := task.Perform(ctx, rt, testdata.Org1.ID)
	require.NoError(t, err)

	testsuite.AssertQueryCount(t, db, `SELECT count(*) FROM contacts_contact WHERE name = 'Norbert' AND language = 'eng'`, nil, 1)
	testsuite.AssertQueryCount(t, db, `SELECT count(*) FROM contacts_contact WHERE name = 'Leah' AND language IS NULL`, nil, 1)
}

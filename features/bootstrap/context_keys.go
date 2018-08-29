package bootstrap

import (
	"net/http"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jmoiron/sqlx"
	"github.com/palantir/stacktrace"
)

type keysContext struct {
	selfURL  string
	db       *sqlx.DB
	response *http.Response
}

type key struct {
	KeyID     string
	Value     string
	CreatedAt string
	ExpiresAt string
}

var k = key{}

func (kc *keysContext) resetResponse(interface{}, error) {
	if kc.response != nil {
		kc.response.Body.Close()
		kc.response = nil
	}
}

func (kc *keysContext) resetDb(interface{}) {
	truncateKeys(kc.db)
}

// RegisterKeysContext registers steps related to "keys"
// with the godog suite
func RegisterKeysContext(s *godog.Suite, selfURL string, database *sqlx.DB) {
	kc := &keysContext{
		selfURL: selfURL,
		db:      database,
	}

	s.BeforeScenario(kc.resetDb)
	s.AfterScenario(kc.resetResponse)
	s.Step(`^there are keys:$`, kc.thereAreKeys)
}

func (kc *keysContext) thereAreKeys(keysTable *gherkin.DataTable) error {
	var result []key
	members := keysTable.Rows[0]
	for i := 1; i < len(keysTable.Rows); i++ {
		row := k
		SerializeTableRow(&row, members, keysTable.Rows[i])
		result = append(result, row)
	}

	for _, item := range result {
		if err := removeKey(kc.db, item.KeyID); err != nil {
			return stacktrace.NewError("Unable to remove keys", err)
		}

		if err := insertKey(kc.db, item.KeyID, item.Value, item.CreatedAt, item.ExpiresAt); err != nil {
			return stacktrace.NewError("Unable to insert keyss", err)
		}
	}

	return nil
}

//
// Helpers
//
func truncateKeys(db *sqlx.DB) error {
	_, err := db.Exec(`TRUNCATE keys CASCADE`)

	return err
}

func removeKey(db *sqlx.DB, keyID string) error {
	query := `DELETE FROM keys WHERE id = $1`

	_, err := db.Exec(query, keyID)

	if err != nil {
		return err
	}

	return nil
}

func insertKey(db *sqlx.DB, keyID string, value string, createdAt string, expiresAt string) error {
	query := `
        INSERT INTO keys
        (
            id,
			value,
            created_at,
			expires_at
        )
        VALUES ($1, $2, $3, $4)
    `

	_, err := db.Exec(
		query,
		keyID,
		value,
		createdAt,
		expiresAt,
	)

	if err != nil {
		return err
	}

	return nil
}

package scanner

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

type PG struct {
	Conn *sql.DB
}

var pg *PG

func TestLoadDB(t *testing.T) {
	db, err := sql.Open("postgres", "postgresql://postgres:test@localhost/template1")
	if err != nil {
		t.Fatalf("Couldnt open connection to test db\n%s", err)
	}
	pg = &PG{Conn: db}
}

type Unnested struct {
	Username string
}

func TestQuery(t *testing.T) {
	r := []*Unnested{}
	q := `SELECT 'test' as username;`

	err := Query(pg.Conn, &r, q)
	if err != nil {
		t.Fatalf("Error scanning query result into struct \n%s", err)
	}

	fmt.Println(r[0])
}

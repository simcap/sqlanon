package main

import (
	"reflect"
	"testing"
)

func TestParseInsertStatement(t *testing.T) {
	cases := []struct {
		in  string
		out *insertStmt
	}{
		{
			in:  "INSERT INTO actor VALUES (1, 'Penelope', 'Guiness', '2013-05-26 14:47:57.62');",
			out: &insertStmt{table: "actor", values: []interface{}{int64(1), "Penelope", "Guiness", "2013-05-26 14:47:57.62"}},
		},
		{
			in:  "INSERT INTO actor (id, first_name, last_name, updated) VALUES (1, 'Penelope', 'Guiness', '2013-05-26 14:47:57.62');",
			out: &insertStmt{table: "actor", columns: []string{"id", "first_name", "last_name", "updated"}, values: []interface{}{int64(1), "Penelope", "Guiness", "2013-05-26 14:47:57.62"}},
		},
	}

	for _, c := range cases {
		result, err := parseInsertStmt([]byte(c.in))
		if err != nil {
			t.Fatalf("error parsing %s\n%s", c.in, err)
		}
		if got, want := result, c.out; !reflect.DeepEqual(got, want) {
			t.Fatalf("parsing [%s]:\ngot\n%s\n\nwant\n%s", c.in, got, want)
		}
	}
}

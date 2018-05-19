package main

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	dump, err := parseDump("sample.sql")
	if err != nil {
		t.Fatal(err)
	}

	inserts := dump.insertIntos
	if got, want := len(inserts), 13; got != want {
		t.Fatalf("got %d want %d", got, want)
	}

	expected := &insert{
		table:  "names",
		fields: []string{"id", "\"language\"", "entity_type", "entity_id", "full_name", "first_name", "middle_name", "last_name", "alias", "abbreviation", "short_name", "prefix", "suffix"},
		values: []interface{}{int64(1), "en-US", "affiliations", int64(1), "Major League Baseball", ident("NULL"), ident("NULL"), ident("NULL"), ident("NULL"), ident("NULL"), ident("NULL"), ident("NULL"), ident("NULL")},
	}
	if got, want := inserts[0], expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot\n%#v\n\nwant\n%#v", got, want)
	}
}

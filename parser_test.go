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
		table: "names",
		fields: []sqlField{
			&field{n: "id", v: int64(1)},
			&field{n: "\"language\"", v: "en-US"},
			&field{n: "entity_type", v: "affiliations"},
			&field{n: "entity_id", v: int64(1)},
			&field{n: "full_name", v: "Major League Baseball"},
			&field{n: "first_name", v: ident("NULL")},
			&field{n: "middle_name", v: ident("NULL")},
			&field{n: "last_name", v: ident("NULL")},
			&field{n: "alias", v: ident("NULL")},
			&field{n: "abbreviation", v: ident("NULL")},
			&field{n: "short_name", v: ident("NULL")},
			&field{n: "prefix", v: ident("NULL")},
			&field{n: "suffix", v: ident("NULL")},
		},
	}
	if got, want := inserts[0], expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot\n%#v\n\nwant\n%#v", got, want)
	}
}

package main

import (
	"reflect"
	"testing"
)

func TestScanCreateTable(t *testing.T) {
	text := []byte(`CREATE TABLE users (
    id serial NOT NULL,
    "language" character varying(100) NOT NULL,
    name character varying(100) NOT NULL,
    age integer NOT NULL,
    date timestamp
);`)

	scan := newScanner(text[12:])
	stmt, err := scanCreateTable(scan)
	if err != nil {
		t.Fatal(err)
	}

	expected := &createStmt{table: "users", columns: []*column{
		{name: "id", category: numberCol},
		{name: "\"language\"", category: stringCol},
		{name: "name", category: stringCol},
		{name: "age", category: numberCol},
		{name: "date", category: timeCol},
	}}
	if got, want := stmt, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot\n%#v\n\nwant\n%#v", got, want)
	}
}

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

	expectedColDef := map[string]*column{
		"names.id":           &column{name: "id", category: numberCol},
		"names.\"language\"": &column{name: "\"language\"", category: stringCol},
		"names.entity_type":  &column{name: "entity_type", category: stringCol},
		"names.entity_id":    &column{name: "entity_id", category: numberCol},
		"names.full_name":    &column{name: "full_name", category: stringCol},
		"names.first_name":   &column{name: "first_name", category: stringCol},
		"names.middle_name":  &column{name: "middle_name", category: stringCol},
		"names.last_name":    &column{name: "last_name", category: stringCol},
		"names.alias":        &column{name: "alias", category: stringCol},
		"names.abbreviation": &column{name: "abbreviation", category: stringCol},
		"names.short_name":   &column{name: "short_name", category: stringCol},
		"names.prefix":       &column{name: "prefix", category: stringCol},
		"names.suffix":       &column{name: "suffix", category: stringCol},
	}
	if got, want := dump.columnDefinitions, expectedColDef; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot\n%#v\n\nwant\n%#v", got, want)
	}
}

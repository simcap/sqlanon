package main

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type ident string

type dump struct {
	insertIntos       []*insert
	columnDefinitions map[string]*column
}

func newDump() *dump {
	return &dump{columnDefinitions: make(map[string]*column)}
}

type sqlField interface {
	name() string
	value() interface{}
}

type field struct {
	n string
	v interface{}
}

func (f *field) name() string       { return f.n }
func (f *field) value() interface{} { return f.v }

type colCategory uint

const (
	stringCol colCategory = iota
	numberCol
	timeCol
	otherCol
)

type column struct {
	name     string
	category colCategory
}

type createStmt struct {
	table   string
	columns []*column
}

type insert struct {
	table  string
	fields []sqlField
}

func (i *insert) addField(name string, val interface{}) {
	i.fields = append(i.fields, &field{n: name, v: val})
}

func (i *insert) names() (out []string) {
	for _, f := range i.fields {
		out = append(out, f.name())
	}
	return
}

func (i *insert) write(w io.Writer) {
	w.Write([]byte("INSERT INTO "))
	w.Write([]byte(i.table))
	w.Write([]byte(" ("))
	w.Write([]byte(strings.Join(i.names(), ", ")))
	w.Write([]byte(") VALUES ("))
	for index, f := range i.fields {
		switch v := f.value().(type) {
		case int64:
			w.Write([]byte(fmt.Sprintf("%d", v)))
		case float64:
			w.Write([]byte(fmt.Sprintf("%f", v)))
		case string:
			w.Write([]byte(fmt.Sprintf("'%s'", v)))
		case ident:
			w.Write([]byte(fmt.Sprintf("%s", v)))
		}
		if index != len(i.fields)-1 {
			w.Write([]byte(", "))
		}
	}
	w.Write([]byte(");"))
}

func parseDump(filename string) (*dump, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(content))
	s.Init(file, content, nil, 0)

	d := newDump()
	for {
		pos, tok, lit := s.Scan()
		switch tok {
		case token.EOF:
			return d, nil
		case token.IDENT:
			switch lit {
			case "INSERT":
				if in, err := scanInsert(s); err != nil {
					return d, fmt.Errorf("scanning insert at %s: %s", fset.Position(pos), err)
				} else {
					d.insertIntos = append(d.insertIntos, in)
				}
			case "CREATE":
				pos, tok, lit := s.Scan()
				if tok == token.IDENT && lit == "TABLE" {
					stmt, err := scanCreateTable(s)
					if err != nil {
						return d, fmt.Errorf("scanning create table at %s: %s", fset.Position(pos), err)
					}
					for _, c := range stmt.columns {
						key := fmt.Sprintf("%s.%s", stmt.table, c.name)
						d.columnDefinitions[key] = c
					}
				}
			}
		}
	}
}

var errWithinStatementEOF = errors.New("EOF within statement")

func newScanner(line []byte) (s scanner.Scanner) {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(line))
	s.Init(file, line, nil, 0)
	return
}

func scanCreateTable(s scanner.Scanner) (*createStmt, error) {
	stmt := new(createStmt)
	if _, tok, lit := s.Scan(); tok != token.IDENT {
		return stmt, errors.New("expecting table name")
	} else {
		stmt.table = lit
	}

	if _, tok, lit := s.Scan(); tok != token.LPAREN {
		return stmt, fmt.Errorf("expecting left paren after table name got token %s with lit %q", tok, lit)
	}

	current := new(column)
start:
	for {
		_, tok, lit := s.Scan()
		switch tok {
		case token.EOF:
			return stmt, errWithinStatementEOF
		case token.IDENT, token.STRING:
			current.name = lit
			_, _, nextLit := s.Scan()
			switch nextLit {
			case "character", "text", "varchar", "char":
				current.category = stringCol
			case "integer", "smallint", "bigint", "decimal", "serial", "numeric", "real", "double precision", "bigserial":
				current.category = numberCol
			case "timestamp", "date", "time", "interval":
				current.category = timeCol
			}
			for {
				_, tok, _ := s.Scan()
				switch tok {
				case token.COMMA:
					stmt.columns = append(stmt.columns, current)
					current = new(column)
					goto start
				case token.SEMICOLON:
					stmt.columns = append(stmt.columns, current)
					current = new(column)
					return stmt, nil
				}
			}
		}
	}
}

func scanInsert(s scanner.Scanner) (*insert, error) {
	in := new(insert)
	if _, tok, lit := s.Scan(); tok != token.IDENT && lit != "INTO" {
		return in, errors.New("expecting INTO keyword")
	}

	if _, tok, lit := s.Scan(); tok != token.IDENT {
		return in, errors.New("expecting table name")
	} else {
		in.table = lit
	}

	if _, tok, lit := s.Scan(); tok != token.LPAREN {
		return in, fmt.Errorf("expecting left paren after table name got token %s with lit %q", tok, lit)
	}

	var names []string
	for {
		_, tok, lit := s.Scan()
		if tok == token.RPAREN {
			break
		}
		if tok == token.COMMA {
			continue
		}
		if tok == token.EOF {
			return in, errWithinStatementEOF
		}
		names = append(names, lit)
	}

	if _, tok, lit := s.Scan(); tok != token.IDENT && lit != "VALUES" {
		return in, fmt.Errorf("expecting VALUES keyword got token %s with lit %q", tok, lit)
	}

	if _, tok, _ := s.Scan(); tok != token.LPAREN {
		return in, errors.New("expecting left paren after VALUES")
	}

	var values []interface{}
	var current interface{}
	for {
		_, tok, lit := s.Scan()
		switch tok {
		case token.EOF:
			return in, errWithinStatementEOF
		case token.COMMA:
			values = append(values, current)
			current = nil
			continue
		case token.INT:
			i, err := strconv.ParseInt(lit, 10, 64)
			if err != nil {
				return in, errors.New("cannot convert int value")
			}
			current = i
		case token.FLOAT:
			f, err := strconv.ParseFloat(lit, 64)
			if err != nil {
				return in, errors.New("cannot convert float value")
			}
			current = f
		case token.CHAR, token.STRING:
			current = lit[1 : len(lit)-1]
		case token.IDENT:
			current = ident(lit)
		case token.RPAREN:
			values = append(values, current)
			current = nil
			if len(names) != len(values) {
				return nil, fmt.Errorf("expecting same count of names and values to insert at names[%v] values[%v]", names, values)
			}
			for i, n := range names {
				in.addField(n, values[i])
			}
			return in, nil
		}
	}
}

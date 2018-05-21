package main

import (
	"errors"
	"fmt"
	"go/token"
	"strconv"
)

type insertStmt struct {
	table   string
	columns []string
	values  []interface{}
}

func parseInsertStmt(b []byte) (*insertStmt, error) {
	s := newScanner(b)
	if _, _, lit := s.Scan(); lit != "INSERT" {
		return nil, errors.New("expecting INSERT keyword")
	}
	if _, _, lit := s.Scan(); lit != "INTO" {
		return nil, errors.New("expecting INTO keyword")
	}

	stmt := new(insertStmt)

	if _, tok, lit := s.Scan(); tok != token.IDENT {
		return nil, errors.New("expecting table name")
	} else {
		stmt.table = lit
	}

	var hasColumns bool
	if _, tok, lit := s.Scan(); tok == token.LPAREN {
		hasColumns = true
		for {
			_, tok, lit = s.Scan()
			switch tok {
			case token.COMMA:
				continue
			case token.RPAREN:
				goto jump
			case token.EOF:
				return stmt, errWithinStatementEOF
			case token.IDENT, token.STRING:
				stmt.columns = append(stmt.columns, lit)
			}
		}
	}

jump:
	if hasColumns {
		if _, tok, lit := s.Scan(); tok != token.IDENT && lit != "VALUES" {
			return stmt, fmt.Errorf("expecting VALUES keyword got token %s with lit %q", tok, lit)
		}
	}

	if _, tok, _ := s.Scan(); tok != token.LPAREN {
		return stmt, errors.New("expecting left paren after VALUES")
	}

	var current interface{}
	for {
		_, tok, lit := s.Scan()
		switch tok {
		case token.EOF:
			return stmt, errWithinStatementEOF
		case token.COMMA:
			stmt.values = append(stmt.values, current)
			current = nil
			continue
		case token.INT:
			i, err := strconv.ParseInt(lit, 10, 64)
			if err != nil {
				return stmt, errors.New("cannot convert int value")
			}
			current = i
		case token.FLOAT:
			f, err := strconv.ParseFloat(lit, 64)
			if err != nil {
				return stmt, errors.New("cannot convert float value")
			}
			current = f
		case token.CHAR, token.STRING:
			current = lit[1 : len(lit)-1]
		case token.IDENT:
			current = ident(lit)
		case token.RPAREN:
			stmt.values = append(stmt.values, current)
			return stmt, nil
		}
	}
}

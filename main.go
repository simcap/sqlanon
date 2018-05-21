package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	dumpFilepath string
	dryRun       bool
)

func init() {
	flag.StringVar(&dumpFilepath, "f", "", "SQL dump filepath")
	flag.BoolVar(&dryRun, "d", false, "Dry run")
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	dump, err := parseDump(dumpFilepath)
	if err != nil {
		log.Fatal(err)
	}

	dumpFile, err := os.Open(dumpFilepath)
	if err != nil {
		log.Fatal(err)
	}

	anonFilename := dumpFilepath + ".anonymized"
	if filepath.Ext(dumpFilepath) == ".sql" {
		anonFilename = dumpFilepath[:len(dumpFilepath)-4] + ".anonymized.sql"
	}

	anonFile, err := os.Create(anonFilename)
	if err != nil {
		log.Fatal(err)
	}

	var anonymizer anonymizer
	if dryRun {
		anonymizer = new(noopAnonymizer)
	} else {
		anonymizer = new(stringScrambler)
	}

	writer := newAnonWriter(dumpFile, anonFile, anonymizer)
	writer.setDump(dump)

	if err := writer.write(); err != nil {
		log.Fatal(err)
	}
}

type anonWriter struct {
	scanner *bufio.Scanner
	out     *bufio.Writer
	anon    anonymizer
	dump    *dump
}

func newAnonWriter(original io.Reader, out io.Writer, anons ...anonymizer) *anonWriter {
	w := &anonWriter{
		scanner: bufio.NewScanner(original),
		out:     bufio.NewWriter(out),
		dump:    newDump(),
		anon:    new(noopAnonymizer),
	}
	if len(anons) > 0 {
		w.anon = anons[0]
	}
	return w
}

func (w *anonWriter) setDump(d *dump) { w.dump = d }

func (w *anonWriter) write() error {
	for w.scanner.Scan() {
		line := w.scanner.Bytes()
		if bytes.HasPrefix(line, []byte("INSERT INTO")) {
			line := line[6:]
			scan := newScanner(line)
			insertStmt, err := scanInsert(scan)
			if err != nil {
				log.Fatal(err)
			}
			for i, f := range insertStmt.fields {
				key := fmt.Sprintf("%s.%s", insertStmt.table, f.name())
				if col, ok := w.dump.columnDefinitions[key]; ok && col.category != timeCol {
					insertStmt.fields[i] = w.anon.anonymize(f)
				}
			}
			insertStmt.write(w.out)
		} else {
			if _, err := w.out.Write(w.scanner.Bytes()); err != nil {
				log.Fatal(err)
			}
		}
		w.out.Write([]byte("\n"))
	}
	return w.out.Flush()
}

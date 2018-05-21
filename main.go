package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
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

	scanner := bufio.NewScanner(dumpFile)
	out := bufio.NewWriter(anonFile)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, []byte("INSERT INTO")) {
			line := line[6:]
			scan := newScanner(line)
			insertStmt, err := scanInsert(scan)
			if err != nil {
				log.Fatal(err)
			}
			for i, f := range insertStmt.fields {
				key := fmt.Sprintf("%s.%s", insertStmt.table, f.name())
				if dump.columnDefinitions[key].category != timeCol {
					insertStmt.fields[i] = anonymizer.anonymize(f)
				}
			}
			insertStmt.write(out)
		} else {
			if _, err := out.Write(scanner.Bytes()); err != nil {
				log.Fatal(err)
			}
		}
		out.Write([]byte("\n"))
	}
	out.Flush()
}

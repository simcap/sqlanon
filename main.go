package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"
	"os"
	"path/filepath"
)

var (
	dumpFilepath string
)

func init() {
	flag.StringVar(&dumpFilepath, "f", "", "SQL dump filepath")
}

func main() {
	log.SetFlags(0)
	flag.Parse()

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

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNoopAnonymizerKeepGoldenfileIntact(t *testing.T) {
	goldenFiles, _ := filepath.Glob(filepath.Join(".", "samples", fmt.Sprintf("*%s", ".golden.sql")))
	for _, path := range goldenFiles {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		w := newAnonWriter(f, &out)
		if err := w.write(); err != nil {
			t.Fatal(err)
		}
		golden, err := ioutil.ReadAll(f)
		if err := w.write(); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(golden, out.Bytes()) {
			t.Fatalf("golden file %s failed", path)
		}
	}
}

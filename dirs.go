package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
)

func getFile(newFile string) string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.Buffer{}
	buf.WriteString(dir)
	buf.WriteString(newFile)
	result := buf.String()

	return result

}

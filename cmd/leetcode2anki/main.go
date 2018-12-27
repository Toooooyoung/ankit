package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/invzhi/ankit"
	"github.com/invzhi/ankit/leetcode"
)

var (
	all     bool
	id      int
	path    string
	csvfile string
	dbfile  string
	lang    string
)

func init() {
	flag.BoolVar(&all, "all", false, "whether parse all leetcode question")
	flag.IntVar(&id, "id", 0, "only parse specified leetcode question by id")
	flag.StringVar(&path, "path", ".", "leetcode repo path")
	flag.StringVar(&dbfile, "dbfile", "leetcode.db", "sqlite3 filename in leetcode repo")
	flag.StringVar(&csvfile, "file", "notes.txt", "exported csv filename")
	flag.StringVar(&lang, "lang", "golang", "programming language")
}

func code(path string, _ leetcode.Lang) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func main() {
	flag.Parse()

	if !all && id == 0 {
		log.Fatal("no thing to do")
	}

	f, err := os.Create(csvfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	l := leetcode.New(path, dbfile, leetcode.Lang(lang), code)

	if all {
		if err := ankit.WriteToCSV(f, l); err != nil {
			log.Fatal(err)
		}
	}
}

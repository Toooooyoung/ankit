package main

import (
	"flag"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/invzhi/ankit"
	"github.com/invzhi/ankit/leetcode"
)

var (
	path   string
	dbfile string
	lang   string
	spec   string
)

func init() {
	flag.StringVar(&path, "path", ".", "leetcode repo path")
	flag.StringVar(&dbfile, "db", "leetcode.db", "sqlite3 filename")
	flag.StringVar(&lang, "lang", "golang", "programming language")

	flag.StringVar(&spec, "spec", "", "optional: the path of leetcode question that should be exported only")
}

func question(path string, info os.FileInfo) (leetcode.Key, error) {
	if path != "." && info.IsDir() {
		id, err := strconv.Atoi(path)
		if err != nil {
			return nil, filepath.SkipDir
		}

		return leetcode.KeyID(id), filepath.SkipDir
	}
	return nil, nil
}

func code(path string, _ leetcode.Lang) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Join(path, "code.go"), nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var w strings.Builder
	if err := format.Node(&w, fset, f.Decls); err != nil {
		return "", err
	}

	return w.String(), nil
}

func main() {
	flag.Parse()

	lang := leetcode.Lang(lang)
	if !lang.Valid() {
		log.Fatalf("%s is unsupported on leetcode", lang)
	}

	repo := leetcode.NewRepo(path, dbfile, lang, code, question)

	var d ankit.Deck = repo
	if spec != "" {
		info, err := os.Lstat(spec)
		if err != nil {
			log.Fatal(err)
		}

		// path
		rel, err := filepath.Rel(path, spec)
		if err != nil {
			log.Fatal(err)
		}

		key, err := repo.KeyFn(rel, info)
		if err != nil && err != filepath.SkipDir {
			log.Fatal(err)
		}

		d = ankit.OneNoteDeck(repo.Note(spec, key))
	}

	if err := ankit.Export(os.Stdout, d); err != nil {
		log.Fatal(err)
	}
}

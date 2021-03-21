// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/stfnmllr/go-resp3/internal/ioutildepr"
)

const (
	patchDir = "../cmd/converter"
)

const (
	envGOARCH    = "GOARCH"
	envGOOS      = "GOOS"
	envGOFILE    = "GOFILE"
	envGOLINE    = "GOLINE"
	envGOPACKAGE = "GOPACKAGE"
	envDOLLAR    = "DOLLAR"
)

func goPackage() string { env, _ := os.LookupEnv(envGOPACKAGE); return env }
func goFile() string    { env, _ := os.LookupEnv(envGOFILE); return env }

func sourceDir() string {
	dir, _ := filepath.Split(goFile())
	if dir == "" {
		return "."
	}
	return dir
}

func astOutName(name string) string {
	_, file := filepath.Split(name)
	ext := filepath.Ext(file)
	return filepath.Join(patchDir, strings.TrimSuffix(file, ext)+".ast")
}

var regexpGenerated = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

func isGenerated(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b := false
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		b = regexpGenerated.MatchString(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return b
}

var (
	sourceDirFlag     = flag.String("sourceDir", sourceDir(), "Source directory")
	outResultFlag     = flag.String("outResult", filepath.Join(sourceDir(), "result_gen.go"), "Result output file name")
	outRedisValueFlag = flag.String("outRedisValue", filepath.Join(sourceDir(), "redisvalue_gen.go"), "Redis value output file name")
	pkgNameFlag       = flag.String("package", goPackage(), "package")
)

func main() {

	flag.Parse()
	if *pkgNameFlag == "" {
		log.Fatalf("package missing")
	}

	// command line arguments
	log.Printf("command line arguments: %s", strings.Join(os.Args, ","))

	// go generate flag values
	log.Print("flags:")
	flag.VisitAll(func(flag *flag.Flag) {
		log.Printf(" %s: %s", flag.Name, flag.Value.String())
	})

	// show progress
	var wg sync.WaitGroup
	defer wg.Wait()

	done := progress(&wg)
	wg.Add(1)
	defer close(done)

	// start
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(
		fset,
		*sourceDirFlag,
		func(fi os.FileInfo) bool { return !isGenerated(fi.Name()) },
		parser.ParseComments,
	)
	if err != nil {
		log.Fatal(err)
	}

	pkg, ok := pkgs[*pkgNameFlag]
	if !ok {
		log.Fatalf("package %s not found", *pkgNameFlag)
	}

	if err := writeAstFile(astOutName(*pkgNameFlag), fset, pkg); err != nil {
		log.Fatalf("write file %s error: %s", astOutName(*pkgNameFlag), err)
	}

	a := newAnalyzer()
	a.analyze(pkg)

	g := newGenerator()

	outSrc := g.generateResultFcts(a, *pkgNameFlag)
	if err := writeFile(*outResultFlag, outSrc); err != nil {
		log.Fatalf("write file %s error: %s", *outResultFlag, err)
	}

	outSrc = g.generateRedisValueFcts(a, *pkgNameFlag)
	if err := writeFile(*outRedisValueFlag, outSrc); err != nil {
		log.Fatalf("write file %s error: %s", *outRedisValueFlag, err)
	}

}

func readFile(filename string) ([]byte, error) {
	b, err := ioutildepr.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func writeFile(filename string, b []byte) error {
	return ioutildepr.WriteFile(filename, b, 0644)
}

func writeAstFile(filename string, fset *token.FileSet, node ast.Node) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	ast.Fprint(file, fset, node, nil)
	return nil
}

func progress(wg *sync.WaitGroup) chan struct{} {
	done := make(chan struct{}, 0)

	go func() {
		loop := true
		for loop {
			select {
			case <-done:
				loop = false
			case <-time.After(100 * time.Millisecond):
				print(".")
			}
		}
		println()
		wg.Done()
	}()

	return done
}

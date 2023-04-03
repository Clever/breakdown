package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/Clever/breakdown/gen-go/models"
	"golang.org/x/sync/errgroup"
)

// Module ...
type Module struct {
	Pckgs map[string]*Pkg
}

// Pkg ...
type Pkg struct {
	Name     string
	Version  string
	IsLocal  bool
	Pkgs     []string
	SeenPkgs map[string]bool `json:",omitempty"`
}

var outputFlag = flag.String("output", "/dev/stdout", "output to file location")
var prettyFlag = flag.Bool("pretty", true, "prettify json output")
var versionFlag = flag.Bool("version", false, "print version")
var dirFlag = flag.String("dir", ".", "directory of where to scan dependencies")

var version string

func findFiles(root string) []string {
	var files []string
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case "vendor", "node_modules", ".git":
				return filepath.SkipDir
			}
		}
		if d.Name() == "go.mod" || d.Name() == "package.json" {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("usage: breakdowncli <flags...> <repo_name> <commit_sha>")
	}
	if *versionFlag {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}
	repoName := flag.Arg(0)
	commitSha := flag.Arg(1)
	repoCommit := &models.RepoCommit{
		RepoName:  &repoName,
		CommitSha: &commitSha,
	}

	repoCommit.PackageFiles = make(models.RepoPackageFiles, 0)
	errList := []string{}

	pkgChan := make(chan *models.RepoPackageFile, 10)
	defer close(pkgChan)
	go func() {
		for pkg := range pkgChan {
			repoCommit.PackageFiles = append(repoCommit.PackageFiles, pkg)
			if pkg.Error != "" {
				if pkg.Path != nil {
					errList = append(errList, fmt.Sprintf("(%s): %s", *pkg.Path, pkg.Error))
				} else {
					errList = append(errList, fmt.Sprintf("(NO PATH): %s", pkg.Error))
				}
			}
		}
	}()
	g, _ := errgroup.WithContext(context.Background())
	for _, file := range findFiles(*dirFlag) {
		fileC := file
		switch filepath.Base(file) {
		case "go.mod":
			g.Go(func() error {
				log.Printf("[GOMOD] processing %s", fileC)
				if err := BreakdownGoMod(fileC, pkgChan); err != nil {
					return fmt.Errorf("processing %s: %s", fileC, err)
				}
				return nil
			})
		case "package.json":
			g.Go(func() error {
				log.Printf("[NPM] processing %s", fileC)
				if err := BreakdownNPMPackages(fileC, pkgChan); err != nil {
					return fmt.Errorf("processing %s: %s", fileC, err)
				}
				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		log.Fatalf("%s", err)
	}

	f, err := os.OpenFile(*outputFlag, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalf("opening %q: %s", *outputFlag, err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	if *prettyFlag {
		encoder.SetIndent("", "    ")
	}
	if err := encoder.Encode(repoCommit); err != nil {
		log.Fatalf("enconding repo info: %s", err)
	}

	if len(errList) > 0 {
		log.Printf("found %d error(s):", len(errList))
		for _, e := range errList {
			log.Printf("\t%s", e)
		}
	}

}

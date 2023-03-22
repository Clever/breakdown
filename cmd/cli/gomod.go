package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Clever/breakdown/gen-go/models"
	"github.com/go-openapi/swag"
	"golang.org/x/tools/go/packages"
)

const (
	// pkgLoadMode determines the amount of information retrieved from running packages.Load
	// See https://pkg.go.dev/golang.org/x/tools/go/packages@v0.1.12#LoadMode for more info.
	pkgLoadMode = packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedModule
)

// BreakdownGoMod breaks down package file information
func realBreakdownGoMod(modLoc string, ch chan<- *models.RepoPackageFile) error {
	cfg := &packages.Config{
		Tests:      true,
		Mode:       pkgLoadMode,
		Dir:        filepath.Dir(modLoc),
		BuildFlags: []string{"-mod=readonly", "-tags", "tools"},
	}
	pkgs, err := packages.Load(cfg, "./...", "./tools")
	if err != nil {
		ch <- &models.RepoPackageFile{Path: swag.String(modLoc), Error: err.Error()}
		return nil
	}

	ms, err := getGoModules(pkgs)
	if err != nil {
		ms = &models.RepoPackageFile{Path: &modLoc, Error: err.Error()}
	}
	ms.Path = swag.String(modLoc)
	ch <- ms
	return nil
}

// BreakdownGoMod ...
func BreakdownGoMod(modLoc string, ch chan<- *models.RepoPackageFile) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	proxyChan := make(chan *models.RepoPackageFile)
	defer close(proxyChan)

	start := time.Now()
	go func(l string, c chan<- *models.RepoPackageFile) {
		realBreakdownGoMod(l, c)
	}(modLoc, proxyChan)

	pkgType := "gomod"
	for {
		select {
		case p := <-proxyChan:
			p.Type = &pkgType
			ch <- p
			return nil
		case t := <-ticker.C:
			dur := t.Sub(start)
			log.Printf("processing %q %.2fs", modLoc, dur.Seconds())
			if dur.Seconds() >= 15 {
				ch <- &models.RepoPackageFile{
					Error: fmt.Sprintf("processing file %s timed out at 15 seconds", modLoc),
					Path:  &modLoc,
					Type:  &pkgType,
				}
				return nil
			}
		}
	}
}

func getGoModules(pkgs []*packages.Package) (*models.RepoPackageFile, error) {
	packageFile := &models.RepoPackageFile{Packages: make(map[string]models.RepoPackages)}
	goMod := &Module{Pckgs: make(map[string]*Pkg)}
	mods, err := getGoModulesUsedByPackage(pkgs)
	if err != nil {
		return nil, err
	}
	goMod.Pckgs = mods

	for _, pkg := range pkgs {
		if pkg.Module != nil && pkg.Module.Main {
			packageFile.GoVersion = pkg.Module.GoVersion
			packageFile.Name = pkg.Module.Path
			break
		}
	}

	for modName, modInfo := range goMod.Pckgs {
		deps := []string{}
		for dep := range modInfo.SeenPkgs {
			deps = append(deps, dep)
		}
		sort.Strings(deps)
		packageFile.Packages[modName] = models.RepoPackages{
			Dependencies: deps,
			IsLocal:      modInfo.IsLocal,
			Name:         modInfo.Name,
			Version:      modInfo.Version,
		}
	}

	return packageFile, nil
}

func getGoPkgName(pkg *packages.Package) (string, string) {
	name, v := pkg.Module.Path, pkg.Module.Version
	if pkg.Module.Main {
		v = pkg.Module.GoVersion
	}
	if pkg.Module.Replace != nil && pkg.Module.Replace.Path == pkg.Module.Path {
		v = pkg.Module.Replace.Version
	}
	return name, v
}

func getGoModulesUsedByPackage(queue []*packages.Package) (map[string]*Pkg, error) {
	visitedPackages := make(map[string]bool)
	modules := make(map[string]*Pkg)

	for len(queue) > 0 {
		pkg := queue[0]
		queue = queue[1:]

		if pkg.Module == nil {
			continue
		}

		visitedPackages[pkg.PkgPath] = true
		path, v := getGoPkgName(pkg)
		name := fmt.Sprintf("%s@%s", path, v)

		modPkg, ok := modules[name]
		if !ok {
			modules[name] = &Pkg{
				Name:     path,
				Version:  v,
				SeenPkgs: make(map[string]bool),
				IsLocal:  pkg.Module.Replace != nil && strings.HasPrefix(pkg.Module.Replace.Path, "./"),
			}
			modPkg = modules[name]
		}

		for _, importedPkg := range pkg.Imports {
			if importedPkg.Module == nil {
				continue
			}
			impPath, impVer := getGoPkgName(importedPkg)
			if impPath == path {
				continue
			}
			modPkg.SeenPkgs[fmt.Sprintf("%s@%s", impPath, impVer)] = true
			if visitedPackages[importedPkg.PkgPath] {
				continue
			}
			queue = append(queue, importedPkg)
		}
	}

	return modules, nil
}

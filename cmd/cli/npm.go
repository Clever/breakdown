package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Clever/breakdown/gen-go/models"
	"github.com/go-openapi/swag"
)

// LockfileV is a npm lock file
type LockfileV struct {
	LockfileVersion int64 `json:"lockfileVersion"`
}

// LockfileV2 is an npm lock file that's on version 2 or up
type LockfileV2 struct {
	Packages DependenciesV2 `json:"packages"`
}

// DependencyV2 is a LockfileV2 dependency info
type DependencyV2 struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dev             *bool             `json:"dev"`
	DevDependencies map[string]string `json:"devDependencies"`
	Dependencies    map[string]string `json:"dependencies"`
	Link            *bool             `json:"link"`
	Resolved        string            `json:"resolved"`
}

// DependenciesV2 ...
type DependenciesV2 map[string]DependencyV2

// LockfileV1 ...
type LockfileV1 struct {
	Dependencies DependenciesV1 `json:"dependencies"`
}

// DependencyV1 ...
type DependencyV1 struct {
	Version      string            `json:"version"`
	Dev          bool              `json:"dev"`
	Dependencies DependenciesV1    `json:"dependencies"`
	Requires     map[string]string `json:"requires"`
}

// DependenciesV1 ...
type DependenciesV1 map[string]DependencyV1

// NpmPackageJSON provides a limited view over package.json
type NpmPackageJSON struct {
	Name         string            `json:"name"`
	Dependencies map[string]string `json:"dependencies"`
}

// BreakdownNPMPackages ...
func BreakdownNPMPackages(packageJSONPath string, ch chan<- *models.RepoPackageFile) error {
	packageJSON, err := getPackageJSON(packageJSONPath)
	pkgType := "npm"
	if err != nil {
		ch <- &models.RepoPackageFile{Path: swag.String(packageJSONPath), Error: err.Error(), Type: &pkgType}
		return nil
	}
	if len(packageJSON.Dependencies) < 2 {
		return nil
	}
	packageLockPath := filepath.Join(filepath.Dir(packageJSONPath), "package-lock.json")
	mod, err := parseLockfile(packageJSON.Dependencies, packageLockPath)
	if err != nil {
		ch <- &models.RepoPackageFile{Path: &packageLockPath, Error: err.Error(), Type: &pkgType}
		return nil
	}
	packageFile := &models.RepoPackageFile{Packages: make(map[string]models.RepoPackages), Type: &pkgType}
	for modName, modInfo := range mod.Pckgs {
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

	packageFile.Path = &packageLockPath
	ch <- packageFile
	return nil
}

func getPackageJSON(path string) (*NpmPackageJSON, error) {
	packageJSONBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	packageJSON := NpmPackageJSON{}
	err = json.Unmarshal(packageJSONBytes, &packageJSON)
	if err != nil {
		return nil, err
	}
	return &packageJSON, nil
}

func getLockfileVersion(path string) (int64, []byte, error) {
	lockfile := LockfileV{}
	lockfileBytes, err := os.ReadFile(path)
	if err != nil {
		return -1, nil, err
	}
	err = json.Unmarshal(lockfileBytes, &lockfile)
	if err != nil {
		return -1, nil, err
	}
	return lockfile.LockfileVersion, lockfileBytes, err
}

func parseLockfile(dirDeps map[string]string, path string) (*Module, error) {
	mod := &Module{Pckgs: make(map[string]*Pkg)}
	version, lockfileBytes, err := getLockfileVersion(path)
	if err != nil {
		return nil, err
	}
	switch version {
	case 1:
		lockfileV1 := LockfileV1{}
		err := json.Unmarshal(lockfileBytes, &lockfileV1)
		if err != nil {
			return nil, err
		}
		lockfileV1.Dependencies[""] = DependencyV1{
			Requires: dirDeps,
		}
		return parseLockfileV1(mod, []DependenciesV1{lockfileV1.Dependencies}, "")
	case 2, 3:
		lockfileV2 := LockfileV2{}
		err := json.Unmarshal(lockfileBytes, &lockfileV2)
		if err != nil {
			return nil, err
		}
		return parseLockfileV2(mod, lockfileV2)
	default:
		return nil, fmt.Errorf("unsupported lockfile verison")
	}
}

func getDepVersionV1(name string, depInfo DependencyV1) string {
	version := depInfo.Version
	if strings.HasPrefix(version, "file://") {
		version = ""
	}
	return fmt.Sprintf("%s@%s", name, version)
}

func parseLockfileV1(mod *Module, depLineage []DependenciesV1, parentNameVer string) (*Module, error) {
	for name, dep := range depLineage[len(depLineage)-1] {
		nameVer := getDepVersionV1(name, dep)
		isLocal := strings.HasPrefix(dep.Version, "file://")
		pckg := &Pkg{
			Name:     name,
			Version:  dep.Version,
			SeenPkgs: make(map[string]bool),
			IsLocal:  isLocal,
		}
		mod.Pckgs[nameVer] = pckg

		localDepLineage := append(depLineage, dep.Dependencies)

		for req := range dep.Requires {
			version := ""
			for i := len(localDepLineage) - 1; i >= 0; i-- {
				deps := localDepLineage[i]
				if reqInfo, ok := deps[req]; ok {
					version = reqInfo.Version
					break
				}
			}
			if version == "" {
				return nil, fmt.Errorf("couldnt find req %q part of %q dep in top", req, name)
			}
			reqNameVer := fmt.Sprintf("%s@%s", req, version)
			pckg.SeenPkgs[reqNameVer] = true
		}

		if _, err := parseLockfileV1(mod, localDepLineage, nameVer); err != nil {
			return nil, err
		}
	}
	return mod, nil
}

// Returns an ordered list of top level packages to check for correct version of given
// dependency and current package
// [n_m/.../n_m/.../depenency, n_m/.../dependency, ...]
func genDepNodeModulePath(pkg, dependency string) []string {
	if len(pkg) == 0 {
		return []string{fmt.Sprintf("node_modules/%s", dependency)}
	} else if pkg[len(pkg)-1] != '/' {
		pkg += "/"
	}
	paths := []string{}
	splt := strings.Split(pkg, "node_modules/")
	for i := len(splt); i > 0; i-- {
		paths = append(paths, fmt.Sprintf("%snode_modules/%s", strings.Join(splt[0:i], "node_modules/"), dependency))
	}
	if !strings.HasPrefix(pkg, "node_modules/") {
		paths = append(paths, fmt.Sprintf("node_modules/%s", dependency))
	}
	return paths
}

func parseLockfileV2(mod *Module, lockfile LockfileV2) (*Module, error) {

	// merge deps and devDeps of root package
	for devDep, depInfo := range lockfile.Packages[""].DevDependencies {
		lockfile.Packages[""].Dependencies[devDep] = depInfo
	}

	// pkgName in format: {node_modules/<parent_dep>/}node_modules/<name>
	for pkgName, pkgInfo := range lockfile.Packages {
		name := ""
		if len(pkgName) > 0 {
			parts := strings.SplitAfter(pkgName, "node_modules/")
			name = parts[len(parts)-1]
		}
		isLocal := false
		if pkgInfo.Link != nil {
			isLocal = *pkgInfo.Link
			// resolved will be the name of the directory where the actual package is
			// along w/ name, version and actual dep info
			var ok bool
			pkgInfo, ok = lockfile.Packages[pkgInfo.Resolved]
			if !ok {
				return nil, fmt.Errorf("resolved %q not found for %q", pkgInfo.Resolved, pkgName)
			}
			name = pkgInfo.Name
		}
		nameVer := fmt.Sprintf("%s@%s", name, pkgInfo.Version)
		pkg := &Pkg{
			Name:     name,
			Version:  pkgInfo.Version,
			IsLocal:  isLocal,
			SeenPkgs: make(map[string]bool),
		}
		mod.Pckgs[nameVer] = pkg

		for dep := range pkgInfo.Dependencies {
			name := dep
			version := ""
			for _, check := range genDepNodeModulePath(pkgName, dep) {
				if info, ok := lockfile.Packages[check]; ok {
					if info.Link != nil && *info.Link {
						resolved := lockfile.Packages[info.Resolved]
						name = resolved.Name
						version = resolved.Version
					} else {
						version = info.Version
					}
					break
				}
			}
			if version == "" {
				return nil, fmt.Errorf("couldn't find dep info for %q part of %q package", dep, pkgName)
			}
			pkg.SeenPkgs[fmt.Sprintf("%s@%s", name, version)] = true
		}
	}
	return mod, nil
}

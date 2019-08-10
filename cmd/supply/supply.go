package supply

import (
  "fmt"
  "errors"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

type Stager interface {
	BuildDir() string
	DepDir() string
	DepsIdx() string
	DepsDir() string
  AddBinDependencyLink(string, string) error
}

type Manifest interface {
	AllDependencyVersions(string) []string
	GetEntry(libbuildpack.Dependency) (*libbuildpack.ManifestEntry, error)
	DefaultVersion(string) (libbuildpack.Dependency, error)
}

type Installer interface {
	InstallDependency(libbuildpack.Dependency, string) error
	InstallOnlyVersion(string, string) error
}

type Command interface {
  Run(string) error
}

type Supplier struct {
	Manifest  Manifest
	Installer Installer
	Stager    Stager
	CompileCommand Command
	Log       *libbuildpack.Logger
}

func (s *Supplier) Run() error {
	s.Log.BeginStep("Downloading HAProxy")

  entry, err := s.VersionToInstall()
  s.Log.Info("Using version %s from %s", entry.Dependency.Version, entry.URI)
	if err != nil {
		return err
	}

  dir, err := s.InstallArchive(entry.Dependency)
	if err != nil {
		return err
	}

	haproxyDir := filepath.Join(dir, fmt.Sprintf("haproxy-%s", entry.Dependency.Version))
	if err := s.CompileAndLink(haproxyDir); err != nil {
		return err
	}

	return nil
}

func (s *Supplier) VersionToInstall() (*libbuildpack.ManifestEntry, error) {
	versions := s.Manifest.AllDependencyVersions("haproxy")
	if len(versions) < 1 {
		return nil, errors.New("Unable to find a version of haproxy to install")
	}
	dep := libbuildpack.Dependency{Name: "haproxy", Version: versions[len(versions)-1]}
	return s.Manifest.GetEntry(dep)
}

func (s *Supplier) InstallArchive(dep libbuildpack.Dependency) (string, error) {
	dir := filepath.Join(s.Stager.DepDir(), "haproxy")
	return dir, s.Installer.InstallDependency(dep, dir)
}

func (s *Supplier) CompileAndLink(dir string) error {
	if err := s.CompileCommand.Run(dir); err != nil {
    return err
  }
  return s.Stager.AddBinDependencyLink(filepath.Join(dir, "haproxy"), "haproxy")
}

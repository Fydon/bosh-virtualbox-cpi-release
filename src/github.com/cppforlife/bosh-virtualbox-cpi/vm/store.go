package vm

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cppforlife/bosh-virtualbox-cpi/driver"
)

type Store struct {
	path   string
	runner driver.Runner
}

func NewStore(path string, runner driver.Runner) Store {
	return Store{path, runner}
}

func (m Store) List() ([]string, error) {
	var out string
	if runtime.GOOS == "windows" {
		_, _, err := m.runner.Execute(fmt.Sprintf("if( -Not ( Test-Path \"%s\" ) ) { New-Item -ItemType Directory \"%s\" }", m.path, m.path))
		if err != nil {
			return nil, err
		}

		out, _, err = m.runner.Execute(fmt.Sprintf("( Get-ChildItem \"%s\" ).Name", m.path))
		if err != nil {
			return nil, err
		}
	} else {
		_, _, err := m.runner.Execute("mkdir", "-p", m.path)
		if err != nil {
			return nil, err
		}

		out, _, err = m.runner.Execute("ls", "-1", m.path)
		if err != nil {
			return nil, err
		}
	}

	return strings.Split(out, "\n"), nil
}

func (m Store) Path(key string) string {
	return filepath.Join(m.path, key)
}

func (m Store) Put(key string, contents []byte) error {
	if runtime.GOOS == "windows" {
		_, _, err := m.runner.Execute(fmt.Sprintf("if( -Not ( Test-Path \"%s\" ) ) { New-Item -ItemType Directory \"%s\" }", m.path, m.path))
		if err != nil {
			return err
		}
	} else {
		_, _, err := m.runner.Execute("mkdir", "-p", m.path)
		if err != nil {
			return err
		}
	}

	return m.runner.Put(filepath.Join(m.path, key), contents)
}

func (m Store) Get(key string) ([]byte, error) {
	return m.runner.Get(filepath.Join(m.path, key))
}

func (m Store) DeleteOne(key string) error {
	_, _, err := m.runner.Execute("rm", "-rf", filepath.Join(m.path, key))
	return err
}

func (m Store) Delete() error {
	_, _, err := m.runner.Execute("rm", "-rf", m.path)
	return err
}

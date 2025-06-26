package core

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"plugin"

	"go.uber.org/dig"
)

type ModuleLoader struct {
	di *dig.Container
}

func NewModuleLoader(di *dig.Container) *ModuleLoader {
	return &ModuleLoader{
		di: di,
	}
}

func (mma *ModuleLoader) Load() error {

	path := os.Getenv("MODULES_PATH")

	fi, err := os.Stat(path)

	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New("expected directory path for modules")
	}

	err = filepath.WalkDir(path, func(modulePath string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(modulePath) != ".so" {
			return errors.New("expected plugin file")
		}

		proxy, err := plugin.Open(
			modulePath,
		)

		if err != nil {
			return err
		}

		fn, err := proxy.Lookup("Wire")

		if err != nil {
			return err
		}

		wire, isWireFunction := fn.(func(di *dig.Container))

		if !isWireFunction {
			return errors.New("wire function expected")
		}

		fmt.Printf("Module loaded: %s\n", modulePath)

		wire(mma.di)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

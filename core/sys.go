package core

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/dig"
)

type ModuleInfo struct {
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Website     string `json:"website"`
}

type ModuleLoader struct {
	di       *dig.Container
	registry []*ModuleInfo
}

func NewModuleLoader(di *dig.Container) *ModuleLoader {
	return &ModuleLoader{
		di:       di,
		registry: make([]*ModuleInfo, 0),
	}
}

func (loader *ModuleLoader) Load() error {

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

		moduleName := filepath.Base(filepath.Dir(modulePath))

		proxy, err := plugin.Open(
			modulePath,
		)

		if err != nil {
			return err
		}

		ui, err := proxy.Lookup("UI")

		if err == nil && ui != nil {
			loader.di.Invoke(func(router chi.Router) {

				content, err := fs.Sub(ui.(*embed.FS), "ui/dist")

				if err != nil {
					log.Fatal(err)
				}

				router.Handle(
					fmt.Sprintf("/ui/%s/*", moduleName),
					http.StripPrefix(fmt.Sprintf("/ui/%s/", moduleName), http.FileServer(http.FS(content))),
				)

			})
		}

		fn, err := proxy.Lookup("Wire")

		if err != nil {
			return err
		}

		wire, isWireFunction := fn.(func(di *dig.Container) *ModuleInfo)

		if !isWireFunction {
			return errors.New("wire function expected")
		}

		fmt.Printf("Module loaded: %s\n", modulePath)

		loader.registry = append(loader.registry, wire(loader.di))

		return nil
	})

	if err != nil {
		return err
	}

	loader.di.Invoke(func(router chi.Router) {
		router.Get("/sys/modules", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, loader.registry)
		})
	})

	return nil
}

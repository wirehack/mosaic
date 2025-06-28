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
	di           *dig.Container
	registry     []*ModuleInfo
	onUiRegister func(mi *ModuleInfo, ui *embed.FS)
}

func NewModuleLoader(di *dig.Container, onUiRegister func(mi *ModuleInfo, ui *embed.FS)) *ModuleLoader {
	return &ModuleLoader{
		di:           di,
		registry:     make([]*ModuleInfo, 0),
		onUiRegister: onUiRegister,
	}
}

func (loader *ModuleLoader) Registry() []*ModuleInfo {
	return loader.registry
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

		wire, isWireFunction := fn.(func(di *dig.Container) *ModuleInfo)

		if !isWireFunction {
			return errors.New("wire function expected")
		}

		current := wire(loader.di)

		loader.registry = append(loader.registry, current)

		ui, err := proxy.Lookup("UI")

		if err == nil && ui != nil {
			loader.onUiRegister(current, ui.(*embed.FS))
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func RegisterModules(di *dig.Container) error {

	loader := NewModuleLoader(di, func(mi *ModuleInfo, ui *embed.FS) {

		di.Invoke(func(router chi.Router) {

			content, err := fs.Sub(ui, "ui/dist")

			if err != nil {
				log.Fatal(err)
			}

			router.Handle(
				fmt.Sprintf("/ui/%s/*", mi.Slug),
				http.StripPrefix(fmt.Sprintf("/ui/%s/", mi.Slug), http.FileServer(http.FS(content))),
			)

		})

	})

	di.Invoke(func(router chi.Router) {
		router.Get("/sys/modules", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, loader.Registry())
		})
	})

	if err := loader.Load(); err != nil {
		return err
	}

	return nil
}

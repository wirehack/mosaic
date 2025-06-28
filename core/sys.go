package core

import (
	"c"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/dig"
)

type ModuleInfo struct {
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
	Author      string `json:"author,omitempty"`
	Website     string `json:"website,omitempty"`
}

type ModuleLoader struct {
	di           *dig.Scope
	registry     []*ModuleInfo
	onUiRegister func(mi *ModuleInfo, ui *embed.FS)
}

func NewModuleLoader(di *dig.Scope, onUiRegister func(mi *ModuleInfo, ui *embed.FS)) *ModuleLoader {
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

		wire, isWireFunction := fn.(func(di *dig.Scope) *ModuleInfo)

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

func RegisterModules(di *dig.Scope) error {

	loader := NewModuleLoader(di, func(mi *ModuleInfo, ui *embed.FS) {

		di.Invoke(func(router chi.Router) {

			content, err := fs.Sub(ui, "ui/dist")

			if err != nil {
				log.Fatal(err)
			}

			uiPath := fmt.Sprintf("/%s/ui", mi.Slug)

			// Serve static files
			router.Method(
				http.MethodGet,
				uiPath+"/*",
				http.StripPrefix("/v1/modules"+uiPath+"/", http.FileServer(http.FS(content))),
			)

			// Serve index.html fallback for SPA routes
			router.Get(uiPath, func(w http.ResponseWriter, r *http.Request) {
				f, err := content.Open("index.html")
				if err != nil {
					http.Error(w, "index.html not found", http.StatusInternalServerError)
					return
				}
				defer f.Close()
				http.ServeContent(w, r, "index.html", time.Now(), f.(io.ReadSeeker))
			})

			router.Get("/test/info", func(w http.ResponseWriter, r *http.Request) {
				render.JSON(w, r, mi)
			})

		})

		c.Log().Infof("Module %s loaded", mi.Slug)
	})

	di.Invoke(func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, loader.Registry())
		})
	})

	if err := loader.Load(); err != nil {
		return err
	}

	return nil
}

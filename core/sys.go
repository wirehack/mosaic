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
	"net/url"
	"os"
	"path/filepath"
	"plugin"
	"strings"

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

			renderIndex := func(w http.ResponseWriter, r *http.Request) {
				f, err := content.Open("index.html")
				if err != nil {
					http.NotFound(w, r)
					return
				}
				if _, err := io.Copy(w, f); err != nil {
					panic(err)
				}
			}

			router.Get(uiPath, renderIndex)
			router.Get(uiPath+"/", renderIndex)

			router.Get(uiPath+"/*", func(w http.ResponseWriter, r *http.Request) {

				parts := strings.Split(r.URL.Path, uiPath+"/")

				if len(parts) != 2 {
					fmt.Println(r.URL.Path)
					http.NotFound(w, r)
					return
				}

				var file = parts[1]

				if file == "" {
					file = "index.html"
				}

				r2 := new(http.Request)
				*r2 = *r
				r2.URL = new(url.URL)
				*r2.URL = *r.URL
				r2.URL.Path = file
				r2.URL.RawPath = file

				h := http.FileServer(http.FS(content))
				h.ServeHTTP(w, r2)
			})

			router.Get(fmt.Sprintf("/%s/info", mi.Slug), func(w http.ResponseWriter, r *http.Request) {
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

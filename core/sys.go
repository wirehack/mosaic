package core

import (
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

type ModuleMetaAware interface {
	Meta() *ModuleInfo
}

type ModuleUIAware interface {
	UI() *embed.FS
}

type ModuleLoader struct {
	di         *dig.Scope
	registry   []any
	onRegister func(proxy any)
}

func NewModuleLoader(di *dig.Scope, onRegister func(proxy any)) *ModuleLoader {
	return &ModuleLoader{
		di:         di,
		registry:   make([]any, 0),
		onRegister: onRegister,
	}
}

func (loader *ModuleLoader) Find(slug string) any {
	for _, item := range loader.registry {
		if proxy, b := item.(ModuleMetaAware); b {
			return proxy
		}
	}
	return nil
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

		wire, isWireFunction := fn.(func(di *dig.Scope) any)

		if !isWireFunction {
			return errors.New("wire function expected")
		}

		current := wire(loader.di)

		loader.registry = append(loader.registry, current)

		loader.onRegister(current)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

type ModuleProxy struct {
	registry []any
}

func NewModuleProxy() *ModuleProxy {
	return &ModuleProxy{registry: make([]any, 0)}
}

func (mp ModuleProxy) Get(slug string) (ret any, exists bool) {

	for _, proxy := range mp.registry {
		if metaAware, isMetaAware := proxy.(ModuleMetaAware); isMetaAware {
			if metaAware.Meta().Slug == slug {
				return proxy, true
			}
		}
	}

	return ret, false
}

func (mp *ModuleProxy) Register(proxy any) {
	mp.registry = append(mp.registry, proxy)
}

func RegisterModules(di *dig.Scope) (proxy *ModuleProxy, err error) {

	mp := NewModuleProxy()

	di.Provide(func() *ModuleProxy { return mp }, dig.Export(true))

	loader := NewModuleLoader(di, func(proxy any) {

		mp.Register(proxy)

		metaAware, isMetaAware := proxy.(ModuleMetaAware)
		uiAware, isUIAware := proxy.(ModuleUIAware)

		if !isUIAware || !isMetaAware {
			return
		}

		di.Invoke(func(router chi.Router) {

			ui := uiAware.UI()

			content, err := fs.Sub(ui, "ui/dist")

			if err != nil {
				log.Fatal(err)
			}

			uiPath := fmt.Sprintf("/%s/ui", metaAware.Meta().Slug)

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

			router.Get(fmt.Sprintf("/%s/info", metaAware.Meta().Slug), func(w http.ResponseWriter, r *http.Request) {
				render.JSON(w, r, metaAware.Meta())
			})

		})
	})

	if err := loader.Load(); err != nil {
		return nil, err
	}

	return proxy, nil
}

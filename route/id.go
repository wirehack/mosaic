package route

import (
	"fmt"
	"mosaic/core"
	"mosaic/types"
	"net/http"
)

func CreateID(w http.ResponseWriter, r *http.Request) {
	if proxy, exists := core.Module[types.TestModuleProxy]("test"); exists {
		fmt.Println(proxy.Sum(2, 2))
	}
}

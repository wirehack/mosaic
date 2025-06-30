package route

import (
	"c"
	"fmt"
	"mosaic/types"
	"net/http"
)

func CreateID(w http.ResponseWriter, r *http.Request) {
	if proxy, exists := c.Module[types.TestModuleProxy]("test"); exists {
		fmt.Println(proxy.Sum(2, 2))
	}
}

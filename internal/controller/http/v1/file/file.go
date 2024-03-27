package file

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/pkg/config"
	"github.com/restaurant/internal/service/hashing"
	"net/http"
	"strings"
	"time"
)

type Controller struct {
	*web.App
}

type onlyFilesFS struct {
	http.FileSystem
}

func NewController(app *web.App) *Controller {
	return &Controller{
		app,
	}
}

func (con *Controller) File(c *gin.Context) {
	fs := gin.Dir("./", false)
	if _, noListing := fs.(*onlyFilesFS); noListing {
		fmt.Println("Error: no listing filesystem here")
		c.Writer.WriteHeader(http.StatusNotFound)
	}

	file := c.Param("filepath")
	OpenH := hashing.ParseHash(file)
	list := strings.Split(OpenH, " ")
	if len(list) == 3 {
		linkTime, err := time.Parse("02.01.2006 15:04:05 ", list[1]+" "+list[2])
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]any{
				"error":  "incorrect link",
				"status": false,
			})
			return
		}
		if linkTime.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, map[string]any{
				"error":  "expired link",
				"status": false,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error":  "incorrect link",
			"status": false,
		})
		return
	}
	cfg := config.NewConfig()

	http.ServeFile(c.Writer, c.Request, cfg.FileDirectory+list[0])
}

package api

import (
	"vue-argon-design-system-master/server/api/middleware"
	"vue-argon-design-system-master/server/api/module/insertUserInformation"

	"github.com/labstack/echo"
)

//Init api binding
func Init(e *echo.Echo) {
	o := e.Group("/o")
	r := e.Group("/r")

	c := r.Group("/c")

	middleware.Init(e, o, r, c)
	insertUserInformation.Init(e, o, r, c)

}

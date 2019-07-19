package middleware

import (
	"net/http"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/confighelper"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/authhelper"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Init(e *echo.Echo, o, r, c *echo.Group) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*", "http://localhost:3030"},
		AllowMethods:  []string{echo.GET, echo.OPTIONS, echo.PUT, echo.POST, echo.DELETE},
		ExposeHeaders: []string{"authorization"},
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${status} ${method} ${uri} \n",
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
}

// This middleware will be called for every restricted URL
func JwtMiddleware() echo.MiddlewareFunc {
	// logginghelper.LogInfo("JwtMiddleware called.............. ")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenFromRequest := c.Request().Header.Get("authorization")
			if tokenFromRequest == "" {
				// logginghelper.LogError("error occured while fetching token from request header")
				return echo.ErrUnauthorized
			}
			_, terr := authhelper.DecodeToken(tokenFromRequest, confighelper.GetConfig("JWTSecretKey"))
			if nil != terr {
				// logginghelper.LogError("error while decoding token ", terr)
				return c.JSON(http.StatusUnauthorized, "ERR_UNAUTHORIZED_USER")
			}
			return next(c)
		}
	}
}

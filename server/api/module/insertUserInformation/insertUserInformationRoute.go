package insertUserInformation

import (
	"net/http"
	"vue-argon-design-system-master/server/api/constants"
	"vue-argon-design-system-master/server/api/model"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"
	"github.com/labstack/echo"
)

func Init(o, r, c *echo.Group) {
	o.POST("/insertUserInformation", insertUserInformationRoute)

}

func insertUserInformationRoute(c echo.Context) error {
	logginghelper.LogDebug("IN: insertUserInformationRoute")
	userInfo := model.UserInformation{}
	err := c.Bind(&userInfo)
	if err != nil {
		logginghelper.LogError("ERR_BINDING_DATA:" + err.Error())
		return c.JSON(http.StatusExpectationFailed, err.Error())
	}
	value, err := insertUserInformationService(userInfo)

	if err != nil {
		if err.Error() == constants.ERRORCODE_USERNAME_ALREADY_EXISTS {
			return c.JSON(http.StatusAlreadyReported, constants.ERRORCODE_USERNAME_ALREADY_EXISTS)
		}
		logginghelper.LogError("ERROR_INSERTING_DATA", err.Error())
		return c.JSON(http.StatusExpectationFailed, value)
	}
	logginghelper.LogDebug("OUT : validateUserDetails")
	return c.JSON(http.StatusOK, value)
}

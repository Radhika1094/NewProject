package insertUserInformation

import (
	"vue-argon-design-system-master/server/api/model"

	"golang.org/x/crypto/bcrypt"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"
)

func insertUserInformationService(userInfo model.UserInformation) (bool, error) {
	logginghelper.LogDebug("IN: insertUserInformationService")
	value, err := insertUserInformationDAO(userInfo)
	if err != nil {
		logginghelper.LogError("ERR_INSERTING_DATA:" + err.Error())
		return value, err
	}
	if value {
		return insertLoginCredentialsService(userInfo.UserName, userInfo.Password)
	}
	return value, err
}

func insertLoginCredentialsService(userName string, password string) (bool, error) {
	logginghelper.LogDebug("IN: insertLoginCredentialsService ")
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	loginObj := model.UserCredential{}
	loginObj.Password = string(bytes)
	loginObj.UserName = userName
	value, err := insertLoginCredentialsDAO(loginObj.UserName, loginObj.Password)
	return value, err

}

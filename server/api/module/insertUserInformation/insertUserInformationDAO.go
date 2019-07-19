package insertUserInformation

import (
	"errors"
	"fmt"
	"vue-argon-design-system-master/server/api/constants"
	"vue-argon-design-system-master/server/api/model"

	"gopkg.in/mgo.v2/bson"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/confighelper"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/dalhelper"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"
)

func insertUserInformationDAO(userInfo model.UserInformation) (bool, error) {
	logginghelper.LogDebug("In: insertUserInformationDAO")

	userInfoObj := model.UserInformation{}
	session, dberr := dalhelper.GetMongoConnection()

	if dberr != nil {
		logginghelper.LogError("ERROR_IN_CONNECTION ", dberr)
		return false, dberr
	}
	logginghelper.LogInfo("Login Details:", userInfoObj)
	cd := session.DB(confighelper.GetConfig("DBNAME")).C(model.COLLECTION_INSERT_USERINFO)

	err := cd.Find(bson.M{"userName": userInfo.UserName}).One(&userInfoObj)

	if err != nil {
		fmt.Println("err.....", err.Error())
		if err.Error() == "not found" {
			findErr := cd.Insert(&userInfo)
			if nil != findErr {

				logginghelper.LogError("ERROR_WHILE_INSERTING", findErr.Error())
				return false, findErr
			}
		} else {
			logginghelper.LogInfo("IN FINDERR:")

		}
	} else {
		logginghelper.LogInfo("IN ERR:")
		return false, errors.New(constants.ERRORCODE_USERNAME_ALREADY_EXISTS)
	}
	logginghelper.LogDebug("Out: insertUserInformationDAO")
	return true, nil
}
func insertLoginCredentialsDAO(username string, password string) (bool, error) {
	logginghelper.LogDebug("In: insertLoginCredentialsDAO")

	session, dberr := dalhelper.GetMongoConnection()

	if dberr != nil {
		logginghelper.LogError("ERROR_WHILE_GETTING_CONNECTION ", dberr)
		return false, dberr
	}

	cd := session.DB(confighelper.GetConfig("DBNAME")).C(model.COLLECTION_INSERT_USERCREDENTIAL)
	findErr := cd.Insert(bson.M{"userName": username, "password": password})

	if nil != findErr {
		logginghelper.LogError("ERROR_WHILE_INSERTING", findErr.Error())
		return false, findErr
	}
	logginghelper.LogDebug("Out: insertLoginCredentialsDAO")
	return true, nil
}

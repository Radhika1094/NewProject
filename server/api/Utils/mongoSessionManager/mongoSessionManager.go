package mongoSessionManager

import (
	"MahaVan/MahaVanServer/api/model"
	"errors"
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/confighelper"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"
	"gopkg.in/mgo.v2/bson"
)

var (
	dbName = confighelper.GetConfig("SSODBNAME")
)

const (
	SESSION_ID = "SESSION_ID"
	EXPIRE_ON  = "EXPIRE_ON"
)

var instance *mgo.Session
var sessionError error
var once sync.Once

func GetSSOStoreConnection() (*mgo.Session, error) {
	once.Do(func() {
		Host := []string{
			confighelper.GetConfig("SSODSN"),
		}
		var (
			Username = ""
			Password = ""
			Database = confighelper.GetConfig("SSODBNAME")
		)

		session, err := mgo.DialWithInfo(&mgo.DialInfo{
			Addrs:    Host,
			Username: Username,
			Password: Password,
			Database: Database,
		})
		if err != nil {
			sessionError = err
		}
		//defer session.Close()

		instance = session

	})
	return instance, sessionError
}

func Set(key string, value string, idealDuration string) error {
	logginghelper.LogInfo("IN: Set session")
	session, err := GetSSOStoreConnection()
	if err != nil {
		logginghelper.LogError("UNABLE TO CONNECT TO DB Error : ", err)
		return err
	}
	sessionEntry := model.SessionEntry{}
	// if ideal duration is not set then no need to set expiredAt
	// because the document should not be deleted
	if len(idealDuration) > 0 {
		ttl, perr := time.ParseDuration(idealDuration)
		if nil != perr {
			logginghelper.LogError("ERR_PARSING_SESSION", perr)
			return errors.New("ERR_PARSING_SESSION")
		}
		sessionEntry.ExpireOn = time.Now().Add(ttl)
	}
	sessionEntry.SessionId = key
	sessionEntry.Token = value
	sessionStoreCollection := session.DB(dbName).C(model.COLLECTION_SESSION_STORE)
	insertErr := sessionStoreCollection.Insert(sessionEntry)
	if nil != insertErr {
		logginghelper.LogError(insertErr)
		return insertErr
	}
	logginghelper.LogInfo("OUT: Set session")
	return nil
}

func getActiveSessionCount(pattern string) (int, error) {
	logginghelper.LogInfo("IN: getActiveSessionCount")

	count := 0
	session, err := GetSSOStoreConnection()
	if err != nil {
		logginghelper.LogError("UNABLE TO CONNECT TO DB Error : ", err)
		return count, err
	}
	query := bson.M{SESSION_ID: bson.M{"$regex": bson.RegEx{"^" + pattern, "i"}}}
	sessionStoreCollection := session.DB(dbName).C(model.COLLECTION_SESSION_STORE)
	count, err = sessionStoreCollection.Find(query).Count()
	if nil != err {
		if err.Error() == "not found" {
			return count, nil
		}
		logginghelper.LogError(err)
		return count, err
	}
	logginghelper.LogInfo("OUT: getActiveSessionCount")
	return count, nil
}

func IsSessionLimitReached(loginId string, maxSessionAllowed int) (bool, error) {
	logginghelper.LogInfo("IN: IsSessionLimitReached")
	activeSessions, err := getActiveSessionCount(loginId)

	if nil != err {
		logginghelper.LogError("IsSessionLimitReached", err)
		return false, err
	}
	logginghelper.LogDebug("activeSessions", activeSessions)
	if activeSessions >= maxSessionAllowed {
		return true, nil
	}

	logginghelper.LogInfo("OUT: IsSessionLimitReached")
	return false, nil
}

func Get(key string) (string, error) {
	logginghelper.LogInfo("IN: Get session")
	token := ""
	session, err := GetSSOStoreConnection()
	if err != nil {
		logginghelper.LogError("UNABLE TO CONNECT TO DB Error : ", err)
		return token, err
	}
	logginghelper.LogDebug("key", key)
	sessionEntry := model.SessionEntry{}
	sessionStoreCollection := session.DB(dbName).C(model.COLLECTION_SESSION_STORE)
	findErr := sessionStoreCollection.Find(bson.M{SESSION_ID: key}).One(&sessionEntry)
	if nil != findErr {
		logginghelper.LogError(findErr)
		if findErr.Error() == "not found" {
			return token, errors.New("ERR_SESSION_NOT_FOUND")
		}
		return token, findErr
	}
	logginghelper.LogInfo("OUT: Get session")
	return sessionEntry.Token, nil
}

func Remove(keys ...string) error {
	logginghelper.LogInfo("IN: Delete session")
	session, err := GetSSOStoreConnection()
	if err != nil {
		logginghelper.LogError("UNABLE TO CONNECT TO DB Error : ", err)
		return err
	}
	sessionStoreCollection := session.DB(dbName).C(model.COLLECTION_SESSION_STORE)
	_, delErr := sessionStoreCollection.RemoveAll(bson.M{SESSION_ID: bson.M{"$in": keys}})
	if nil != delErr {
		if delErr.Error() == "not found" {
			logginghelper.LogError(delErr)
			return errors.New("ERR_SESSION_NOT_FOUND")
		}
		logginghelper.LogError(delErr)
		return delErr
	}
	logginghelper.LogInfo("OUT: Delete session")
	return nil
}

func SlideSession(key, idealDuration string) error {
	logginghelper.LogInfo("IN:: SlideSession")
	session, err := GetSSOStoreConnection()
	if err != nil {
		logginghelper.LogError("UNABLE TO CONNECT TO DB Error : ", err)
		return err
	}
	if len(idealDuration) == 0 {
		// if ideal duration is not set then no need to slide session
		// because it is not going to get deleted
		return nil
	}
	ttl, perr := time.ParseDuration(idealDuration)
	if nil != perr {
		logginghelper.LogError("ERR_SLIDING_SESSION: UNABLE TO PARSE IDEAL DURATION STRING: ", idealDuration, ": ", perr)
		return errors.New("ERR_SLIDING_SESSION")
	}
	selector := bson.M{SESSION_ID: key}
	updator := bson.M{"$set": bson.M{EXPIRE_ON: time.Now().Add(ttl)}}
	sessionStoreCollection := session.DB(dbName).C(model.COLLECTION_SESSION_STORE)
	updateErr := sessionStoreCollection.Update(selector, updator)
	if nil != updateErr {
		if updateErr.Error() == "not found" {
			logginghelper.LogError(updateErr)
			return errors.New("ERR_SESSION_NOT_FOUND")
		}
		logginghelper.LogError(updateErr)
		return updateErr
	}
	logginghelper.LogInfo("OUT:: SlideSession")
	return nil
}

func KillSession(username string) error {
	logginghelper.LogInfo("IN: KillSession")
	session, err := GetSSOStoreConnection()
	if err != nil {
		logginghelper.LogError("UNABLE TO CONNECT TO DB Error : ", err)
		return err
	}
	change := mgo.Change{}
	change.Remove = true
	sessionEntry := model.SessionEntry{}
	logginghelper.LogDebug("constants.SESSION_KEY_PREFIX + username", bson.RegEx{username, "i"})
	sessionStoreCollection := session.DB(dbName).C(model.COLLECTION_SESSION_STORE)
	_, findErr := sessionStoreCollection.Find(bson.M{SESSION_ID: bson.M{"$regex": bson.RegEx{username, "i"}}}).Sort(EXPIRE_ON).Apply(change, &sessionEntry)
	if nil != findErr {
		if findErr.Error() == "not found" {
			logginghelper.LogError(findErr)
			return errors.New("ERR_SESSION_NOT_FOUND")
		}
		logginghelper.LogError(findErr)
		return findErr
	}
	logginghelper.LogInfo("OUT: KillSession")
	return nil
}

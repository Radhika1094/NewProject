package jwtService

import (
	"MahaVan/MahaVanServer/api/constants"
	"MahaVan/MahaVanServer/api/model"
	"errors"
	"time"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/confighelper"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/authhelper"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//GENERATE JWT TOKEN
func GenerateJwtToken(login model.CitizenLoginDetails, sessionId string) (string, error) {
	logginghelper.LogDebug("Inside: generateJwtToken")
	standardClaim := jwt.StandardClaims{}
	expAt := time.Now().Add(constants.EXPIRATION_DURATION).Unix()
	standardClaim.ExpiresAt = expAt
	// expAt := time.Now().Add(constants.JWT_TOKEN_EXPAT).Unix()
	claims := model.JwtCustomClaims{
		login.LoginID,
		login.LoginID,
		sessionId,
		standardClaim,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// jwtSecretKey := confighelper.GetConfig("JWTSecretKey")
	tokenString, err := token.SignedString([]byte(constants.JWT_SECRETKEY))
	if nil != err {
		logginghelper.LogError("GenerateToken SignedString() Error: ", err)
		return "", nil
	}
	tokenString = "Bearer " + tokenString
	logginghelper.LogDebug("TOKEN: ", tokenString)
	return tokenString, nil
}

// GetDecodedLoginFromToken ptovide decoded login object from JWT Token
func GetDecodedLoginFromToken(c echo.Context) (string, error) {
	logginghelper.LogDebug("IN : GetDecodedLoginFromToken")
	login := model.CitizenLoginDetails{}
	decodedToken, err := authhelper.DecodeToken(c.Request().Header.Get("Authorization"), constants.JWT_SECRETKEY)
	if err != nil {
		logginghelper.LogError("GetDecodedLoginFromToken DecodeToken() Error: ", err)
		return "", errors.New("ERA_ERRORCODE_JWTUTILS_FAILED_TO_DECODE_TOKEN")
	}
	// login ID is the compulsary field, so haven't added check for nil
	if decodedToken["username"] == nil || decodedToken["username"] == "" {
		return "", errors.New("ERA_ERRORCODE_JWTUTILS_LOGIN_ID_NOT_FOUND")
	}
	login.LoginID = decodedToken["username"].(string)

	logginghelper.LogDebug("OUT : sGetDecodedLoginFromToken")
	return login.LoginID, nil
}

// GetSessionIdFromToken extract sessionId from JWT Token
func GetSessionIdFromToken(c echo.Context) (string, error) {
	logginghelper.LogDebug("IN : GetDecodedLoginFromToken")
	decodedToken, err := authhelper.DecodeToken(c.Request().Header.Get("Authorization"), constants.JWT_SECRETKEY)
	if err != nil {
		logginghelper.LogError("GetDecodedLoginFromToken DecodeToken() Error: ", err)
		return "", errors.New("ERA_ERRORCODE_JWTUTILS_FAILED_TO_DECODE_TOKEN")
	}
	// login ID is the compulsary field, so haven't added check for nil
	if decodedToken["sessionId"] == nil || decodedToken["sessionId"] == "" {
		return "", errors.New("ERA_ERRORCODE_JWT_UTILS_SESSION_ID_NOT_FOUND")
	}
	sessionId := decodedToken["sessionId"].(string)

	logginghelper.LogDebug("OUT : GetDecodedLoginFromToken")
	return sessionId, nil
}

// GetTokenRestrictTokenService create jwt token which is not used for login but for other purpose
func GetTokenRestrictTokenService(username, clientID string) (string, error) {
	standardClaim := jwt.StandardClaims{}
	expAt := time.Now().Add(constants.EXPIRATION_DURATION).Unix()
	standardClaim.ExpiresAt = expAt
	// expAt := time.Now().Add(constants.JWT_TOKEN_EXPAT).Unix()
	claims := model.JwtCustomClaims{
		username,
		username,
		username,
		standardClaim,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	restrictTokenKey := confighelper.GetConfig("restrictTokenKey")
	tokenString, err := token.SignedString([]byte(restrictTokenKey))
	if nil != err {
		logginghelper.LogError("GetTokenRestrictedurlService SignedString() Error: ", err)
		return "", nil
	}
	tokenString = "Bearer " + tokenString
	return tokenString, nil
}

package jwtService

import (
	"vue-argon-design-system-master/server/api/constants"
	"vue-argon-design-system-master/server/api/model"
	"errors"
	"fmt"
	"net/http"
	"time"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/confighelper"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/authhelper"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//GENERATE JWT TOKEN
func GenerateJwtToken(login model.LoginDetails, sessionId string) (string, error) {
	logginghelper.LogDebug("Inside: generateJwtToken")
	standardClaim := jwt.StandardClaims{}
	expAt := time.Now().Add(constants.EXPIRATION_DURATION).Unix()
	standardClaim.ExpiresAt = expAt
	// expAt := time.Now().Add(constants.JWT_TOKEN_EXPAT).Unix()
	claims := model.JwtCustomClaims{
		login.UserName,
		// sessionId,
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

func GenerateJwtTokenUser(login model.UserCredentials) (string, error) {
	logginghelper.LogDebug("Inside: generateJwtToken")
	standardClaim := jwt.StandardClaims{}
	expAt := time.Now().Add(constants.EXPIRATION_DURATION).Unix()
	standardClaim.ExpiresAt = expAt
	// expAt = time.Now().Add(constants.JWT_TOKEN_EXPAT).Unix()
	claims := model.JwtCustomClaims{
		login.Username,
		// login.Password,
		// sessionId,
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
	login := model.LoginDetails{}
	decodedToken, err := authhelper.DecodeToken(c.Request().Header.Get("Authorization"), constants.JWT_SECRETKEY)
	fmt.Println("decode token", decodedToken)
	fmt.Println("request header", c.Request().Header)
	if err != nil {
		logginghelper.LogError("GetDecodedLoginFromToken DecodeToken() Error: ", err)
		return "", errors.New("ERA_ERRORCODE_JWTUTILS_FAILED_TO_DECODE_TOKEN")
	}
	// login ID is the compulsary field, so haven't added check for nil
	if decodedToken["username"] == nil || decodedToken["username"] == "" {
		return "", errors.New("ERA_ERRORCODE_JWTUTILS_LOGIN_ID_NOT_FOUND")
	}
	login.UserName = decodedToken["username"].(string)
	fmt.Println(login.UserName)

	logginghelper.LogDebug("OUT : sGetDecodedLoginFromToken")
	return login.UserName, nil
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
func GetTokenRestrictTokenService(username string) (string, error) {
	standardClaim := jwt.StandardClaims{}
	expAt := time.Now().Add(constants.EXPIRATION_DURATION).Unix()
	standardClaim.ExpiresAt = expAt
	// expAt := time.Now().Add(constants.JWT_TOKEN_EXPAT).Unix()
	claims := model.JwtCustomClaims{
		username,
		// username,
		standardClaim,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	restrictTokenKey := confighelper.GetConfig("restrictTokenKey")
	fmt.Println("restricttokenkey", restrictTokenKey)
	tokenString, err := token.SignedString([]byte(restrictTokenKey))
	fmt.Println("tokenstring", tokenString)
	if nil != err {
		logginghelper.LogError("GetTokenRestrictedurlService SignedString() Error: ", err)
		return "", nil
	}
	tokenString = "Bearer " + tokenString
	return tokenString, nil
}
func GetUsernameFromRequestToken(request *http.Request) (string, error) {
	fmt.Println("in request token")
	token := request.Header.Get("authorization")
	fmt.Println("token", token)
	username := ""
	if len(token) == 0 || token == "" {
		logginghelper.LogDebug("ERR_TOKEN_NOT_FOUND : ")
		return username, errors.New("ERR_TOKEN_NOT_FOUND")
	}
	fmt.Println(username)
	secretKey := confighelper.GetConfig("JWTSecretKey")
	fmt.Println(secretKey)
	tokenClaimsMap, err := authhelper.DecodeToken(token, secretKey)
	fmt.Println("tokenclaimspmap", tokenClaimsMap)
	if nil != err {
		logginghelper.LogDebug("ERR_DECODING_TOKEN : ", err)
		return username, errors.New("ERR_DECODING_TOKEN")
	}
	logginghelper.LogDebug(tokenClaimsMap)
	username = tokenClaimsMap["username"].(string)
	fmt.Println(username)
	return username, nil
}

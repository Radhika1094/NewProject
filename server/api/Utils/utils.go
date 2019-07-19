package Utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	cryptorand "crypto/rand"
	"io"
	"math/rand"

	"strings"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/authhelper"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/confighelper"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"

	"time"
)

func IsEmptyString(data string) bool {
	data = strings.TrimSpace(data)
	if len(data) == 0 {
		return true
	}
	return false
}

func IsArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetUsernameFromRequestToken(request *http.Request) (string, error) {
	token := request.Header.Get("authorization")
	username := ""
	if len(token) == 0 || token == "" {
		logginghelper.LogDebug("ERR_TOKEN_NOT_FOUND : ")
		return username, errors.New("ERR_TOKEN_NOT_FOUND")
	}
	secretKey := confighelper.GetConfig("JWTSecretKey")
	tokenClaimsMap, err := authhelper.DecodeToken(token, secretKey)
	if nil != err {
		logginghelper.LogDebug("ERR_DECODING_TOKEN : ", err)
		return username, errors.New("ERR_DECODING_TOKEN")
	}
	logginghelper.LogDebug(tokenClaimsMap)
	username = tokenClaimsMap["username"].(string)
	return username, nil
}

func GetUsernameFromRequestTokenForCitizen(request *http.Request) (string, error) {
	token := request.Header.Get("authorization")
	username := ""
	if len(token) == 0 || token == "" {
		logginghelper.LogDebug("ERR_TOKEN_NOT_FOUND : ")
		return username, errors.New("ERR_TOKEN_NOT_FOUND")
	}
	secretKey := confighelper.GetConfig("JWTSecretKey")
	tokenClaimsMap, err := authhelper.DecodeToken(token, secretKey)
	if nil != err {
		logginghelper.LogDebug("ERR_DECODING_TOKEN : ", err)
		return username, errors.New("ERR_DECODING_TOKEN")
	}
	username = tokenClaimsMap["username"].(string)
	return username, nil
}

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

//GET RANDOM STRING
func GetRandamStringOfLen(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//ENCODE STRING
func encodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(cryptorand.Reader, b, max)
	if n != max {
		logginghelper.LogError(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func SendPostRequest(jsonStr []byte, url string, headerMap map[string]string) ([]byte, int, error) {
	logginghelper.LogDebug("IN: sendPostRequest")
	logginghelper.LogDebug("sendPostRequest to URL ---> ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	if nil != err {
		logginghelper.LogError(err)
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for headerKey, headerValue := range headerMap {
		req.Header.Set(headerKey, headerValue)
	}
	timeout := time.Duration(100 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		logginghelper.LogError(err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	logginghelper.LogDebug("response Status:", resp.StatusCode)
	responseByte, berr := ioutil.ReadAll(resp.Body)
	if nil != berr {
		logginghelper.LogError(berr)
		return nil, 0, berr
	}
	// status, _ := strconv.Atoi(resp.StatusCode)
	return responseByte, resp.StatusCode, nil
}
func SendGetRequest(url string, headerMap map[string]string) ([]byte, int, error) {
	logginghelper.LogDebug("IN: sendGetRequest")
	logginghelper.LogDebug("sendGetRequest to URL ---> ", url)
	req, err := http.NewRequest("GET", url, nil)
	// req.Header.Set("X-Custom-Header", "myvalue")
	if nil != err {
		logginghelper.LogError(err)
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for headerKey, headerValue := range headerMap {
		req.Header.Set(headerKey, headerValue)
	}
	timeout := time.Duration(100 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		logginghelper.LogError(err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	logginghelper.LogDebug("response Status:", resp.StatusCode)
	responseByte, berr := ioutil.ReadAll(resp.Body)
	if nil != berr {
		logginghelper.LogError(berr)
		return nil, 0, berr
	}
	// status, _ := strconv.Atoi(resp.StatusCode)
	return responseByte, resp.StatusCode, nil
}

func IsCFRApplication(applicationId string) bool {
	if (strings.HasPrefix(applicationId, "PC")) || (strings.HasPrefix(applicationId, "NC")) {
		return true
	}
	return false
}
func IsIFRApplication(applicationId string) bool {
	if (strings.HasPrefix(applicationId, "PI")) || (strings.HasPrefix(applicationId, "NI")) {
		return true
	}
	return false
}

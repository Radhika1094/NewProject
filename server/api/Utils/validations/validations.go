package validations

import (
	"regexp"
	"strconv"
	"strings"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"
)

//CHECKS IF FIELD IS EMPTY
func IsEmpty(fields ...string) bool {
	logginghelper.LogInfo("Inside Validation util:: IsEmpty")
	for _, field := range fields {
		if len(field) == 0 {
			return true
		}
	}
	return false
}

//CHECH IF FIELD IS NUMBER
func IsNumber(fields ...string) bool {
	logginghelper.LogInfo("Inside Validation util:: IsNumber")
	for _, field := range fields {
		_, err := strconv.Atoi(field)
		if err != nil {
			return false
		}
	}
	return true
}

//CHECK IF FIELD IS MOBILE NUMBER
func IsMobileNumber(field string) bool {
	logginghelper.LogInfo("Inside Validation util:: isMobileNumber")
	matched, _ := regexp.MatchString("^[7-9][0-9]{9}$", field)
	if matched {
		return true
	} else {
		return false
	}
}

//CHECK IF FIELD IS VALID EMAIL
func IsValidEmail(email string) bool {
	logginghelper.LogInfo("Inside Validation util:: IsValidEmail")
	// matched, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email)
	matched, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, email)

	if matched {
		// valid email
		return true
	} else {
		// Invalid email
		return false
	}
}

func ReturnCount(i int) int {
	count := 0
	for i != 0 {
		i /= 10
		count = count + 1
	}
	return count
}

func LeftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}

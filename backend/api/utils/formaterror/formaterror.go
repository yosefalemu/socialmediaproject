package formaterror

import (
	"fmt"
	"strings"
)

var errorMessages = make(map[string]string)
var err error

func FormatError(errString string) map[string]string {
	fmt.Println("FORMAT STRING ERROR", errString)
	if strings.Contains(errString, "username") {
		errorMessages["Taken_username"] = "Username Already Taken"
	}
	if strings.Contains(errString, "hashedPassword") {
		errorMessages["Incorrect_password"] = "Incorrect Password"
	}
	if len(errorMessages) > 0 {
		return errorMessages
	}
	if len(errorMessages) == 0 {
		errorMessages["Incorrect_details"] = errString
		return errorMessages
	}
	return nil
}

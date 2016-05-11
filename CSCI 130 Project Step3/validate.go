package myapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)


func getCode(data string) string {
	//wow matt. our key is so secure.....
	h := hmac.New(sha256.New, []byte("key"))
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}


func tampered(s string) bool {
	xs := strings.Split(s, "|")
	//1 is our model data
	usrData := xs[1]
	//2 is our hmac
	usrCode := xs[2]

	//so we hmac our model data and it should
	//equal to our hmac code otherwise it means the user
	//messed up the cookie.
	return usrCode != getCode(usrData)
}

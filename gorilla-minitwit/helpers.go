package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func format_datetime(timestamp string) string {
	// Convert the string to an int64 (timestamp in seconds)
	unixTime, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return "format_datetime error"
	}

	// Create a time.Time object from the Unix timestamp
	t := time.Unix(unixTime, 0)

	// Format the time into the desired display format
	return t.Format("2006-01-02 @ 15:04")
}

func isNil(i interface{}) bool {
	if i == nil || i == interface{}(nil) {
		return true
	} else {
		return false
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(userEnteredPwd string, dbpwd string) bool {
	hash := md5.Sum([]byte(userEnteredPwd))
	str := hex.EncodeToString(hash[:])
	fmt.Println("str: ", str)
	fmt.Println("dbpwd: ", dbpwd)
	return str == dbpwd
}

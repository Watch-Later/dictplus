package util

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

// WrapErrors 把多个错误合并为一个错误.
func WrapErrors(allErrors ...error) (wrapped error) {
	for _, err := range allErrors {
		if err != nil {
			if wrapped == nil {
				wrapped = err
			} else {
				wrapped = fmt.Errorf("%v | %v", err, wrapped)
			}
		}
	}
	return
}

// ErrorContains returns NoCaseContains(err.Error(), substr)
// Returns false if err is nil.
func ErrorContains(err error, substr string) bool {
	if err == nil {
		return false
	}
	return noCaseContains(err.Error(), substr)
}

// noCaseContains reports whether substr is within s case-insensitive.
func noCaseContains(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// Panic panics if err != nil
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func PathIsNotExist(name string) (ok bool) {
	_, err := os.Lstat(name)
	if os.IsNotExist(err) {
		ok = true
		err = nil
	}
	Panic(err)
	return
}

// PathIsExist .
func PathIsExist(name string) bool {
	return !PathIsNotExist(name)
}

func TimeNow() int64 {
	return time.Now().Unix()
}

// Base64Encode .
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode .
func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// StringIndex returns the index of a string in the slice, case-insensitivly.
// returns -1 if not found.
func StringIndexNoCase(slice []string, item string) int {
	item = strings.ToLower(item)
	for i, v := range slice {
		if strings.ToLower(v) == item {
			return i
		}
	}
	return -1
}

// DeleteFromSlice .
func DeleteFromSlice(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}

func BoolToInt(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func IntToBool(i int64) bool {
	return i >= 1
}

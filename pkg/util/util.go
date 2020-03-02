package util

import (
	"strings"
	"time"
)

func ParseName(index string) string {
	array := strings.Split(index, "-")
	if len(array) > 2 {
		return strings.Join(array[:len(array)-1], "-")
	}
	return array[0]
}

func ParseDate(index string) string {
	array := strings.Split(index, "-")
	if len(array) >= 2 {
		return array[len(array)-1]
	}
	return ""
}

func IsExpired(expiredDay int, date, format string) (bool, error) {
	indexTime, err := time.Parse(format, date)
	if err != nil {
		return false, err
	}

	if int(time.Now().Sub(indexTime).Hours()/24) >= expiredDay {
		return true, nil
	}
	return false, nil
}

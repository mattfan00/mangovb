package hash

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

func Hash(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func HashInt(i int) string {
	return Hash(strconv.Itoa(i))
}

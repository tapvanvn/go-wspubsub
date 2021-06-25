package utility

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"strings"
)

var randArray string = "abcdefghijklmnopqrstuvwxyz0123456789"
var randArrayCase string = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

//GenVerifyCode generate a verify code
func GenVerifyCode(length int) string {

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArray)
	for i := 0; i < length; i++ {
		code += string(randArray[rand.Intn(arrayLen)])
	}

	return code
}

//GenCode generate otp code
func GenCode(length int) string {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArray)
	for i := 0; i < length; i++ {
		code += string(randArray[rand.Intn(arrayLen)])
	}

	return strings.ToUpper(code)
}

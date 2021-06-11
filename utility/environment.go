package utility

import (
	"log"
	"os"
	"strings"
)

//GetEnv return environment or panic if fail
func GetEnv(k string) string {
	v := os.Getenv(k)
	return strings.TrimSpace(v)
}

//MustGetEnv return environment or panic if fail
func MustGetEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return strings.TrimSpace(v)
}

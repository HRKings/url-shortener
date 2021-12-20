package shortener

import (
	"os"
	"strconv"

	"github.com/speps/go-hashids/v2"
)

const DEFAULT_MIN_LENGTH = 5

func GetMinLength() int {
	minLength := os.Getenv("SHORTENER_MINIMUM_LENGTH")

	if minLength == "" {
		return DEFAULT_MIN_LENGTH
	}

	minLengthInt, err := strconv.Atoi(minLength)
	if err == nil {
		return DEFAULT_MIN_LENGTH
	}

	return minLengthInt
}

func GenerateShortLink(initialLink string, sequentialId int64) string {
	hashIdsData := hashids.NewData()

	hashMinLength := GetMinLength()
	hashIdsData.Salt = os.Getenv("HASHID_SALT")
	hashIdsData.MinLength = hashMinLength

	hash, _ := hashids.NewWithData(hashIdsData)
	encodedHash, _ := hash.EncodeInt64([]int64{sequentialId})

	return encodedHash
}

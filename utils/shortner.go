package shortener

import (
	"github.com/speps/go-hashids/v2"
)

func GenerateShortLink(initialLink string, sequentialId int) string {
	hashIdsData := hashids.NewData()
	hashIdsData.MinLength = 5
	hash, _ := hashids.NewWithData(hashIdsData)
	encodedHash, _ := hash.Encode([]int{sequentialId})
	return encodedHash[:5]
}

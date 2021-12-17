package shortener

import (
	"github.com/speps/go-hashids/v2"
)

func GenerateShortLink(initialLink string, sequentialId int) string {
	hd := hashids.NewData()
	hd.MinLength = 5
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{sequentialId})
	return e[:5]
}

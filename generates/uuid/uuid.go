package uuid

import (
	"strings"

	"github.com/golibs/uuid"
)

// New : 生成一个UUID
func New() string {
	// var id UUID = uuid.Rand()
	// fmt.Println(id.Hex())
	// fmt.Println(id.Raw())

	// id1, err := uuid.FromStr("1870747d-b26c-4507-9518-1ca62bc66e5d")
	// id2 := uuid.MustFromStr("1870747db26c450795181ca62bc66e5d")
	// fmt.Println(id1 == id2) // true

	return uuid.Rand().Hex()
}

// NewShow : 生成不带—的UUID
func NewShow() string {
	return strings.Replace(uuid.Rand().Hex(), "-", "", -1)
}

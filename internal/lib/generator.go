package lib

import "github.com/nrednav/cuid2"

func GetNewID() string {
	return cuid2.Generate()
}

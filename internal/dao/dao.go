package dao

import "github.com/normegil/evevulcan/internal/eveapi"

type DAOs struct {
	API eveapi.API
}

func (d DAOs) Character() Character {
	return Character{
		API: d.API,
	}
}

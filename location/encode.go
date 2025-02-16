package location

import (
	"encoding/json"
)

func Encode(self Location) ([]byte, error) {
	val, err := json.Marshal(&self)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func Decode(data []byte) (Location, error) {
	loc := Location{}
	err := json.Unmarshal(data, &loc)
	if err != nil {
		return loc, err
	}
	return loc, nil
}

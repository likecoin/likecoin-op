package model

import (
	"strconv"
	"time"

	"likecollective-indexer/openapi/api"
)

func TimeFromOptString(optStr api.OptString) (*time.Time, error) {
	v, ok := optStr.Get()
	if !ok {
		return nil, nil
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(i, 0)
	return &t, nil
}

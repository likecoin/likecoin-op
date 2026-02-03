package model

import (
	"encoding/json"

	"github.com/go-faster/jx"
)

type APIAdditionalProps map[string]jx.Raw

func MakeAPIAdditionalProps(additionalProps map[string]any) (APIAdditionalProps, error) {
	if additionalProps == nil {
		return nil, nil
	}

	additionalBytes, err := json.Marshal(additionalProps)
	if err != nil {
		return nil, err
	}

	a := make(APIAdditionalProps)
	jxDecoder := jx.DecodeBytes(additionalBytes)
	err = jxDecoder.Obj(func(d *jx.Decoder, key string) error {
		value, err := d.Raw()
		if err != nil {
			return err
		}
		a[key] = value
		return nil
	})

	if err != nil {
		return nil, err
	}

	return a, nil
}

package utils

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/asaskevich/govalidator"
)

type QueryUrl map[string]interface{}

// ParseBodyData parse json-formatted request body into given struct.
func ParseBodyArrayData(ctx context.Context, r *http.Request, data interface{}) error {
	if r.Body == nil {
		return errors.New("invalid data")
	}

	bBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "read")
	}

	err = json.Unmarshal([]byte(bBody), data)
	if err != nil {
		return errors.Wrap(err, "json")
	}

	return nil
}

// ParseBodyData parse json-formatted request body into given struct.
func ParseBodyData(ctx context.Context, r *http.Request, data interface{}) error {
	if r.Body == nil {
		return errors.New("invalid data")
	}

	bBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "read")
	}

	err = json.Unmarshal(bBody, data)
	if err != nil {
		return errors.Wrap(err, "json")
	}

	valid, err := govalidator.ValidateStruct(data)
	if err != nil {
		return errors.Wrap(err, "validate")
	}

	if !valid {
		return errors.New("invalid data")
	}

	return nil
}

func IsEmpty(data string) bool {
	if len(data) < 1 {
		return true
	}
	return false
}

func UrlQuery(q QueryUrl) string {
	tempQuery := url.Values{}
	for i, j := range q {
		tempQuery.Add(i, j.(string))
	}

	return tempQuery.Encode()
}

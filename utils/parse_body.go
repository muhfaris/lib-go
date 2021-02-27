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

// ParseBodyparse json-formatted request body into given struct.
func ParseBodyArray(ctx context.Context, r *http.Request, data interface{}) error {
	if r.Body == nil {
		return errors.New("lib-go: error parsebody, body is empty!")
	}

	bBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrapf(err, "lib-go: error read body data, %v", err)
	}

	err = json.Unmarshal([]byte(bBody), data)
	if err != nil {
		return errors.Wrapf(err, "lib-go: error unmarshal data, %v", err)
	}

	return nil
}

// ParseBody parse json-formatted request body into given struct.
func ParseBody(ctx context.Context, r *http.Request, data interface{}) error {
	if r.Body == nil {
		return errors.New("lib-go: error parsebody, body is empty!")
	}

	bBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrapf(err, "lib-go: error read body data, %v", err)
	}

	err = json.Unmarshal(bBody, data)
	if err != nil {
		return errors.Wrapf(err, "lib-go: error unmarshal data, %v", err)
	}

	valid, err := govalidator.ValidateStruct(data)
	if err != nil {
		return errors.Wrapf(err, "lib-go: error govalidator struct, %v", err)
	}

	if !valid {
		return errors.Errorf("lib-go: error validate struct, %t", valid)
	}

	return nil
}

// IsEmptyString is check parameter is empty
func IsEmptyString(data string) bool {
	if data == "" {
		return true
	}
	return false
}

// QueryURL is wrap data for query string
type QueryURL map[string]interface{}

// URLQuery is generate data to URL query string
func URLQuery(q QueryURL) string {
	tempQuery := url.Values{}
	for i, j := range q {
		tempQuery.Add(i, j.(string))
	}

	return tempQuery.Encode()
}

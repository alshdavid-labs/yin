package yin

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

func MockHTTPBody(body interface{}) io.ReadCloser {
	b, _ := json.Marshal(body)
	return ioutil.NopCloser(bytes.NewReader(b))
}

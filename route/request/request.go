package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//FromRequest get entity from request
func FromRequest(entity interface{}, r *http.Request) error {

	body, err := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &entity)

	if err != nil {
		return err
	}

	return nil
}

//ReceiveFile receive the uploaded file
func ReceiveFile(fieldName string, r *http.Request) (string, []byte, error) {

	r.ParseMultipartForm(32 << 20) // limit your max input length!

	file, header, err := r.FormFile(fieldName)

	if err != nil {

		return "", nil, err
	}

	defer file.Close()

	var buf bytes.Buffer

	_, err = io.Copy(&buf, file)

	if err != nil {
		return "", nil, err
	}
	var res = buf.Bytes()

	defer buf.Reset()

	return header.Filename, res, nil
}

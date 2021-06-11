package request

import (
	"errors"
	"strings"
)

//FormCreateTopic infomation to crete a topic
type FormCreateTopic struct {
	Path string `json:"path"`
}

//CheckValid is form data valid
func (form *FormCreateTopic) CheckValid() error {

	form.Path = strings.TrimSpace(form.Path)
	if len(form.Path) == 0 {
		return errors.New("path cannot be empty")
	}
	return nil
}

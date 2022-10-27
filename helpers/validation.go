package helpers

import (
	"fmt"

	"github.com/astaxie/beego/validation"
	"github.com/sirupsen/logrus"
)

func IsValid(obj interface{}) (bool, error) {
	valid := validation.Validation{}
	if v, err := valid.Valid(obj); !v {
		logrus.Errorf("validation system error: %v", err)
		return v, mergeValidationErrors(valid.ErrorsMap, err)
	}

	return true, nil
}

func mergeValidationErrors(m map[string][]*validation.Error, err error) error {
	finalError := err
	for _, v := range m {
		for _, e := range v {
			if finalError != nil {
				finalError = fmt.Errorf("%s:%s", finalError.Error(), e.Error())
			} else {
				finalError = e
			}
		}
	}
	return finalError
}

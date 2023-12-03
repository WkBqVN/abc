package utils

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
)

func ConvertJsonToStruct(filePath string, outputStruct interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	outputValue := reflect.ValueOf(outputStruct)
	if outputValue.Kind() != reflect.Ptr || outputValue.IsNil() {
		return errors.New("outputStruct must be a non-nil pointer to a struct")
	}
	if err := decoder.Decode(&outputStruct); err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return ErrCannotCloseJsonfile
	}
	return nil
}

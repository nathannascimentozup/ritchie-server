package slicer

import (
	"errors"
	"reflect"
	"ritchie-server/server"
)

type Slicer struct {
	slice interface{}
}

func NewSlicer(s interface{}) server.Slicer {
	return Slicer{
		slice: s,
	}
}

func (s Slicer) Interface() ([]interface{}, error) {
	sliceInterface := reflect.ValueOf(s.slice)
	if sliceInterface.Kind() != reflect.Slice {
		return nil, errors.New("slicer.Interface() given a non-slicer type")
	}

	ret := make([]interface{}, sliceInterface.Len())

	for i := 0; i < sliceInterface.Len(); i++ {
		ret[i] = sliceInterface.Index(i).Interface()
	}

	return ret, nil
}

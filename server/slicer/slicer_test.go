package slicer

import (
	"reflect"
	"testing"
)

func TestInterface(t *testing.T) {
	type args struct {
		slice interface{}
	}
	tests := []struct {
		name    string
		in      args
		out     []interface{}
		wantErr bool
	}{
		{
			name:    "empty interface",
			in:      args{slice: []interface{}{}},
			out:     []interface{}{},
			wantErr: false,
		},
		{
			name:    "interface with data",
			in:      args{slice: []string{"1", "2", "3"}},
			out:     []interface{}{"1", "2", "3"},
			wantErr: false,
		},
		{
			name:    "wrong type",
			in:      args{slice: nil},
			out:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSlicer(tt.in.slice).Interface()
			if (err != nil) != tt.wantErr {
				t.Errorf("Interface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("Interface() got = %v, want %v", got, tt.out)
			}
		})
	}
}

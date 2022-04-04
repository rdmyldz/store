package store

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestPut(t *testing.T) {
	type args struct {
		key []byte
		r   io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "1st case",
			args:    args{key: []byte("key 1"), r: strings.NewReader("data for key 1")},
			wantErr: false,
		},
		{
			name:    "2nd case",
			args:    args{key: []byte("key 2"), r: bytes.NewReader([]byte("data for key 2"))},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Put(tt.args.key, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPutBytes(t *testing.T) {
	type args struct {
		key  []byte
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "1st case",
			args:    args{key: []byte("key 3"), data: []byte("data for key 3")},
			wantErr: false,
		},
		{
			name:    "2nd case",
			args:    args{key: []byte("key 4"), data: []byte([]byte("data for key 4"))},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PutBytes(tt.args.key, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PutBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		args    args
		want    io.Reader
		wantErr bool
	}{
		{
			name:    "1st case",
			args:    args{key: []byte("key 1")},
			want:    strings.NewReader("data for key 1"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := io.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			want, err := io.ReadAll(tt.want)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("got: %q - want: %q\n", got, want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBytes(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "1st case",
			args:    args{[]byte("key 1")},
			want:    []byte("data for key 1"),
			wantErr: false,
		},
		{
			name:    "2st case - none exist file",
			args:    args{[]byte("non exist file")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBytes(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

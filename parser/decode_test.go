package parser

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_parseInt(t *testing.T) {
	type args struct {
		reader    io.ByteReader
		firstChar byte
	}
	tests := []struct {
		name    string
		args    args
		want    BnCode
		wantErr bool
	}{
		{
			name:    "Positive test case",
			args:    args{reader: bytes.NewReader([]byte("-123456e")), firstChar: 'i'},
			want:    BnCode{State: BnInt, Value: int(-123456)},
			wantErr: false,
		},
		{
			name:    "Empty int",
			args:    args{reader: bytes.NewReader([]byte("e")), firstChar: 'i'},
			want:    BnCode{State: BnInt},
			wantErr: true,
		},
		{
			name:    "Invalid first charater",
			args:    args{reader: bytes.NewReader([]byte("-123456e")), firstChar: ' '},
			want:    BnCode{},
			wantErr: true,
		},
		{
			name:    "Terminated stream",
			args:    args{reader: bytes.NewReader([]byte("")), firstChar: 'i'},
			want:    BnCode{},
			wantErr: true,
		},
		{
			name:    "Invalid character inside the stream",
			args:    args{reader: bytes.NewReader([]byte("ae")), firstChar: 'i'},
			want:    BnCode{},
			wantErr: true,
		},
		{
			name:    "Multiple zeros",
			args:    args{reader: bytes.NewReader([]byte("00e")), firstChar: 'i'},
			want:    BnCode{},
			wantErr: true,
		},
		{
			name:    "Negative zero",
			args:    args{reader: bytes.NewReader([]byte("-0e")), firstChar: 'i'},
			want:    BnCode{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInt(tt.args.reader, tt.args.firstChar)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseString(t *testing.T) {
	type args struct {
		reader    io.ByteReader
		firstChar byte
	}
	tests := []struct {
		name    string
		args    args
		want    BnCode
		wantErr bool
	}{
		{
			name:    "Positive test case",
			args:    args{reader: bytes.NewReader([]byte("1:foooooooooo")), firstChar: '1'},
			want:    BnCode{State: BnString, Value: "foooooooooo"},
			wantErr: false,
		},
		{
			name:    "Empty string",
			args:    args{reader: bytes.NewReader([]byte(":")), firstChar: '0'},
			want:    BnCode{State: BnString, Value: ""},
			wantErr: false,
		},
		{
			name:    "Invalid chars in length",
			args:    args{reader: bytes.NewReader([]byte(":foo")), firstChar: 'd'},
			want:    BnCode{State: BnString},
			wantErr: true,
		},
		{
			name:    "No length, terminator right away",
			args:    args{reader: bytes.NewReader([]byte("foo")), firstChar: ':'},
			want:    BnCode{State: BnString},
			wantErr: true,
		},
		{
			name:    "Empty stream",
			args:    args{reader: bytes.NewReader([]byte("")), firstChar: '3'},
			want:    BnCode{State: BnString},
			wantErr: true,
		},
		{
			name:    "No separator",
			args:    args{reader: bytes.NewReader([]byte("foo")), firstChar: '3'},
			want:    BnCode{State: BnString},
			wantErr: true,
		},
		{
			name:    "Empty string with non-zero length",
			args:    args{reader: bytes.NewReader([]byte(":")), firstChar: '3'},
			want:    BnCode{State: BnString},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseString(tt.args.reader, tt.args.firstChar)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDict(t *testing.T) {
	type args struct {
		reader    io.ByteReader
		firstChar byte
	}
	tests := []struct {
		name    string
		args    args
		want    BnCode
		wantErr bool
	}{
		{
			name: "Positive test case",
			args: args{reader: bytes.NewReader([]byte("5:helloi-3e4:spam3:foo3:zooli42e3:fooee")), firstChar: 'd'},
			want: BnCode{State: BnDict, Value: map[string]BnCode{
				"hello": {State: BnInt, Value: -3},
				"spam":  {State: BnString, Value: "foo"},
				"zoo":   {State: BnList, Value: []BnCode{{State: BnInt, Value: 42}, {State: BnString, Value: "foo"}}},
			}},
			wantErr: false,
		},
		{
			name:    "Empty dictionary",
			args:    args{reader: bytes.NewReader([]byte("e")), firstChar: 'd'},
			want:    BnCode{State: BnDict, Value: map[string]BnCode{}},
			wantErr: false,
		},
		{
			name:    "Positive test case",
			args:    args{reader: bytes.NewReader([]byte("")), firstChar: ' '},
			want:    BnCode{State: BnDict},
			wantErr: true,
		},
		{
			name:    "Empty buffer",
			args:    args{reader: bytes.NewReader([]byte("")), firstChar: 'd'},
			want:    BnCode{State: BnDict},
			wantErr: true,
		},
		{
			name:    "Key without a value",
			args:    args{reader: bytes.NewReader([]byte("5:hello")), firstChar: 'd'},
			want:    BnCode{State: BnDict},
			wantErr: true,
		},
		{
			name:    "Wrong sequence of keys",
			args:    args{reader: bytes.NewReader([]byte("4:spam3:foo5:helloi-3ee")), firstChar: 'd'},
			want:    BnCode{State: BnDict},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDict(tt.args.reader, tt.args.firstChar)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseList(t *testing.T) {
	type args struct {
		reader    io.ByteReader
		firstChar byte
	}
	tests := []struct {
		name    string
		args    args
		want    BnCode
		wantErr bool
	}{
		{
			name:    "Positive test case",
			args:    args{reader: bytes.NewReader([]byte("i42e3:fooe")), firstChar: 'l'},
			want:    BnCode{State: BnList, Value: []BnCode{{State: BnInt, Value: 42}, {State: BnString, Value: "foo"}}},
			wantErr: false,
		},
		{
			name:    "Empty List",
			args:    args{reader: bytes.NewReader([]byte("e")), firstChar: 'l'},
			want:    BnCode{State: BnList, Value: []BnCode{}},
			wantErr: false,
		},
		{
			name:    "Invalid first character",
			args:    args{reader: bytes.NewReader([]byte("i42e3:fooe")), firstChar: ' '},
			want:    BnCode{State: BnList},
			wantErr: true,
		},
		{
			name:    "Empty stream",
			args:    args{reader: bytes.NewReader([]byte("")), firstChar: 'l'},
			want:    BnCode{State: BnList},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseList(tt.args.reader, tt.args.firstChar)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decode(t *testing.T) {
	type args struct {
		reader    io.ByteReader
		firstChar byte
	}
	tests := []struct {
		name    string
		args    args
		want    BnCode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decode(tt.args.reader, tt.args.firstChar)
			if (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		reader io.ByteReader
	}
	tests := []struct {
		name    string
		args    args
		want    BnCode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

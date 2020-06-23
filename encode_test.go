package bencode

import (
	"reflect"
	"testing"
)

func Test_flattenInt(t *testing.T) {
	type args struct {
		src BnCode
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Positive test case",
			args: args{
				src: BnCode{State: BnInt, Value: -42},
			},
			want:    []byte("i-42e"),
			wantErr: false,
		},
		{
			name: "Invalid state",
			args: args{
				src: BnCode{State: BnString, Value: -42},
			},
			want:    []byte(""),
			wantErr: true,
		},
		{
			name: "Not an int value",
			args: args{
				src: BnCode{State: BnInt, Value: nil},
			},
			want:    []byte(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flattenInt(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("flattenInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flattenInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flattenString(t *testing.T) {
	type args struct {
		src BnCode
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Positive test case",
			args: args{
				src: BnCode{State: BnString, Value: "foobar"},
			},
			want:    []byte("6:foobar"),
			wantErr: false,
		},
		{
			name: "Empty string",
			args: args{
				src: BnCode{State: BnString, Value: ""},
			},
			want:    []byte("0:"),
			wantErr: false,
		},
		{
			name: "Invalid state",
			args: args{
				src: BnCode{State: BnInt, Value: "foobar"},
			},
			want:    []byte(""),
			wantErr: true,
		},
		{
			name: "Not an string value",
			args: args{
				src: BnCode{State: BnString, Value: nil},
			},
			want:    []byte(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flattenString(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("flattenString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flattenString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flattenList(t *testing.T) {
	type args struct {
		src BnCode
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Positive test case",
			args: args{
				src: BnCode{State: BnList, Value: []BnCode{{State: BnInt, Value: 42}, {State: BnString, Value: "foobar"}}},
			},
			want:    []byte("li42e6:foobare"),
			wantErr: false,
		},
		{
			name: "Invalid state",
			args: args{
				src: BnCode{State: BnDict, Value: []BnCode{{State: BnInt, Value: 42}, {State: BnString, Value: "foobar"}}},
			},
			want:    []byte(""),
			wantErr: true,
		},
		{
			name: "Not an list value",
			args: args{
				src: BnCode{State: BnList, Value: ""},
			},
			want:    []byte(""),
			wantErr: true,
		},
		{
			name: "Empty list",
			args: args{
				src: BnCode{State: BnList, Value: []BnCode{}},
			},
			want:    []byte("le"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flattenList(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("flattenList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flattenList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flattenDict(t *testing.T) {
	type args struct {
		src BnCode
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Positive test case",
			args: args{
				src: BnCode{State: BnDict, Value: map[string]BnCode{"z": {State: BnInt, Value: 42}, "a": {State: BnString, Value: "foobar"}}},
			},
			want:    []byte("d1:a6:foobar1:zi42ee"),
			wantErr: false,
		},
		{
			name: "Invalid state",
			args: args{
				src: BnCode{State: BnList, Value: map[string]BnCode{"a": {State: BnInt, Value: 42}, "b": {State: BnString, Value: "foobar"}}},
			},
			want:    []byte(""),
			wantErr: true,
		},
		{
			name: "Not an dictionary value",
			args: args{
				src: BnCode{State: BnDict, Value: ""},
			},
			want:    []byte(""),
			wantErr: true,
		},
		{
			name: "Empty dictionary",
			args: args{
				src: BnCode{State: BnDict, Value: map[string]BnCode{}},
			},
			want:    []byte("de"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flattenDict(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("flattenDict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flattenDict() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		src BnCode
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Integer",
			args: args{
				src: BnCode{State: BnInt, Value: 42},
			},
			want:    []byte("i42e"),
			wantErr: false,
		},
		{
			name: "String",
			args: args{
				src: BnCode{State: BnString, Value: "foo"},
			},
			want:    []byte("3:foo"),
			wantErr: false,
		},
		{
			name: "List",
			args: args{
				src: BnCode{State: BnList, Value: []BnCode{}},
			},
			want:    []byte("le"),
			wantErr: false,
		},
		{
			name: "Dictionary",
			args: args{
				src: BnCode{State: BnDict, Value: map[string]BnCode{}},
			},
			want:    []byte("de"),
			wantErr: false,
		},
		{
			name: "Invalid type",
			args: args{
				src: BnCode{State: 7, Value: nil},
			},
			want:    []byte(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

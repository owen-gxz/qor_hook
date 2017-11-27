package utils

import (
	"github.com/qor/admin"
	"reflect"
	"github.com/qor/media/oss"
)

var (
	TString  = reflect.TypeOf("")
	TBool    = reflect.TypeOf(true)
	TInt     = reflect.TypeOf(int(0))
	TInt8    = reflect.TypeOf(int8(0))
	TInt32   = reflect.TypeOf(int32(0))
	TInt64   = reflect.TypeOf(int64(0))
	TUint    = reflect.TypeOf(uint(0))
	TUint8   = reflect.TypeOf(uint8(0))
	TUint16  = reflect.TypeOf(uint16(0))
	TUint32  = reflect.TypeOf(uint32(0))
	TUint64  = reflect.TypeOf(uint64(0))
	TFloat32 = reflect.TypeOf(float32(0.0))
	TFloat64 = reflect.TypeOf(float64(0.0))

	ModelName = "ModelName"
)

func SetNewValue(re *admin.Resource, m map[string]string) {
	rts := reflect.TypeOf(re.Value).Elem()
	nrs := make([]reflect.StructField, 0)
	for i := 0; i < rts.NumField(); i++ {
		nu := rts.Field(i)
		nrs = append(nrs, nu)
	}
	n := getReflect(m)
	for i := 0; i < len(n); i++ {
		nrs = append(nrs, n[i])
	}
	nsnrs := reflect.StructOf(nrs)
	re.Value = reflect.New(nsnrs).Interface()
}

func DelValueKey(re *admin.Resource, k string) {
	rts := reflect.TypeOf(re.Value).Elem()
	nrs := make([]reflect.StructField, 0)
	for i := 0; i < rts.NumField(); i++ {
		if rts.Field(i).Name == Upper(k) {
			continue
		}
		nu := rts.Field(i)
		nrs = append(nrs, nu)
	}
	nsnrs := reflect.StructOf(nrs)
	re.Value = reflect.New(nsnrs).Interface()
}

func getReflect(m map[string]string) []reflect.StructField {
	nrs := make([]reflect.StructField, 0)
	for k, v := range m {
		nu := reflect.StructField{
			Name: Upper(k),
			Type: getType(v),
		}
		nrs = append(nrs, nu)
	}
	return nrs
}

func getType(ts string) reflect.Type {
	var t reflect.Type
	switch ts {
	case "string":
		fallthrough
	case "password":
		fallthrough
	case "rich_editor":
		fallthrough
	case "single_edit":
		fallthrough
	case "select_one":
		fallthrough
	case "datetime":
		t = TString
	case "number":
		fallthrough
	case "int":
		t = TInt
	case "OSS":
		fallthrough
	case "file":
		t = reflect.TypeOf(oss.OSS{})
	}
	return t
}




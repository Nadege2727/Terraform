// generated by stringer -type=ValueType valuetype.go; DO NOT EDIT

package schema

import "fmt"

const _ValueType_name = "TypeInvalidTypeBoolTypeIntTypeFloatTypeStringTypeListTypeMapTypeSetTypeJSONtypeObject"

var _ValueType_index = [...]uint8{0, 11, 19, 26, 35, 45, 53, 60, 67, 75, 85}

func (i ValueType) String() string {
	if i < 0 || i >= ValueType(len(_ValueType_index)-1) {
		return fmt.Sprintf("ValueType(%d)", i)
	}
	return _ValueType_name[_ValueType_index[i]:_ValueType_index[i+1]]
}

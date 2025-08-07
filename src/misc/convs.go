package misc

import (
	"fmt"
	"reflect"
	"strconv"

)

func AnyToInt(v any) (int, error) {
    switch val := v.(type) {
    case int:
        return val, nil

    case int8, int16, int32, int64:
        return int(reflect.ValueOf(val).Int()), nil

    case uint, uint8, uint16, uint32, uint64:
        return int(reflect.ValueOf(val).Uint()), nil

    case float32, float64:
        return int(reflect.ValueOf(val).Float()), nil

    case string:
        i, err := strconv.Atoi(val)
        return i, err

    default:
        return 0, fmt.Errorf("type %T can't be converted to a integer", v)

    }
	
}
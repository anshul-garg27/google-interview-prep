package util

import (
	"strconv"
)

func InterfaceToInt(val interface{}) int {
	if val == nil {
		return 0
	} else if intVal, ok := val.(int); ok {
		return intVal
	} else if floatVal, ok := val.(float64); ok {
		return int(floatVal)
	} else if boolVal, ok := val.(bool); ok {
		if boolVal {
			return 1
		} else {
			return 0
		}
	} else if strVal, ok := val.(string); ok {
		if intVal, err := strconv.Atoi(strVal); err == nil {
			return intVal
		}
	}
	return 0
}

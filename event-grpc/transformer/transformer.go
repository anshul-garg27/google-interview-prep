package transformer

import (
	"strconv"
)

func IntSliceToIntPointerSlice(a []int) []*int {
	b := make([]*int, len(a))
	for i := 0; i < len(a); i++ {
		b[i] = &a[i]
	}
	return b
}

func InterfaceToInt(val interface{}) int {
	if intVal, ok := val.(int); ok {
		return intVal
	} else if intVal, ok := val.(int32); ok {
		return int(intVal)
	} else if floatVal, ok := val.(float64); ok {
		return int(floatVal)
	} else if strVal, ok := val.(string); ok {
		if intVal, err := strconv.Atoi(strVal); err == nil {
			return intVal
		}
	}
	return 0
}

func InterfaceToInt64(val interface{}) int64 {
	if intVal, ok := val.(int); ok {
		return int64(intVal)
	} else if floatVal, ok := val.(float64); ok {
		return int64(floatVal)
	} else if strVal, ok := val.(string); ok {
		if intVal, err := strconv.Atoi(strVal); err == nil {
			return int64(intVal)
		}
	} else if intVal, ok := val.(int64); ok {
		return intVal
	}
	return 0
}

func FormatCount(count int) *string {

	if count < 0 {
		count = 0
	}

	formattedCount := strconv.Itoa(count)

	if count >= 1000000 {
		formattedCount = strconv.Itoa(int(count/1000000)) + "M"
	} else if count >= 100000 {
		formattedCount = strconv.Itoa(int(count/1000)) + "K"
	} else if count >= 1000 {
		thousands := int(count / 1000)
		decimal := int((count - (thousands * 1000)) / 100)

		if decimal > 0 {
			formattedCount = strconv.Itoa(thousands) + "." + strconv.Itoa(decimal) + "K"
		} else {
			formattedCount = strconv.Itoa(thousands) + "K"
		}
	}

	return &formattedCount
}

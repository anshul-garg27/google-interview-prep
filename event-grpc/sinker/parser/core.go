package parser

import "gorm.io/gorm/utils"

func getMap(m interface{}) map[string]interface{} {
	if m == nil {
		return make(map[string]interface{})
	} else {
		return m.(map[string]interface{})
	}
}

func getString(i interface{}) string {
	if i == nil {
		return ""
	} else {
		return utils.ToString(i)
	}
}

func getFloat64(i interface{}) float64 {
	if i == nil {
		return 0
	} else {
		return i.(float64)
	}
}

func getInt64(i interface{}) int64 {
	if i == nil {
		return 0
	} else {
		return i.(int64)
	}
}

func getBool(i interface{}) bool {
	if i == nil {
		return false
	} else {
		return i.(bool)
	}
}

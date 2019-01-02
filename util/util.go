package util

import (
	"reflect"
	"strings"
	"github.com/gin-gonic/gin/json"
	jsonAlias "encoding/json"
	"time"
	"runtime"
	"bytes"
	"strconv"
)

/**
	结构体对象转map(key由驼峰转蛇形)
 */
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		snakeName, _ := SnakeString(obj1.Field(i).Name)
		data[snakeName] = obj2.Field(i).Interface()
	}
	return data
}

func ResolvePointValue(value interface{}) interface{} {
	valueType := reflect.TypeOf(value).Kind()

	if valueType == reflect.Ptr {
		direct := reflect.Indirect(reflect.ValueOf(value))
		if direct.CanAddr() {
			return direct.Interface()
		} else {
			return new(interface{})
		}
	} else {
		return value
	}
}

/**
	结构体对象转map(key由驼峰转蛇形)(只转换指定列)
 */
func StructToMapWithColumns(obj interface{}, columns map[string]bool) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		fieldName := obj1.Field(i).Name
		if _, exist := columns[fieldName]; exist {
			snakeName, _ := SnakeString(fieldName)
			data[snakeName] = obj2.Field(i).Interface()
		}
	}
	return data
}

func SnakeString(s string) (string, bool) {
	data := make([]byte, 0, len(s)*2)
	change := false
	j := false
	pre := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if d >= 'A' && d <= 'Z' {
			if i > 0 && j && pre {
				change = true
				data = append(data, '_')
			}
		} else {
			pre = true
		}

		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:])), change
}

func MapToString(obj map[string]interface{}) string {
	jsonObj, _ := json.Marshal(obj)
	return string(jsonObj)
}

func SliceToString(obj []interface{}) string {
	jsonObj, _ := json.Marshal(obj)
	return string(jsonObj)
}

//将数组里面的对象key从驼峰LoadAve转化为蛇形load_ave
func SliceToSnakeSlice(obj []interface{}) []interface{} {
	for _, item := range obj {
		item := item.(map[string]interface{})
		for key, value := range item {
			snakeKey, _ := SnakeString(key)
			if snakeKey != key {
				delete(item, key)
				item[snakeKey] = value
			}
		}
	}
	return obj
}

func StringToMap(obj string) interface{} {
	rawByte := []byte(obj)
	var result map[string]interface{}
	jsonAlias.Unmarshal(rawByte, &result)
	return result
}

func StringToStringArray(obj string) []string {
	var lists []string
	dec := jsonAlias.NewDecoder(strings.NewReader(obj))
	if err := dec.Decode(&lists); err != nil {
		panic(err)
	}
	return lists
}

func StringToMapArray(obj string) []interface{} {
	var lists []interface{}
	dec := jsonAlias.NewDecoder(strings.NewReader(obj))
	if err := dec.Decode(&lists); err != nil {
		panic(err)
	}
	return lists
}

func FetchMapUnixTime(obj map[string]interface{}, key string) (i int) {
	if t, exist := obj[key]; exist {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			i = int(t.(float64) / 1000)
		case reflect.Int64:
			i = int(t.(int64) / 1000)
		}

	} else {
		i = 0
	}
	return
}

func FetchMapInterface(obj map[string]interface{}, key string) (s interface{}) {
	if t, exist := obj[key]; exist {
		s = t.(interface{})
	}
	return
}

func FetchMapString(obj map[string]interface{}, key string) (s string) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = strconv.FormatFloat(t.(float64), 'f', 6, 64)
		case reflect.Int64:
			s = strconv.FormatInt(t.(int64), 10)
		case reflect.Int32:
			s = strconv.FormatInt(int64(t.(int32)), 10)
		case reflect.Int:
			s = strconv.Itoa(t.(int))
		default:
			s = t.(string)
		}
	}
	return
}

func FetchMapInt(obj map[string]interface{}, key string) (s int) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = int(t.(float64))
		case reflect.Int64:
			s = int(t.(int64))
		case reflect.Int32:
			s = int(t.(int32))
		case reflect.String:
			if temp, err := strconv.Atoi(t.(string)); err != nil {
				panic(err)
			} else {
				s = temp
			}
		default:
			s = t.(int)
		}
	}
	return
}

func FetchMapInt64(obj map[string]interface{}, key string) (s int64) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = int64(t.(float64))
		case reflect.Int32:
			s = int64(t.(int32))
		case reflect.String:
			if temp, err := strconv.Atoi(t.(string)); err != nil {
				panic(err)
			} else {
				s = int64(temp)
			}
		default:
			s = t.(int64)
		}
	}
	return
}

func FetchMapFloat64(obj map[string]interface{}, key string) (s float64) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = t.(float64)
		case reflect.Int64:
			s = float64(t.(int64))
		}
	}
	return
}

func FetchMapFloat32(obj map[string]interface{}, key string) (s float32) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Float64:
			s = float32(t.(float64))
		default:
			s = t.(float32)
		}
	}
	return
}

func FetchMapBool(obj map[string]interface{}, key string) (b bool) {
	if t, exist := obj[key]; exist && t != nil {
		b = t.(bool)
	}
	return
}

func FetchMapSlice(obj map[string]interface{}, key string) (s []interface{}) {
	if t, exist := obj[key]; exist {
		s = t.([]interface{})
	}
	return
}

func FetchMapMap(obj map[string]interface{}, key string) (map[string]interface{}) {
	if t, exist := obj[key]; exist && t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Ptr:
			temp := ResolvePointValue(t)
			if temp != nil {
				return temp.(map[string]interface{})
			}
		default:
			return t.(map[string]interface{})
		}
	}
	return nil
}

func BoolToInt(obj bool) int {
	if obj == true {
		return 1
	} else {
		return 0
	}
}

func GetCurrentMilliSecond() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func MillisecondToTime(mtimestamp int64) time.Time {
	second := int64(mtimestamp / 1000)
	msecond := int64(mtimestamp % 1000)
	return time.Unix(second, msecond*1e6)
}

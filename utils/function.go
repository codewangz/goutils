package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Ucfirst(str string) string {
	if len(str) < 1 {
		return ""
	}
	stdArr := []rune(str)
	if stdArr[0] >= 97 && stdArr[0] <= 122 {
		stdArr[0] -= 32
	}
	return string(stdArr)
}

func Lcfirst(str string) string {
	if len(str) < 1 {
		return ""
	}
	stdArr := []rune(str)
	if stdArr[0] >= 65 && stdArr[0] <= 90 {
		stdArr[0] += 32
	}
	return string(stdArr)
}

func ItoString(val interface{}) string {
	str := ""
	switch val.(type) {
	case float64:
		str = strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case string:
		str = val.(string)
	case int64:
		str = strconv.FormatInt(val.(int64), 10)
	case int:
		str = strconv.Itoa(val.(int))
	case error:
		str = val.(error).Error()
	}
	return str
}

func ToSliceString(val interface{}) []string {
	result := []string{}
	switch x := val.(type) {
	case []interface{}:
		for _, v := range x {
			result = append(result, ItoString(v))
		}
	case []string:
		for _, v := range x {
			result = append(result, v)
		}
	}
	return result
}

func ToSliceInt64(val interface{}) []int64 {
	var result []int64
	switch x := val.(type) {
	case []interface{}:
		for _, v := range x {
			result = append(result, ToInt64(v))
		}
	}
	return result
}

func ToMapInterface(val interface{}) map[string]interface{} {
	switch val.(type) {
	case map[string]interface{}:
		return val.(map[string]interface{})
	}
	return nil
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

func CheckPhone(tel string) bool {

	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`

	rgx := regexp.MustCompile(reg)

	return rgx.MatchString(tel)
}

func SortMap(data map[string]interface{}) []interface{} {
	var keys []string
	for key, _ := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var sliceData []interface{}
	for _, key := range keys {
		sliceData = append(sliceData, data[key])
	}
	return sliceData
}

func InSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func ToInt(val interface{}) int {
	num := 0
	switch val.(type) {
	case float64:
		num = int(val.(float64))
	case string:
		num, _ = strconv.Atoi(val.(string))
	case int64:
		num = int(val.(int64))
	case int:
		num = val.(int)
	}
	return num
}

func ToInt64(val interface{}) int64 {
	var num int64
	switch val.(type) {
	case float64:
		num = int64(val.(float64))
	case string:
		num, _ = strconv.ParseInt(val.(string), 10, 64)
	case int64:
		num = val.(int64)
	case int:
		num = int64(val.(int))
	}
	return num
}

func DeleteFromSlice(element interface{}, slice interface{}) {
	switch slice := (slice).(type) {
	case *[]string:
		for key, val := range *slice {
			if val == element {
				*slice = append((*slice)[:key], (*slice)[key+1:]...)
			}
		}
	}
}

func SliceInt64Sum(slice []int64) (result int64) {
	for _, val := range slice {
		result += val
	}
	return
}

func ToSliceInterface(val interface{}) (result []interface{}) {
	switch val.(type) {
	case []interface{}:
		return val.([]interface{})
	}
	return
}

func MapMerge(map1, map2 map[string]interface{}) (result map[string]interface{}) {
	result = Copy(map1).(map[string]interface{})
	for key, val := range map2 {
		result[key] = val
	}
	return
}

func ToSnakes(data []map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, val := range data {
		temp := map[string]interface{}{}
		for key, value := range val {
			temp[snakeString(key)] = value
		}
		result = append(result, temp)
	}
	return result
}

func ToSnake(data map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for key, value := range data {
		result[snakeString(key)] = value
	}
	return result
}

/**
 * 驼峰转蛇形 snake string
 * @description XxYy to xx_yy , XxYY to xx_y_y
 * @date 2020/7/30
 * @param s 需要转换的字符串
 * @return string
 **/
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

/**
 * 蛇形转驼峰
 * @description xx_yy to XxYx  xx_y_y to XxYY
 * @date 2020/7/30
 * @param s要转换的字符串
 * @return string
 **/
func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

func JsonDecode(jsonStr string) (result interface{}) {
	if jsonStr == "" {
		return
	}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		panic(err)
	}
	return
}

func JsonEncode(obj interface{}) string {
	if obj == nil {
		return ""
	}
	result, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(result)
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Copy(data interface{}) interface{} {
	var result interface{}
	jsonByte, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonByte, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func MD5(str string) string {
	md5 := md5.New()
	md5.Write([]byte(str))
	return hex.EncodeToString(md5.Sum(nil))
}

func MapJoin(m map[string]interface{}, split string) (str string) {
	for _, v := range m {
		if str != "" {
			str += split + ItoString(v)
		} else {
			str += ItoString(v)
		}
	}
	return
}

func ToSliceMap(val interface{}) (result []map[string]interface{}) {
	switch item := val.(type) {
	case []interface{}:
		for _, va := range item {
			switch v := va.(type) {
			case map[string]interface{}:
				result = append(result, v)
			}
		}
	}
	return
}

func ToFloat64(val interface{}) (result float64) {
	switch item := val.(type) {
	case string:
		result, _ = strconv.ParseFloat(item, 64)
		return
	case int64:
		return float64(item)
	case int:
		return float64(item)
	case float64:
		return item
	case float32:
		return float64(item)
	case nil:
		return 0
	}
	return 0
}

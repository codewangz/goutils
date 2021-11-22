package service

import (
	"fmt"
	"goutils/utils"
	"reflect"
	"strconv"
	"strings"
)

type base struct {
}

var serviceMap map[string]interface{}

type Result interface{}

type Params map[string]interface{}

func (p Params) GetInt(key string, defVal ...int64) int64 {
	val := p[key]
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	switch val.(type) {
	case float64:
		return int64(val.(float64))
	case int64:
		return val.(int64)
	case int:
		return int64(val.(int))
	case string:
		if v, err := strconv.ParseInt(val.(string), 10, 64); err == nil {
			return v
		}
	}

	return 0
}

func (p Params) GetFloat(key string, defVal ...float64) float64 {
	val := p[key]
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	switch val.(type) {
	case float64:
		return val.(float64)
	case string:
		if v, err := strconv.ParseFloat(val.(string), 64); err == nil {
			return v
		}
	}

	return 0
}

func (p Params) GetString(key string, defVal ...string) string {
	val := p[key]
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	switch val.(type) {
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case string:
		return val.(string)
	}

	return ""
}

func (p Params) GetMap(key string, defVal ...map[string]interface{}) map[string]interface{} {
	val := p[key]
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	switch val.(type) {
	case map[string]interface{}:
		return val.(map[string]interface{})
	}
	return nil
}

func (p Params) GetSlice(key string, defVal ...[]interface{}) []interface{} {
	val := p[key]
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	switch val.(type) {
	case []interface{}:
		return val.([]interface{})
	}
	return nil
}

func (p Params) GetIntSlice(key string, defVal ...[]int64) []int64 {
	result := []int64{}
	val := p[key]
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	switch val.(type) {
	case []interface{}:
		for _, v := range val.([]interface{}) {
			switch v.(type) {
			case float64:
				result = append(result, int64(v.(float64)))
			case string:
				if v, err := strconv.ParseInt(v.(string), 10, 64); err == nil {
					result = append(result, v)
				}
			}
		}
	}
	return result
}

func (p Params) GetStringSlice(key string, defVal ...[]string) []string {
	result := []string{}
	val := p[key]
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	switch val.(type) {
	case []interface{}:
		for _, v := range val.([]interface{}) {
			switch v.(type) {
			case float64:
				result = append(result, strconv.FormatFloat(v.(float64), 'f', -1, 64))
			case int:
				result = append(result, strconv.Itoa(v.(int)))
			case int64:
				result = append(result, strconv.FormatInt(v.(int64), 10))
			case string:
				result = append(result, v.(string))
			}
		}
	case []string:
		for _, v := range val.([]string) {
			result = append(result, v)
		}
	case []int:
		for _, v := range val.([]int) {
			result = append(result, strconv.Itoa(v))
		}
	case []int64:
		for _, v := range val.([]int64) {
			result = append(result, strconv.FormatInt(v, 10))
		}
	}
	return result
}

func (p Params) GetStringSliceSlice(key string, defVal ...[][]string) [][]string {
	result := [][]string{}
	val := p[key]
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	switch val.(type) {
	case []interface{}:
		for _, v := range val.([]interface{}) {
			tempV := []string{}
			switch v.(type) {
			case []interface{}:
				for _, va := range v.([]interface{}) {
					switch va.(type) {
					case float64:
						tempV = append(tempV, strconv.FormatFloat(va.(float64), 'f', -1, 64))
					case string:
						tempV = append(tempV, va.(string))
					}
				}
			}
			result = append(result, tempV)
		}
	}
	return result
}

func (p Params) GetVal(key string, defVal ...interface{}) interface{} {
	val := p[key]
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	return val
}

func Call(service string, params Params) Result {
	arr := strings.Split(service, ".")
	if len(arr) != 2 {
		panic(fmt.Errorf("params error"))
	}
	name := arr[0]
	if name == "" {
		panic(fmt.Errorf("service name error"))
	}
	method := utils.Ucfirst(arr[1])
	if method == "" {
		panic(fmt.Errorf("service method error"))
	}
	obj := serviceMap[name+"Service"]
	if obj == nil {
		panic(fmt.Errorf("%sService is not exist", name))
	}
	Method := reflect.ValueOf(obj).MethodByName(method)
	if Method.Kind() == reflect.Invalid {
		panic(fmt.Errorf("service method %s is not exist", method))
	}
	args := make([]reflect.Value, Method.Type().NumIn())
	if len(args) > 0 {
		args[0] = reflect.ValueOf(params)
	}
	result := Method.Call(args)
	if len(result) > 1 && result[1].Interface() != nil {
		panic(result[1].Interface().(error))
	}
	if (len)(result) > 0 {
		return result[0].Interface()
		/*		if result[0].Interface() != nil {
				return result[0].Interface().(Result)
			}*/
	}
	return nil
}

func RegisterService(services ...interface{}) {
	if serviceMap == nil {
		serviceMap = map[string]interface{}{}
	}
	for _, service := range services {
		value := reflect.ValueOf(service)
		value = reflect.Indirect(value)
		name := value.Type().Name()
		serviceMap[name] = service
	}
}

func init() {
	serviceMap = make(map[string]interface{})
}

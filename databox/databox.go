package databox

import (
	"goutils/utils"
	"reflect"
	"strconv"
	"strings"
)

type dataBox struct {
	data interface{}
}

func NewDataBox(args ...interface{}) dataBox {
	var data interface{}
	if len(args) > 0 {
		data = args[0]
	}
	//复制一份，防止原数据被修改
	return dataBox{utils.Copy(data)}
}

func (dbx *dataBox) Get(key string) interface{} {
	var data interface{}
	paths := strings.Split(key, ".")
	data = dbx.data
	for _, path := range paths {
		if reflect.ValueOf(data).Kind() == reflect.Invalid {
			return nil
		}
		dataType := reflect.TypeOf(data).String()
		i, err := strconv.Atoi(path)
		if err == nil && dataType == "[]interface {}" { //list
			tempData := data.([]interface{})
			if len(tempData) > i {
				data = tempData[i]
			} else {
				return nil
			}
		} else if dataType == "map[string]interface {}" {
			tempData := data.(map[string]interface{})
			data = tempData[path]
		} else {
			return nil
		}
	}

	return data
}

func (dbx *dataBox) Set(key string, val interface{}) {
	paths := strings.Split(key, ".")
	dbx.retSet(paths, val)
}

func (dbx *dataBox) retSet(paths []string, val interface{}) {
	data := dbx.createData(paths, val)
	dbx.data = dbx.mergeData(dbx.data, data, paths, 0)
}

func (dbx *dataBox) createData(paths []string, val interface{}) interface{} {

	var tempData interface{}
	if index, err := strconv.Atoi(paths[len(paths)-1]); err == nil {
		tempData = []interface{}{}
		for i := 0; i <= index; i++ {
			tempData = append(tempData.([]interface{}), nil)
		}
		tempData.([]interface{})[index] = val
	} else {
		tempData = map[string]interface{}{paths[len(paths)-1]: val}
	}

	if len(paths) == 1 {
		return tempData
	} else {
		return dbx.createData(paths[0:len(paths)-1], tempData)
	}
}

func (dbx *dataBox) mergeData(dst interface{}, src interface{}, paths []string, deep int) interface{} {
	if dst == nil {
		return src
	}

	if deep == len(paths) {
		return src
	}

	ok, kind := dbx.isSameKind(dst, src)
	if !ok {
		return src
		//panic("结构类型不一致，不能设置")
	}

	if kind == reflect.Slice {
		index, err := strconv.Atoi(paths[deep])
		if err != nil {
			panic(err)
		}
		srcSlice := utils.Copy(src.([]interface{})[index])
		if ok, dstSlice := dbx.isInPath(kind, dst, paths[deep:]); ok {
			srcSlice = dbx.mergeData(dstSlice, srcSlice, paths, deep+1)
		} else {
			copy(src.([]interface{}), dst.([]interface{}))
			dst = src
		}

		dst.([]interface{})[index] = srcSlice

	} else if kind == reflect.Map {
		srcval := utils.Copy(src.(map[string]interface{})[paths[deep]])
		if ok, dstMap := dbx.isInPath(kind, dst, paths[deep:]); ok {
			srcval = dbx.mergeData(dstMap, srcval, paths, deep+1)
		}
		dst.(map[string]interface{})[paths[0]] = srcval
	}

	return dst

}

func (dbx *dataBox) isSameKind(dst, src interface{}) (bool, reflect.Kind) {
	return reflect.ValueOf(dst).Kind() == reflect.ValueOf(src).Kind(), reflect.ValueOf(src).Kind()
}

func (dbx *dataBox) isInPath(kind reflect.Kind, dst interface{}, path []string) (bool, interface{}) {
	if kind == reflect.Map {
		if _, ok := dst.(map[string]interface{})[path[0]]; ok {
			return ok, utils.Copy(dst.(map[string]interface{})[path[0]])
		}
	} else if kind == reflect.Slice {
		if i, err := strconv.Atoi(path[0]); err == nil {
			if len(dst.([]interface{})) > i {
				return true, utils.Copy(dst.([]interface{})[i])
			}
			return false, nil
		} else {
			panic(err)
		}
	}
	return false, nil
}

func (dbx *dataBox) Data() interface{} {
	return dbx.data
}

func (dbx *dataBox) GetSlice(key string, defVal ...[]interface{}) []interface{} {
	val := dbx.Get(key)
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	return utils.ToSliceInterface(val)
}

func (dbx *dataBox) GetInt64(key string, defVal ...int64) int64 {
	val := dbx.Get(key)
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	return utils.ToInt64(val)
}

func (dbx *dataBox) GetString(key string, defVal ...string) string {
	val := dbx.Get(key)
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	return utils.ItoString(val)
}

func (dbx *dataBox) GetSliceString(key string, defVal ...[]string) []string {
	val := dbx.Get(key)
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	return utils.ToSliceString(val)
}

func (dbx *dataBox) GetSliceMap(key string, defVal ...[]map[string]interface{}) []map[string]interface{} {
	val := dbx.Get(key)
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	return utils.ToSliceMap(val)
}

func (dbx *dataBox) GetMapInterface(key string, defVal ...map[string]interface{}) map[string]interface{} {
	val := dbx.Get(key)
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	return utils.ToMapInterface(val)
}

func (dbx *dataBox) GetFloat64(key string, defVal ...float64) float64 {
	val := dbx.Get(key)
	if val == nil && len(defVal) > 0 {
		return defVal[0]
	}
	return utils.ToFloat64(val)
}

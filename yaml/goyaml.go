package yaml

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type yamlConf struct {
	filePath string
	data     interface{}
}

func NewYamlConf(filePath string) *yamlConf {
	conf := yamlConf{filePath: filePath}
	conf.init()
	return &conf
}

func (conf *yamlConf) init() {
	yamlFile, err := ioutil.ReadFile(conf.filePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &conf.data)
	if err != nil {
		panic(err)
	}
	conf.data = conf.convert(conf.data)
}

func (conf *yamlConf) Get(key string) interface{} {
	//fmt.Println(key)
	var data interface{}
	paths := strings.Split(key, ".")
	data = conf.data
	for _, path := range paths {
		dataType := reflect.TypeOf(data).String()
		i, err := strconv.Atoi(path)
		if err == nil { //list
			if dataType == "[]interface {}" {
				tempData := data.([]interface{})
				if len(tempData) > i {
					data = tempData[i]
				} else {
					return nil
				}
			} else {
				return nil
			}
		} else { //string
			if dataType == "map[string]interface {}" {
				tempData := data.(map[string]interface{})
				data = tempData[path]
			} else {
				return nil
			}
		}
	}

	return data
}

func (conf *yamlConf) convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = conf.convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = conf.convert(v)
		}
	}
	return i
}

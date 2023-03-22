package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type localCache struct {
	expireTime time.Time
	key        string
	data       interface{}
}

type localCaches map[string]localCache

var localcaches localCaches

func NewLocalCaches() localCaches {
	return localcaches
}

func (caches localCaches) Get(key string) interface{} {
	key = caches.md5Key(key)
	if cache, ok := caches[key]; ok {
		return cache.data
	}
	return nil
}

func (caches localCaches) Set(key string, args ...interface{}) {
	key = caches.md5Key(key)
	var seconds time.Duration
	if len(args) > 1 {
		seconds = args[1].(time.Duration)
	} else {
		seconds = time.Minute * 3
	}
	localCache := localCache{expireTime: time.Now().Add(seconds), key: key, data: args[0]}
	caches[key] = localCache
	//fmt.Println(key,"生成缓存")
}

func (caches localCaches) md5Key(uri string) string {
	md5 := md5.New()
	md5.Write([]byte(uri))
	return hex.EncodeToString(md5.Sum(nil))
}

func (caches localCaches) release() {
	d := time.Duration(time.Second * 5)
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<-t.C
		for key, value := range caches {
			if time.Now().After(value.expireTime) {
				//fmt.Println(key,"缓存失效")
				delete(caches, key)
			}
		}
	}
}

func (caches localCaches) CallCache(obj interface{}, method string, time time.Duration, params ...interface{}) interface{} {
	if method == "" {
		panic(fmt.Errorf("method is nil"))
	}
	if obj == nil {
		panic(fmt.Errorf("obj is nil"))
	}
	objName := reflect.ValueOf(obj).Type().String()
	jsonObjName, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	jsonp, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	key := objName + string(jsonObjName) + method + string(jsonp)
	ret := caches.Get(key)
	if ret != nil {
		return ret
	}
	Method := reflect.ValueOf(obj).MethodByName(method)
	if Method.Kind() == reflect.Invalid {
		panic(fmt.Errorf("the object %s method %s is not exist", reflect.ValueOf(obj).Type(), method))
	}
	args := make([]reflect.Value, len(params))
	for key, val := range params {
		args[key] = reflect.ValueOf(val)
	}
	result := Method.Call(args)
	if len(result) > 1 && result[1].Interface() != nil {
		panic(result[1].Interface().(error))
	}
	if (len)(result) > 0 {
		ret := result[0].Interface().(interface{})
		caches.Set(key, ret, time)
		return ret
	}
	return nil
}

func init() {
	localcaches = map[string]localCache{}
	go localcaches.release()
}

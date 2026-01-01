package GoMiniblink

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// fnJsData JavaScript 函数数据包装结构
type fnJsData struct {
	jsData
	fnName string
	fn     reflect.Value
	mb     Miniblink
}

// 初始化 JavaScript 函数数据对象
func (_this *fnJsData) init(name string) *fnJsData {
	_this.name = [100]byte{'f', 'u', 'n', 'c', 't', 'i', 'o', 'n'}
	_this.fnName = name
	return _this
}

// 将 Go 值转换为 JavaScript 值
func toJsValue(mb Miniblink, es jsExecState, value interface{}) jsValue {
	if value == nil {
		return mbApi.jsUndefined()
	}
	switch value.(type) {
	case int:
		return mbApi.jsInt(int32(value.(int)))
	case int8:
		return mbApi.jsInt(int32(value.(int8)))
	case int16:
		return mbApi.jsInt(int32(value.(int16)))
	case int32:
		return mbApi.jsInt(value.(int32))
	case int64:
		return mbApi.jsDouble(float64(value.(int64)))
	case uint:
		return mbApi.jsInt(int32(value.(uint)))
	case uint8:
		return mbApi.jsInt(int32(value.(uint8)))
	case uint16:
		return mbApi.jsInt(int32(value.(uint16)))
	case uint32:
		return mbApi.jsInt(int32(value.(uint32)))
	case uint64:
		return mbApi.jsDouble(float64(value.(uint64)))
	case float32:
		return mbApi.jsDouble(float64(value.(float32)))
	case float64:
		return mbApi.jsDouble(value.(float64))
	case bool:
		return mbApi.jsBoolean(value.(bool))
	case string:
		return mbApi.jsString(es, value.(string))
	case time.Time:
		return mbApi.jsDouble(float64(value.(time.Time).Unix()))
	default:
		break
	}
	rt := reflect.TypeOf(value)
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		length := rv.Len()
		arr := mbApi.jsEmptyArray(es)
		mbApi.jsSetLength(es, arr, uint32(length))
		for i := 0; i < length; i++ {
			v := toJsValue(mb, es, rv.Index(i).Interface())
			mbApi.jsSetAt(es, arr, uint32(i), v)
		}
		return arr
	case reflect.Map:
		obj := mbApi.jsEmptyObject(es)
		kv := rv.MapRange()
		for kv.Next() && kv.Key().Kind() == reflect.String {
			k := kv.Key().Interface().(string)
			v := toJsValue(mb, es, kv.Value().Interface())
			mbApi.jsSet(es, obj, k, v)
		}
		return obj
	case reflect.Struct:
		obj := mbApi.jsEmptyObject(es)
		for i := 0; i < rv.NumField(); i++ {
			f := rt.Field(i)
			if strings.ToUpper(f.Name)[0] == f.Name[0] {
				fname := rt.Field(i).Name
				fvalue := rv.Field(i).Interface()
				v := toJsValue(mb, es, fvalue)
				mbApi.jsSet(es, obj, fname, v)
			}
		}
		return obj
	case reflect.Func:
		rsName := "__tmpFnRs" + strconv.FormatUint(seq(), 10)
		jsFn := new(fnJsData).init("__tmpFn" + strconv.FormatUint(seq(), 10))
		jsFn.fn = rv
		if is64 {
			jsFn.callAsFunction = syscall.NewCallbackCDecl(execTempFunc)
		} else {
			jsFn.callAsFunction = syscall.NewCallbackCDecl(execTempFuncX86)
		}
		jsFn.finalize = syscall.NewCallbackCDecl(deleteTempFunc)
		keepRef[jsFn.fnName] = jsFn
		fv := mbApi.jsFunction(es, &jsFn.jsData)
		mbApi.jsSetGlobal(es, jsFn.fnName, fv)
		js := `return function(){
                 var rs=%q;
                 var fn=%q;
                 var arr=Array.prototype.slice.call(arguments);
                 var args=[fn,rs].concat(arr);
                 window[fn].apply(null,args);
                 var fnrs=window.top[rs];
                 window.top[rs]=undefined;
                 window[fn]=undefined;
                 return fnrs;
               }`
		js = fmt.Sprintf(js, rsName, jsFn.fnName)
		return mbApi.jsEval(es, js)
	}
	panic("不支持的go类型：" + rv.Kind().String() + "(" + rv.Type().String() + ")")
}

// 删除临时 JavaScript 函数的回调
func deleteTempFunc(ptr uintptr) uintptr {
	data := (*fnJsData)(unsafe.Pointer(ptr))
	delete(keepRef, data.fnName)
	return 0
}

// 在 x86 架构下执行临时函数的包装函数
func execTempFuncX86(es jsExecState, _, _, _ uintptr, count uint32) uintptr {
	return execTempFunc(es, 0, 0, count)
}

// 执行临时 JavaScript 函数
func execTempFunc(es jsExecState, _, _ jsValue, count uint32) uintptr {
	wke := mbApi.jsGetWebView(es)
	mb := views[wke]
	arr := make([]reflect.Value, count)
	for i := uint32(0); i < count; i++ {
		jv := mbApi.jsArg(es, i)
		arr[i] = reflect.ValueOf(toGoValue(mb, es, jv))
	}
	dataName := arr[0].String()
	if v, ok := keepRef[dataName]; ok {
		rsName := arr[1].String()
		rs := v.(*fnJsData).fn.Call(arr[2:])
		if len(rs) > 0 {
			jv := toJsValue(mb, es, rs[0].Interface())
			mbApi.jsSetGlobal(es, rsName, jv)
		}
	}
	delete(keepRef, dataName)
	return 0
}

// 将 JavaScript 值转换为 Go 值
func toGoValue(mb Miniblink, es jsExecState, value jsValue) interface{} {
	switch mbApi.jsTypeOf(value) {
	case jsType_NULL, jsType_UNDEFINED:
		return nil
	case jsType_NUMBER:
		return mbApi.jsToDouble(es, value)
	case jsType_BOOLEAN:
		return mbApi.jsToBoolean(es, value)
	case jsType_STRING:
		return mbApi.jsToTempString(es, value)
	case jsType_ARRAY:
		length := mbApi.jsGetLength(es, value)
		ps := make([]interface{}, length)
		for i := 0; i < length; i++ {
			v := mbApi.jsGetAt(es, value, uint32(i))
			ps[i] = toGoValue(mb, es, v)
		}
		return ps
	case jsType_OBJECT:
		ps := make(map[string]interface{})
		keys := mbApi.jsGetKeys(es, value)
		for _, k := range keys {
			v := mbApi.jsGet(es, value, k)
			ps[k] = toGoValue(mb, es, v)
		}
		return ps
	case jsType_FUNCTION:
		name := "__pofn" + strconv.FormatUint(seq(), 10)
		return JsFunc(func(param ...interface{}) interface{} {
			jses := mbApi.wkeGlobalExec(mb.GetHandle())
			ps := make([]jsValue, len(param))
			for i, v := range param {
				ps[i] = toJsValue(mb, jses, v)
			}
			rs := mbApi.jsCall(jses, value, mbApi.jsUndefined(), ps)
			mbApi.jsSetGlobal(jses, name, mbApi.jsUndefined())
			return toGoValue(mb, jses, rs)
		})
	default:
		panic("不支持的js类型：" + strconv.Itoa(int(value)))
	}
}

var seed uint64 = 0

// 生成一个递增的唯一序列号
func seq() uint64 {
	seed++
	return seed
}

// 将 Go 布尔值转换为 C 布尔值（1或0）
func toBool(b bool) byte {
	if b {
		return 1
	} else {
		return 0
	}
}

// 将 Go 字符串转换为 C 风格的 null 终止字符串
func toCallStr(str string) []byte {
	buf := []byte(str)
	rs := make([]byte, len(str)+1)
	for i, v := range buf {
		rs[i] = v
	}
	return rs
}

// 将 C 字符串指针转换为 Go UTF-8 字符串
func ptrToUtf8(ptr uintptr) string {
	var seq []byte
	for {
		b := *((*byte)(unsafe.Pointer(ptr)))
		if b != 0 {
			seq = append(seq, b)
			ptr++
		} else {
			break
		}
	}
	return string(seq)
}

// 检查指定路径是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

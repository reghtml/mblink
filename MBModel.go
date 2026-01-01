package GoMiniblink

type JsFunc func(param ...interface{}) interface{}

// GoFnContext Go 函数调用上下文
type GoFnContext struct {
	Miniblink Miniblink
	Name      string
	State     interface{}
	Param     []interface{}
}

type GoFn func(context GoFnContext) interface{}

// JsFnBinding JavaScript 函数绑定信息
type JsFnBinding struct {
	Name  string
	Fn    GoFn
	State interface{}
	core  wkeJsNativeFunction
}

// 调用绑定的 JavaScript 函数
func (_this *JsFnBinding) Call(mb Miniblink, param []interface{}) interface{} {
	ctx := GoFnContext{
		Miniblink: mb,
		Name:      _this.Name,
		State:     _this.State,
		Param:     param,
	}
	return _this.Fn(ctx)
}

package controls

import f "github.com/reghtml/mblink/forms"

// 默认加载事件处理函数
func (_this *BaseUI) defOnLoad() {
	for _, v := range _this.EvLoad {
		v(_this.instance)
	}
}

// 默认销毁事件处理函数
func (_this *BaseUI) defOnDestroy() {
	for _, v := range _this.EvDestroy {
		v(_this.instance)
	}
}

// 默认失去焦点事件处理函数
func (_this *BaseUI) defOnLostFocus() {
	for _, v := range _this.EvLostFocus {
		v(_this.instance)
	}
}

// 默认获得焦点事件处理函数
func (_this *BaseUI) defOnFocus() {
	for _, v := range _this.EvFocus {
		v(_this.instance)
	}
}

// 默认按键按下事件处理函数
func (_this *BaseUI) defOnKeyPress(e *f.KeyPressEvArgs) {
	for _, v := range _this.EvKeyPress {
		v(_this.instance, e)
	}
}

// 默认按键释放事件处理函数
func (_this *BaseUI) defOnKeyUp(e *f.KeyEvArgs) {
	for _, v := range _this.EvKeyUp {
		v(_this.instance, e)
	}
}

// 默认按键按下事件处理函数
func (_this *BaseUI) defOnKeyDown(e *f.KeyEvArgs) {
	for _, v := range _this.EvKeyDown {
		v(_this.instance, e)
	}
}

// 默认绘制事件处理函数
func (_this *BaseUI) defOnPaint(e f.PaintEvArgs) {
	for _, v := range _this.EvPaint {
		v(_this.instance, e)
	}
}

// 默认鼠标点击事件处理函数
func (_this *BaseUI) defOnMouseClick(e *f.MouseEvArgs) {
	for _, v := range _this.EvMouseClick {
		v(_this.instance, e)
	}
}

// 默认鼠标滚轮事件处理函数
func (_this *BaseUI) defOnMouseWheel(e *f.MouseEvArgs) {
	for _, v := range _this.EvMouseWheel {
		v(_this.instance, e)
	}
}

// 默认鼠标释放事件处理函数
func (_this *BaseUI) defOnMouseUp(e *f.MouseEvArgs) {
	for _, v := range _this.EvMouseUp {
		v(_this.instance, e)
	}
}

// 默认鼠标按下事件处理函数
func (_this *BaseUI) defOnMouseDown(e *f.MouseEvArgs) {
	for _, v := range _this.EvMouseDown {
		v(_this.instance, e)
	}
}

// 默认鼠标移动事件处理函数
func (_this *BaseUI) defOnMouseMove(e *f.MouseEvArgs) {
	for _, v := range _this.EvMouseMove {
		v(_this.instance, e)
	}
}

// 默认显示事件处理函数
func (_this *BaseUI) defOnShow() {
	for _, v := range _this.EvShow {
		v(_this.instance)
	}
}

// 默认大小改变事件处理函数
func (_this *BaseUI) defOnResize(e f.Rect) {
	for _, v := range _this.EvResize {
		v(_this.instance, e)
	}
}

// 默认移动事件处理函数
func (_this *BaseUI) defOnMove(e f.Point) {
	for _, v := range _this.EvMove {
		v(_this.instance, e)
	}
}

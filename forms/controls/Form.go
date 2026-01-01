package controls

import (
	fm "github.com/reghtml/mblink/forms"
	br "github.com/reghtml/mblink/forms/bridge"
)

// Form 窗体类
type Form struct {
	BaseUI
	BaseContainer

	EvState map[string]func(s GUI, state fm.FormState)
	OnState func(state fm.FormState)

	impl          br.Form
	title         string
	showInTaskbar bool
	border        fm.FormBorder
	state         fm.FormState
	startPos      fm.FormStart
}

// 获取窗体实现对象
func (_this *Form) getFormImpl() br.Form {
	return _this.impl
}

// 使用参数初始化窗体
func (_this *Form) InitEx(param br.FormParam) *Form {
	_this.impl = App.NewForm(param)
	_this.BaseUI.Init(_this, _this.impl)
	_this.BaseContainer.Init(_this)
	_this.EvState = make(map[string]func(GUI, fm.FormState))
	_this.title = ""
	_this.state = fm.FormState_Normal
	_this.border = fm.FormBorder_Default
	_this.startPos = fm.FormStart_Default
	_this.showInTaskbar = param.HideInTaskbar == false
	_this.SetSize(300, 400)
	_this.setOn()
	return _this
}

// 初始化窗体
func (_this *Form) Init() *Form {
	return _this.InitEx(br.FormParam{})
}

// 设置窗体是否置顶
func (_this *Form) SetTopMost(isTop bool) {
	_this.impl.SetTopMost(isTop)
}

// 启用无边框窗口的调整大小功能
func (_this *Form) NoneBorderResize() {
	_this.impl.NoneBorderResize()
}

// 转换为控件容器接口
func (_this *Form) toControls() br.Controls {
	return _this.impl
}

// 设置事件回调函数
func (_this *Form) setOn() {
	_this.OnState = _this.defOnState
	var bakState br.FormStateProc
	bakState = _this.impl.SetOnState(func(state fm.FormState) {
		if bakState != nil {
			bakState(state)
		}
		_this.state = state
		_this.OnState(state)
	})
	bakLoad := _this.OnLoad
	_this.OnLoad = func() {
		switch _this.startPos {
		case fm.FormStart_Screen_Center:
			scr := App.GetScreen()
			size := _this.GetBound().Rect
			x, y := scr.WorkArea.Width/2-size.Width/2, scr.WorkArea.Height/2-size.Height/2
			_this.impl.SetLocation(x, y)
		case fm.FormStart_Default:
			_this.impl.SetLocation(100, 100)
		}
		bakLoad()
	}
}

// 默认状态改变处理函数
func (_this *Form) defOnState(state fm.FormState) {
	for _, v := range _this.EvState {
		v(_this, state)
	}
}

// 设置窗体标题
func (_this *Form) SetTitle(title string) {
	_this.title = title
	_this.impl.SetTitle(_this.title)
}

// 设置窗体边框样式
func (_this *Form) SetBorderStyle(style fm.FormBorder) {
	_this.border = style
	_this.impl.SetBorderStyle(_this.border)
}

// 获取窗体边框样式
func (_this *Form) GetBorderStyle() fm.FormBorder {
	return _this.border
}

// 设置窗体状态
func (_this *Form) SetState(state fm.FormState) {
	if _this.state == state {
		return
	}
	switch state {
	case fm.FormState_Max:
		_this.impl.ShowToMax()
	case fm.FormState_Min:
		_this.impl.ShowToMin()
	case fm.FormState_Normal:
		_this.impl.Show()
	}
}

// 获取窗体状态
func (_this *Form) GetState() fm.FormState {
	return _this.state
}

// 设置窗体启动位置
func (_this *Form) SetStartPosition(pos fm.FormStart) {
	_this.startPos = pos
}

// 设置是否显示最大化按钮
func (_this *Form) SetMaximizeBox(isShow bool) {
	_this.impl.SetMaximizeBox(isShow)
}

// 设置是否显示最小化按钮
func (_this *Form) SetMinimizeBox(isShow bool) {
	_this.impl.SetMinimizeBox(isShow)
}

// 关闭窗体
func (_this *Form) Close() {
	_this.impl.Close()
}

// 设置窗体图标（从文件）
func (_this *Form) SetIcon(file string) {
	_this.impl.SetIcon(file)
}

// 设置窗体图标（从数据）
func (_this *Form) SetIconData(iconData []byte) {
	_this.impl.SetIconData(iconData)
}

// 以模态方式显示窗体
func (_this *Form) ShowDialog() {
	_this.SetStartPosition(fm.FormStart_Screen_Center)
	_this.impl.ShowDialog()
}

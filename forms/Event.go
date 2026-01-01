package forms

import (
	"time"
)

//todo 应该全是接口

// MouseEvArgs 鼠标事件参数
type MouseEvArgs struct {
	Button           MouseButtons
	X, Y, Delta      int
	IsDouble         bool
	Time             time.Time
	IsHandle         bool
	ScreenX, ScreenY int
}

// PaintEvArgs 绘制事件参数
type PaintEvArgs struct {
	Clip     Bound
	Graphics Graphics
}

// KeyEvArgs 键盘事件参数
type KeyEvArgs struct {
	Key      Keys
	Value    uintptr
	IsHandle bool
	IsSys    bool
}

// KeyPressEvArgs 按键按下事件参数
type KeyPressEvArgs struct {
	KeyChar  string
	Value    uintptr
	IsHandle bool
	IsSys    bool
}

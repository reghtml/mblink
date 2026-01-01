package controls

import (
	fm "github.com/reghtml/mblink/forms"
	br "github.com/reghtml/mblink/forms/bridge"
)

// Container 容器接口
type Container interface {
	GUI

	toControls() br.Controls
}

// Child 子控件接口
type Child interface {
	GUI

	setOwner(owner GUI)
	setParent(parent GUI)
	toControl() br.Control
	GetAnchor() fm.AnchorStyle
	SetAnchor(style fm.AnchorStyle)
}

// BaseContainer 容器基类，管理子控件的布局和锚定
type BaseContainer struct {
	Childs map[uintptr]Child

	container Container
	logs      map[uintptr]fm.Bound2
}

// 初始化容器
func (_this *BaseContainer) Init(container Container) *BaseContainer {
	_this.Childs = make(map[uintptr]Child)
	_this.logs = make(map[uintptr]fm.Bound2)
	_this.container = container
	var bakResize br.WindowResizeProc
	bakResize = container.toControls().SetOnResize(func(e fm.Rect) {
		_this.onAnchor(e)
		if bakResize != nil {
			bakResize(e)
		}
	})
	return _this
}

func (_this *BaseContainer) onAnchor(rect fm.Rect) {
	for _, n := range _this.Childs {
		b := _this.logs[n.GetHandle()]
		anc := n.GetAnchor()
		bound := n.GetBound()
		pos := bound.Point
		sz := bound.Rect
		if anc == fm.AnchorStyle_Fill {
			sz = _this.container.GetBound().Rect
			pos = fm.Point{}
		} else {
			pos = fm.Point{
				X: b.Left,
				Y: b.Top,
			}
			sz = fm.Rect{
				Width:  b.Right - b.Left,
				Height: b.Bottom - b.Top,
			}
			if anc&fm.AnchorStyle_Left != 0 {
				pos.X = b.Left
			}
			if anc&fm.AnchorStyle_Right != 0 {
				if anc&fm.AnchorStyle_Left != 0 {
					sz.Width = rect.Width - b.Left - b.Right
				} else {
					pos.X = rect.Width - b.Right - sz.Width
				}
			}
			if anc&fm.AnchorStyle_Top != 0 {
				pos.Y = b.Top
			}
			if anc&fm.AnchorStyle_Bottom != 0 {
				if anc&fm.AnchorStyle_Top != 0 {
					sz.Height = rect.Height - b.Top - b.Bottom
				} else {
					pos.Y = rect.Height - b.Bottom - sz.Height
				}
			}
		}

		n.SetSize(sz.Width, sz.Height)
		n.SetLocation(pos.X, pos.Y)
	}
}

// 添加子控件
func (_this *BaseContainer) AddChild(child Child) {
	if _, ok := _this.Childs[child.GetHandle()]; ok == false {
		_this.container.toControls().AddControl(child.toControl())
		bn := child.GetBound()
		pb := _this.container.GetBound()
		rect := fm.Bound2{
			Left:   bn.X,
			Top:    bn.Y,
			Right:  pb.Width - bn.Width - bn.X,
			Bottom: pb.Height - bn.Height - bn.Y,
		}
		child.setParent(_this.container)
		ow := _this.container.GetOwner()
		if ow == nil {
			ow = _this.container
		}
		child.setOwner(ow)
		_this.logs[child.GetHandle()] = rect
		_this.Childs[child.GetHandle()] = child
	}
}

// 移除子控件
func (_this *BaseContainer) RemoveChild(child Child) {
	if _, ok := _this.Childs[child.GetHandle()]; ok {
		_this.container.toControls().RemoveControl(child.toControl())
		child.setParent(nil)
		child.setOwner(nil)
		delete(_this.Childs, child.GetHandle())
		delete(_this.logs, child.GetHandle())
	}
}

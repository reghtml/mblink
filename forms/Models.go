package forms

import (
	"image"
)

// Point 点坐标结构
type Point struct {
	X, Y int
}

// 判断两个点是否相等
func (_this Point) IsEqual(point Point) bool {
	return _this.X == point.X && _this.Y == point.Y
}

// Rect 矩形尺寸结构
type Rect struct {
	Width, Height int
}

// 判断两个矩形尺寸是否相等
func (_this Rect) IsEqual(rect Rect) bool {
	return _this.Width == rect.Width && _this.Height == rect.Height
}

// 判断矩形尺寸是否为空
func (_this Rect) IsEmpty() bool {
	return _this.Width == 0 || _this.Height == 0
}

// Bound 边界结构，包含位置和尺寸
type Bound struct {
	Point
	Rect
}

// Bound2 边界结构，使用上下左右四个值定义
type Bound2 struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

// Screen 屏幕信息结构
type Screen struct {
	Full     Rect
	WorkArea Rect
}

// Graphics 图形绘制接口
type Graphics interface {
	GetHandle() uintptr
	DrawImage(src *image.RGBA, xSrc, ySrc, width, height, xDst, yDst int) Graphics
	Close()
}

// MsgBoxParam 消息框参数
type MsgBoxParam struct {
	Title  string
	Text   string
	Icon   MsgBoxIcon
	Button MsgBoxButton
}

package bridge

// Controls 容器控件接口，可以添加和移除子控件
type Controls interface {
	Control

	AddControl(control Control)
	RemoveControl(control Control)
	GetChilds() []Control
}

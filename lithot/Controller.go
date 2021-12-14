package lithot

type Controller interface {
	Build(lithot *Lithot) //参数和方法名必须一致
	Name() string
}
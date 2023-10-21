package cp_mq

type ILog interface {
	Debug(str string, isForce ...bool)
	Info(str string, isForce ...bool)
	Error(str string, isForce ...bool)
	Critical(str string, isForce ...bool)
}

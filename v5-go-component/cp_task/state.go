package cp_task

import "fmt"

const (
	TASK_STATE_STOP     TaskState = iota //任务停止中
	TASK_STATE_STOPPING                  //任务正在停止中
	TASK_STATE_WAIT                      //等待启动中，只有设定StartTime时，才会进入这个状态
	TASK_STATE_RUN                       //任务执行中
	TASK_STATE_SLEEP                     //任务休眠，等待下一次执行
)

type TaskState int

func (t TaskState) String() string {
	return stateToString(t)
}

//获取状态文字说明
func stateToString(s TaskState) string {
	str := ""
	switch s {
	case TASK_STATE_STOP:
		str = "停止"
	case TASK_STATE_STOPPING:
		str = "停止中"
	case TASK_STATE_WAIT:
		str = "等待"
	case TASK_STATE_RUN:
		str = "执行"
	case TASK_STATE_SLEEP:
		str = "休眠"
	default:
		str = fmt.Sprintf("未知[%v]", s)
	}

	return str
}

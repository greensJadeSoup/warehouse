package cp_task

import (
	"errors"
	"fmt"
	"runtime"
	"time"
	"warehouse/v5-go-component/cp_obj"
)

// Task 任务主体
type Task struct {
	Name      string      //任务名称
	IsSingle  bool        //是否单例运行
	StartTime cp_obj.Datetime //小于等于当前时间马上执行

	OnInterval func(t *Task) time.Duration //事件：必须，任务循环间隔时间
	OnRun      func(t *Task)               //事件：必须，任务执行主体
	OnPanic    func(t *Task, err error)    //事件：必须，发生异常时

	lastRunTime cp_obj.Datetime //上一次执行时间
	nextRunTime cp_obj.Datetime //下一次执行时间
	runCount    int         //运行次数累计
	panicCount  int         //异常次数
	state       TaskState   //任务状态
	isInit      bool
	stopChan    chan struct{} //任务中止通知

	data map[string]interface{}
}

// Init 初始化Task，并返回StartTime与当前时间的时差秒数
func (t *Task) Init() error {
	if t.isInit {
		return nil
	}

	//检查名字设定
	if t.Name == "" {
		return errors.New("Name没有设置")
	}

	//检查时间循环
	if t.OnInterval == nil {
		return errors.New("OnInterval事件没有设置")
	}

	//检查Run函数是否设置
	if t.OnRun == nil {
		return errors.New("OnRun事件没有设置")
	}

	if t.OnPanic == nil {
		return errors.New("OnPanic事件没有设置")
	}

	var startTime cp_obj.Datetime
	if t.StartTime == startTime {
		t.StartTime = cp_obj.NewDatetime()
	}

	t.isInit = true
	t.stopChan = make(chan struct{})
	t.data = make(map[string]interface{})
	return nil
}

// Start 启动任务
func (t *Task) Start() error {
	if !t.isInit {
		err := t.Init()
		if err != nil {
			return err
		}
	}

	if t.state != TASK_STATE_STOP {
		//运行中的任务，直接返回
		return nil
	}

	go func(tt *Task) {
		tt.state = TASK_STATE_WAIT

		//异常捕获
		defer recoverTaskPanic(tt)

		//指定时间启动
		d := t.StartTime.Time().Unix() - time.Now().Unix()

		if d > 0 {
			select {
			case <-tt.stopChan:
				goto TaskStop
			case <-time.After(time.Duration(d) * time.Second):
			}
		}

		for {
			tt.state = TASK_STATE_RUN

			//统计运行记录
			i := tt.OnInterval(tt)
			tt.runCount++
			tt.lastRunTime.SetTime(time.Now())
			tt.nextRunTime.SetTime(time.Now().Add(i))

			runStart := time.Now().UnixNano()
			tt.OnRun(tt)
			if tt.state == TASK_STATE_STOPPING {
				<-tt.stopChan
				goto TaskStop
			}
			runEnd := time.Now().UnixNano()
			runPay := runEnd - runStart

			tt.state = TASK_STATE_SLEEP
			select {
			case <-tt.stopChan:
				goto TaskStop
			case <-time.After(i - time.Duration(runPay)):
			}
		}

	TaskStop:
		t.state = TASK_STATE_STOP
	}(t)

	return nil
}

// Stop 中止任务
func (t *Task) Stop() {
	if t.state == TASK_STATE_STOPPING || t.state == TASK_STATE_STOP {
		return
	}
	t.state = TASK_STATE_STOPPING

	go func() {
		t.stopChan <- struct{}{}
	}()
}

// Wait 等待任务中止完成
func (t *Task) Wait() {
	for t.state != TASK_STATE_STOP {
		time.Sleep(100 * time.Millisecond)
	}
}

// Data 获取缓存对象
func (t *Task) Data() map[string]interface{} {
	return t.data
}

// LastRunTime 上一次执行时间
func (t *Task) LastRunTime() cp_obj.Datetime {
	return t.lastRunTime
}

// NextRunTime 下一次执行时间
func (t *Task) NextRunTime() cp_obj.Datetime {
	return t.nextRunTime
}

// RunCount 任务执行周期内的累计执行次数
func (t *Task) RunCount() int {
	return t.runCount
}

// PanicCount 执行周期内的异常次数
func (t *Task) PanicCount() int {
	return t.panicCount
}

// State 当前任务状态
func (t *Task) State() TaskState {
	return t.state
}

func recoverTaskPanic(t *Task) {
	if e := recover(); e != nil {
		t.state = TASK_STATE_STOP
		t.panicCount++

		errorInfo := fmt.Sprintf("%s任务执行故障：%v\n故障堆栈：", t.Name, e)
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			} else {
				errorInfo += "\n"
			}
			errorInfo += fmt.Sprintf("%v %v", file, line)
		}

		t.OnPanic(t, errors.New(errorInfo))
	}
}

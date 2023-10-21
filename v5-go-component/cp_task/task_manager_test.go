package cp_task

import (
	"testing"
	"time"
)

func TestTaskManager(t *testing.T) {
	tm := &TaskManager{}

	t1 := &Task{
		Name:     "单例任务",
		IsSingle: true,
		OnInterval: func(tt1 *Task) time.Duration {
			return 1 * time.Second
		},
		OnRun: func(tt1 *Task) {
			t.Log("-单例任务运行：", tt1.RunCount())
		},
		OnPanic: func(tt1 *Task, e1 error) {
			t.Error("单例任务异常：", e1)
		},
	}

	t2 := &Task{
		Name: "多例任务",
		OnInterval: func(tt2 *Task) time.Duration {
			return 1 * time.Second
		},
		OnRun: func(tt2 *Task) {
			t.Log("多例任务运行：", tt2.RunCount())
		},
		OnPanic: func(tt2 *Task, e2 error) {
			t.Error("多例任务异常：", e2)
		},
	}

	err := tm.Add(t1)
	if err != nil {
		t.Error("单例任务添加失败：", err)
	}

	err = tm.Add(t2)
	if err != nil {
		t.Error("多例任务添加失败：", err)
	}

	err = tm.StartSingle()
	if err != nil {
		t.Error("单例任务启动失败：", err)
	}

	err = tm.StartMultiple()
	if err != nil {
		t.Error("多例任务启动失败：", err)
	}

	time.Sleep(5 * time.Second)

	t.Log("停止单例任务")
	tm.StopSingle()
	tm.WaitSingle()
	t.Log("停止单例任务-成功")
	time.Sleep(2 * time.Second)
	t.Log("重新启动单例任务")

	err = tm.StartSingle()
	if err != nil {
		t.Error("单例任务重新启动失败：", err)
	}

	time.Sleep(5 * time.Second)
	t.Log("停止所有任务")
	tm.Stop()
	tm.Wait()
	t.Log("结束")
}

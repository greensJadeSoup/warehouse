package cp_task

import (
	"fmt"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	task := &Task{
		Name:     "testTask",
		IsSingle: true,
		OnInterval: func(tt *Task) time.Duration {
			return 1 * time.Second
		},
		OnRun: func(tt *Task) {
			cp_log.Info("任务执行", tt.RunCount())

			if tt.RunCount() == 2 {
				panic("测试异常")
			}

			if tt.RunCount() > 5 {
				if tt.IsSingle {
					cp_log.Info("停止任务")
					tt.Stop()
				}
			}
		},
		OnPanic: func(tt *Task, err error) {
			cp_log.Info("运行中异常：", err)
			//恢复任务运行
			err = tt.Start()
			if err != nil {
				cp_log.Info("异常后重启失败：", err)
			}
		},
	}

	err := task.Start()
	if err != nil {
		cp_log.Info("任务启动失败：", err)
	}

	time.Sleep(2 * time.Second)
	cp_log.Info("等待结束")
	task.Wait()
	cp_log.Info("结束")
}

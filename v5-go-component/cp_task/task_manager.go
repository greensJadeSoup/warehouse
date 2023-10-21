package cp_task

import (
	"fmt"
	"go.uber.org/zap"
	"sync"
	"warehouse/v5-go-component/cp_log"
)

// TaskManager 任务中心
type TaskManager struct {
	tasks []*Task
}

var tm TaskManager

func GetTaskManager() *TaskManager {
	return &tm
}

// Add 添加任务
func (m *TaskManager) Add(t *Task) error {
	err := t.Init()
	if err != nil {
		return err
	}

	m.tasks = append(m.tasks, t)
	cp_log.Info("任务添加成功", zap.String("name", t.Name))
	return nil
}

// Start 启动所有任务
func (m *TaskManager) Start() error {
	for i, l := 0, len(m.tasks); i < l; i++ {
		err := m.tasks[i].Start()
		if err != nil {
			m.Stop()
			return err
		}
	}

	return nil
}

// StartMultiple 启动多例任务
func (m *TaskManager) StartMultiple() error {
	for i, l := 0, len(m.tasks); i < l; i++ {
		if m.tasks[i].IsSingle {
			continue
		}

		err := m.tasks[i].Start()
		if err != nil {
			m.StopMultiple()
			return err
		}
	}

	return nil
}

// StartSingle 启动单例任务
func (m *TaskManager) StartSingle() error {
	for i, l := 0, len(m.tasks); i < l; i++ {
		if !m.tasks[i].IsSingle {
			continue
		}

		err := m.tasks[i].Start()
		if err != nil {
			m.StopSingle()
			return err
		}
	}

	return nil
}

// Stop 停止所有任务
func (m *TaskManager) Stop() {
	for i, l := 0, len(m.tasks); i < l; i++ {
		m.tasks[i].Stop()
	}
}

// StopMultiple 停止多例任务
func (m *TaskManager) StopMultiple() {
	for i, l := 0, len(m.tasks); i < l; i++ {
		if m.tasks[i].IsSingle {
			continue
		}

		m.tasks[i].Stop()
	}
}

// StopSingle 停止单例任务
func (m *TaskManager) StopSingle() {
	for i, l := 0, len(m.tasks); i < l; i++ {
		if !m.tasks[i].IsSingle {
			continue
		}

		m.tasks[i].Stop()
	}
}

// Wait 等待所有任务停止
func (m *TaskManager) Wait() {
	var wg sync.WaitGroup

	for i, l := 0, len(m.tasks); i < l; i++ {
		wg.Add(1)

		go func(index int) {
			m.tasks[index].Wait()
			wg.Done()
		}(i)

	}

	wg.Wait()
}

// WaitMultiple 等待所有多例任务停止
func (m *TaskManager) WaitMultiple() {
	var wg sync.WaitGroup

	for i, l := 0, len(m.tasks); i < l; i++ {
		if m.tasks[i].IsSingle {
			continue
		}

		wg.Add(1)

		go func(index int) {
			m.tasks[index].Wait()
			wg.Done()
		}(i)

	}

	wg.Wait()
}

// WaitSingle 等待所有单例任务停止
func (m *TaskManager) WaitSingle() {
	var wg sync.WaitGroup

	for i, l := 0, len(m.tasks); i < l; i++ {
		if !m.tasks[i].IsSingle {
			continue
		}

		wg.Add(1)

		go func(index int) {
			m.tasks[index].Wait()
			wg.Done()
		}(i)

	}

	wg.Wait()
}

// Info 任务中心运行信息
func (m *TaskManager) Info() map[string]interface{} {
	info := make(map[string]interface{}, 2)

	tasksLen := len(m.tasks)
	items := make([]map[string]interface{}, tasksLen)

	for i := 0; i < tasksLen; i++ {
		item := make(map[string]interface{}, 8)
		item["Name"] = m.tasks[i].Name
		item["IsSingle"] = m.tasks[i].IsSingle
		item["StartTime"] = m.tasks[i].StartTime
		item["Interval"] = fmt.Sprintf("%v", m.tasks[i].OnInterval(m.tasks[i]))
		item["LastRunTime"] = m.tasks[i].LastRunTime()
		item["NextRunTime"] = m.tasks[i].NextRunTime()
		item["RunCount"] = m.tasks[i].RunCount()
		item["PanicCount"] = m.tasks[i].PanicCount()
		item["State"] = m.tasks[i].State().String()

		items[i] = item
	}

	info["Items"] = items

	return info
}

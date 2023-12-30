package clash

import (
	"fmt"
	"testing"
)

func Init() {
	if err := SetSecretFromEnv("clash-api-secret"); err != nil {
		panic(err)
	}
}

func TestGetLogs(t *testing.T) {
	Init()
	// 调用 GetLogs 获取日志 channel
	logChan, err := GetLogs(LevelDebug)
	if err != nil {
		t.Errorf("Error retrieving logs: %s", err)
		return
	}

	receivedLogs := make([]*Log, 0)

	// 使用 select 来等待日志写入 channel，并获取前三条日志内容
	for i := 0; i < 3; i++ {
		select {
		case log, ok := <-logChan:
			if !ok {
				t.Error("Channel closed unexpectedly")
				return
			}
			receivedLogs = append(receivedLogs, log)
		}
	}

	// 打印前三条日志内容
	fmt.Println("Received logs:")
	for _, log := range receivedLogs {
		fmt.Println(log)
	}
}

func TestGetMemory(t *testing.T) {
	Init()
	// 调用 GetLogs 获取日志 channel
	count := 0
	exitChan := make(chan struct{}) // 创建一个退出通道
	_ = GetMemory(func(memory *Memory) (stop bool) {
		count++
		fmt.Printf("%+v\n", memory)
		if count > 3 {
			close(exitChan) // 发送退出通道信号
			return true
		}
		return false
	})

	for {
		select {
		case <-exitChan: // 当接收到退出通道信号时退出循环
			fmt.Println("Exiting loop")
			return
		}
	}
}

func TestRestart(t *testing.T) {
	Init()
	Restart()
}

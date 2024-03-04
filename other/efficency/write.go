package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// SessionRequest 定义创建会话时的请求体结构
type SessionRequest struct {
	KeyName   string `json:"keyName"`
	OwnerName string `json:"ownerName"`
	TaskName  string `json:"taskName"`
	IpnsName  string `json:"ipnsName"`
	LocalPath string `json:"localPath"`
}

// SessionResponse 定义用于解析会话创建响应的结构体
type SessionResponse struct {
	SessionID string `json:"session_id"`
}

// WriteDataRequest 定义写入数据请求的结构体
type WriteDataRequest struct {
	Data   string `json:"data"`
	Height int    `json:"height"`
}

func main() {
	// 定义并发执行的次数
	concurrencyLevel := 1000
	maxConcurrency := 10

	// 创建一个缓冲channel来限制并发数量
	concurrencySemaphore := make(chan struct{}, maxConcurrency)

	// 创建WaitGroup以等待所有goroutine完成
	var wg sync.WaitGroup
	wg.Add(concurrencyLevel)

	// 记录开始时间
	startTime := time.Now()

	for i := 0; i < concurrencyLevel; i++ {
		concurrencySemaphore <- struct{}{}

		go func() {
			defer wg.Done() // 在goroutine完成时通知WaitGroup
			defer func() { <-concurrencySemaphore }()

			// 创建SessionRequest实例
			sessionReq := SessionRequest{
				KeyName:   "image",
				OwnerName: "pinge",
				TaskName:  "api",
				IpnsName:  "k51qzi5uqu5did01y4bfh94mbd1olkqyyyj1hqhtrrqsxh97funiqyod9l2dx8",
				LocalPath: "/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile",
			}

			// 执行过程
			err := runSessionAndWriteDataProcess(sessionReq, i)
			if err != nil {
				fmt.Println("过程执行错误:", err)
			}
		}()
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 计算并打印总耗时
	elapsedTime := time.Since(startTime)
	fmt.Printf("总耗时: %s\n", elapsedTime)
}

// runSessionAndWriteDataProcess 执行创建会话并写入数据的过程
func runSessionAndWriteDataProcess(sessionReq SessionRequest, i int) error {
	// 创建会话
	sessionID, err := createSession(sessionReq)
	if err != nil {
		return err
	}
	// 使用会话ID写入数据
	err = writeData(sessionID, "test", i)
	if err != nil {
		return err
	}

	return nil
}

// createSession 发送POST请求到会话创建接口
func createSession(sessionReq SessionRequest) (string, error) {
	url := "http://127.0.0.1:23333/image/create"
	jsonData, err := json.Marshal(sessionReq)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var sessionResp SessionResponse
	if err := json.Unmarshal(body, &sessionResp); err != nil {
		return "", err
	}

	return sessionResp.SessionID, nil
}

// writeData 使用会话ID和提供的数据发送POST请求到数据写入接口
func writeData(sessionID, data string, height int) error {
	url := fmt.Sprintf("http://127.0.0.1:23333/image/addstring?session_id=%s", sessionID)
	requestData := WriteDataRequest{
		Data:   data,
		Height: height,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

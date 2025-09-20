package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// PayloadCollection 请求数据
type PayloadCollection struct {
	Token    string    `json:"token"`
	Payloads []Payload `json:"data"`
}

// Payload 单个任务，包含 HTTP 响应上下文
type Payload struct {
	Data string
	W    http.ResponseWriter
	R    *http.Request
}

// Process 执行任务逻辑并响应 HTTP
func (p *Payload) Process() {
	// 模拟耗时操作
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Processed payload: %v\n", p.Data)

	// 处理完成后直接响应 HTTP
	if p.W != nil {
		p.W.Write([]byte(fmt.Sprintf("Payload %v processed\n", p.Data)))
	}
}

// Job 封装任务
type Job struct {
	Payload Payload
}

// WorkerPool 协程池
type WorkerPool struct {
	taskQueue chan Job
}

// NewWorkerPool 创建协程池
func NewWorkerPool(queueSize int) *WorkerPool {
	return &WorkerPool{
		taskQueue: make(chan Job, queueSize),
	}
}

// StartWorkers 启动 n 个 worker
func (p *WorkerPool) StartWorkers(count int) {
	for i := 0; i < count; i++ {
		go func() {
			for job := range p.taskQueue {
				job.Payload.Process() // worker 处理任务，并响应请求
			}
		}()
	}
}

// Submit 提交任务到协程池
func (p *WorkerPool) Submit(job Job) {
	p.taskQueue <- job
}

var (
	pool   *WorkerPool
	taskID int64
)

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 每个请求启动 goroutine 解析请求
	go func() {
		var content PayloadCollection
		err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&content)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid request\n"))
			return
		}
		for _, payload := range content.Payloads {
			id := atomic.AddInt64(&taskID, 1)
			fmt.Printf("Submitting task %d\n", id)

			// 封装 HTTP 响应上下文到 Payload
			payload.W = w
			payload.R = r

			pool.Submit(Job{Payload: payload}) // 提交给协程池处理
		}
	}()
}

func main() {
	queueSize := 1000 // 能容纳多条
	pool = NewWorkerPool(queueSize)

	workerCount := 10 // 启动多少个子线程去执行任务
	pool.StartWorkers(workerCount)

	http.HandleFunc("/payload", payloadHandler)
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

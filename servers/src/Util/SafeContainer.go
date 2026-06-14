package Util

import (
	"context"
)

// Action 定义操作闭包
type Action[T any] func(data *[]T)

// SafeContainer 绑定上下文的安全容器
type SafeContainer[T any] struct {
	data     []T
	actionCh chan Action[T]
}

// NewSafeContainer 构造函数
// ctx: 传入父级的 Context（比如对局的 Context）
func NewSafeContainer[T any](ctx context.Context, bufferSize int) *SafeContainer[T] {
	s := &SafeContainer[T]{
		data:     make([]T, 0),
		actionCh: make(chan Action[T], bufferSize),
	}

	// 启动管理员协程
	go s.run(ctx)

	return s
}

func (s *SafeContainer[T]) run(ctx context.Context) {
	for {
		select {
		case action := <-s.actionCh:
			action(&s.data)
		case <-ctx.Done():
			return
		}
	}
}

// Push 异步添加
func (s *SafeContainer[T]) Push(item T) {
	// 注意：如果 ctx 已经结束，再往里塞可能会阻塞或报错，可以加个 select 保护

	s.actionCh <- func(data *[]T) {
		*data = append(*data, item)
	}
}

// Do 执行复合原子操作（同步等待）
func (s *SafeContainer[T]) Do(f func(data *[]T)) {
	done := make(chan struct{})
	s.actionCh <- func(data *[]T) {
		f(data)
		close(done)
	}
	<-done
}

// GetLen 获取长度（同步）
func (s *SafeContainer[T]) GetLen() int {
	resCh := make(chan int)
	s.actionCh <- func(data *[]T) {
		resCh <- len(*data)
	}
	return <-resCh
}

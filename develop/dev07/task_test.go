package main

import (
	"testing"
	"time"
)

// helper функция для создания каналов, которые закрываются через заданный интервал времени
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

// Тест, проверяющий, что объединённый канал закрывается при закрытии одного из входных каналов.
func TestOr(t *testing.T) {
	t.Run("one channel closes", func(t *testing.T) {
		start := time.Now()
		<-or(
			sig(2*time.Hour),
			sig(5*time.Minute),
			sig(1*time.Second),
			sig(1*time.Hour),
			sig(1*time.Minute),
		)

		duration := time.Since(start)
		if duration < 1*time.Second || duration > 2*time.Second {
			t.Errorf("Expected close after around 1 second, but got %v", duration)
		}
	})

	t.Run("multiple channels close immediately", func(t *testing.T) {
		start := time.Now()
		<-or(
			sig(0),
			sig(0),
			sig(0),
		)

		duration := time.Since(start)
		if duration > 1*time.Millisecond {
			t.Errorf("Expected immediate close, but got %v", duration)
		}
	})

	t.Run("no channels provided", func(t *testing.T) {
		merged := or()
		select {
		case <-merged:
			t.Errorf("Expected merged channel to remain open, but it closed")
		case <-time.After(100 * time.Millisecond):
			// Expected behavior: the merged channel should remain open
		}
	})

	t.Run("longest channel closes last", func(t *testing.T) {
		start := time.Now()
		<-or(
			sig(500*time.Millisecond),
			sig(1*time.Second),
			sig(150*time.Millisecond),
		)

		duration := time.Since(start)
		if duration < 150*time.Millisecond || duration > 200*time.Millisecond {
			t.Errorf("Expected close after around 150 milliseconds, but got %v", duration)
		}
	})
}

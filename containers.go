package containers

import (
	linkedmap "github.com/elvislee1126/go-containers/linked-map"
	rollingwindow "github.com/elvislee1126/go-containers/rolling-window"
)

// 滑动窗口
func NewRollingWindow[K comparable, V any](opts ...rollingwindow.RollingWindowOption) *rollingwindow.RollingWindow[K, V] {
	return rollingwindow.New[K, V](opts...)
}

// linked hash map
func NewLinkedMap[K comparable, V any]() *linkedmap.LinkedMap[K, V] {
	return linkedmap.New[K, V]()
}

package linkedmap

import (
	"sync"

	"container/list"
)

type LinkedMap[K comparable, V any] struct {
	lock sync.RWMutex
	l    *list.List
	m    map[K]*node[K, V]
}

type node[K comparable, V any] struct {
	key K
	val V
	ele *list.Element
}

type LinkedMapRanger[K comparable, V any] func(k K, v V) bool

func New[K comparable, V any]() *LinkedMap[K, V] {
	return &LinkedMap[K, V]{
		lock: sync.RWMutex{},
		l:    list.New(),
		m:    make(map[K]*node[K, V]),
	}
}

// 设置 k 的键值为 v
func (lm *LinkedMap[K, V]) Store(k K, v V) {
	lm.lock.Lock()
	defer lm.lock.Unlock()
	if node := lm.m[k]; node != nil {
		node.val = v
		return
	}
	node := &node[K, V]{
		key: k,
		val: v,
	}
	node.ele = lm.l.PushBack(node)
	lm.m[k] = node
}

// 读取 k 的键值
func (lm *LinkedMap[K, V]) Load(k K) (V, bool) {
	lm.lock.RLock()
	defer lm.lock.RUnlock()
	if node := lm.m[k]; node != nil {
		return node.val, true
	}
	var zero V
	return zero, false
}

// 若 key 的键值存在，返回 (原始键值, false)
// 否则，键值设置为 v，返回 (v, true)
func (lm *LinkedMap[K, V]) LoadOrStore(k K, v V) (V, bool) {
	lm.lock.Lock()
	defer lm.lock.Unlock()
	if node := lm.m[k]; node != nil {
		return node.val, false
	}
	node := &node[K, V]{
		key: k,
		val: v,
	}
	node.ele = lm.l.PushBack(node)
	lm.m[k] = node
	return node.val, true
}

// 删除元素
func (lm *LinkedMap[K, V]) Delete(k K) (V, bool) {
	lm.lock.RLock()
	node := lm.m[k]
	if node == nil {
		lm.lock.RUnlock()
		var zero V
		return zero, false
	}
	lm.lock.RUnlock()
	lm.lock.Lock()
	defer lm.lock.Unlock()
	lm.l.Remove(node.ele)
	delete(lm.m, k)
	return node.val, true
}

// 返回按插入时间倒序排序的键名
func (lm *LinkedMap[K, V]) Keys() []K {
	keys := make([]K, 0, lm.l.Len())
	lm.Range(func(k K, v V) bool {
		keys = append(keys, k)
		return false
	})
	return keys
}

// 按插入时间倒序遍历
func (lm *LinkedMap[K, V]) Range(fn LinkedMapRanger[K, V]) {
	lm.lock.RLock()
	defer lm.lock.RUnlock()
	for cursor := lm.l.Back(); cursor != nil; cursor = cursor.Prev() {
		node := cursor.Value.(*node[K, V])
		if isBreak := fn(node.key, node.val); isBreak {
			break
		}
	}
}

// 返回元素数量
func (lm *LinkedMap[K, V]) Len() int {
	return lm.l.Len()
}

// 返回目前最老的元素
func (lm *LinkedMap[K, V]) Oldest() (V, bool) {
	var zero V
	lm.lock.RLock()
	defer lm.lock.RUnlock()
	ele := lm.l.Front()
	if ele == nil {
		return zero, false
	}
	n := ele.Value.(*node[K, V])
	return n.val, true
}

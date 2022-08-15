package rollingwindow

import (
	"errors"
	"log"
	"time"

	linkedmap "github.com/elvislee1126/go-containers/linked-map"
)

type PositionTag int

const (
	LeadOfWindowTag PositionTag = iota
	InWindowTag
	BehindOfWindowTag
)

var (
	ErrElementBehindOfWindow = errors.New("element is behind of current window")
	ErrElementLeafOfWindow   = errors.New("element is lead of current window")
)

type RollingWindow[K comparable, V any] struct {
	verbose bool
	logger  *log.Logger

	// 窗口 slot 数量
	slotAmount int64

	// 窗口 slot 容量，即这个 slot 时间片大小
	slotSize time.Duration

	// 提供系统运行时时间
	timeProvider TimeProvider

	// 释放窗口外 slot 的间隔
	drainTimer    *time.Timer
	drainInterval time.Duration

	// 当前所有 slot，包含窗口内何窗口外的
	slots *linkedmap.LinkedMap[int64, *windowSlot[K, V]]

	drainRequestChan chan *WindowPosition
}

type WindowPosition struct {
	SlotIdx          int64
	WindowIdx        [2]int64
	RelativePosition PositionTag
}

func NewRollingWindow[K comparable, V any](opts ...RollingWindowOption) *RollingWindow[K, V] {
	o := DefaultNewOptions()
	for _, cfg := range opts {
		o = cfg.apply(o)
	}
	win := &RollingWindow[K, V]{
		slotAmount:       o.SlotAmount,
		slotSize:         o.SlotSize,
		timeProvider:     o.TimeProvider,
		slots:            linkedmap.New[int64, *windowSlot[K, V]](),
		logger:           log.Default(),
		drainRequestChan: make(chan *WindowPosition, 1),
		drainTimer:       time.NewTimer(o.DrainInterval),
		drainInterval:    o.DrainInterval,
	}
	win.startDrain()
	return win
}

func (r *RollingWindow[K, V]) GetWindowPosition(eleTime *time.Time) WindowPosition {
	sysNowTime := r.timeProvider()
	windowRightIdx := sysNowTime.UnixMilli() / r.slotSize.Milliseconds()
	windowLeftIdx := windowRightIdx - r.slotAmount
	position := WindowPosition{
		WindowIdx: [2]int64{windowLeftIdx, windowRightIdx},
	}
	if eleTime != nil {
		slotIdx := eleTime.UnixMilli() / r.slotSize.Milliseconds()
		var relativePosition PositionTag
		if slotIdx < windowLeftIdx {
			relativePosition = BehindOfWindowTag
		} else if slotIdx > windowRightIdx {
			relativePosition = LeadOfWindowTag
		} else {
			relativePosition = InWindowTag
		}
		position.RelativePosition = relativePosition
		position.SlotIdx = slotIdx
	}

	return position
}

func (r *RollingWindow[K, V]) SetWithTime(key K, value V, t time.Time) error {
	positionInfo := r.GetWindowPosition(&t)
	if positionInfo.RelativePosition == BehindOfWindowTag {
		return ErrElementBehindOfWindow
	}
	if positionInfo.RelativePosition == LeadOfWindowTag {
		return ErrElementLeafOfWindow
	}
	slot, isNew := r.slots.LoadOrStore(positionInfo.SlotIdx, newWindowSlot[K, V](positionInfo.SlotIdx))
	slot.Store(key, value)
	// 窗口前移
	// 新建 slot
	if isNew {
		if r.verbose {
			r.logger.Printf("create new slot: slot_idx=%d", positionInfo.SlotIdx)
		}
		// todo 回收旧的 slot
		// r.Drain(positionInfo)
	}
	return nil
}

func (r *RollingWindow[K, V]) Set(key K, value V) error {
	return r.SetWithTime(key, value, time.Now())
}

func (r *RollingWindow[K, V]) Get(key K, opts ...RollingWindowGetElementOption) (val V, ok bool) {
	var zeroVal V
	val = zeroVal

	cfg := getConfig{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.CurrentWindow {
		position := r.GetWindowPosition(nil)
		for i := position.WindowIdx[0]; i <= position.WindowIdx[1]; i++ {
			slot, _ := r.slots.Load(i)
			if slot == nil {
				continue
			}
			if valInSlot, ok2 := slot.Load(key); ok2 {
				val = valInSlot
				ok = true
				break
			}
		}
		return
	}
	return
}

func (r *RollingWindow[K, V]) Size() int {
	eleSize := 0
	r.slots.Range(func(k int64, v *windowSlot[K, V]) bool {
		eleSize += v.Size()
		return false
	})
	return eleSize
}

func (r *RollingWindow[K, V]) doDrain(position *WindowPosition) {
	oldestSlot, ok := r.slots.Oldest()
	if !ok {
		return
	}
	for slotIdx := oldestSlot.slotIdx; slotIdx <= position.WindowIdx[0]; slotIdx++ {
		slot, loaded := r.slots.Delete(slotIdx)
		if !loaded {
			continue
		}
		if r.verbose {
			r.logger.Printf("drain slot: slot_idx=%d  sloat_ele_amount=%d window_left=%d", slot.slotIdx, slot.Size(), position.WindowIdx[0])
		}
	}
}

func (r *RollingWindow[K, V]) startDrain() {
	if r.verbose {
		r.logger.Printf("定时清理窗口外 slot：interval=%v", r.drainInterval.String())
	}
	timer := time.NewTimer(r.drainInterval)
	go func() {
		for range timer.C {
			now := r.timeProvider()
			pos := r.GetWindowPosition(&now)
			r.drainRequestChan <- &pos
			timer.Reset(r.drainInterval)
		}
		if r.verbose {
			r.logger.Printf("drain timer closed")
		}
	}()
	go func() {
		for req := range r.drainRequestChan {
			r.doDrain(req)
		}
		if r.verbose {
			r.logger.Printf("drain worker closed")
		}
	}()
}

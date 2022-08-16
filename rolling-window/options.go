package rollingwindow

import (
	"time"
)

type TimeProvider func() time.Time

var (
	DefaultTimeProvider TimeProvider = func() time.Time { return time.Now() }
)

type config struct {
	Verbose       bool
	DrainInterval time.Duration
	SlotAmount    int64
	SlotSize      time.Duration
	TimeProvider  TimeProvider
}

type RollingWindowOption interface {
	apply(config) config
}

type rollingWindowOptionFunc func(config) config

func (fn rollingWindowOptionFunc) apply(cfg config) config {
	return fn(cfg)
}

func DefaultNewOptions() config {
	return config{
		DrainInterval: time.Second * 10,
		SlotAmount:    10,
		SlotSize:      time.Second,
		TimeProvider:  DefaultTimeProvider,
	}
}

func WithSlotAmountAndSize(amount int64, size time.Duration) RollingWindowOption {
	return rollingWindowOptionFunc(func(cfg config) config {
		cfg.SlotAmount = amount
		cfg.SlotSize = size
		return cfg
	})
}

func WithTimeProvider(fn TimeProvider) RollingWindowOption {
	return rollingWindowOptionFunc(func(cfg config) config {
		cfg.TimeProvider = fn
		return cfg
	})
}

func WithDrainInterval(d time.Duration) RollingWindowOption {
	return rollingWindowOptionFunc(func(cfg config) config {
		cfg.DrainInterval = d
		return cfg
	})
}

func WithVerbose() RollingWindowOption {
	return rollingWindowOptionFunc(func(cfg config) config {
		cfg.Verbose = true
		return cfg
	})
}

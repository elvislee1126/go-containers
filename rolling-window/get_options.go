package rollingwindow

import "time"

type getConfig struct {
	CurrentWindow bool
	WindowLeft    *time.Time
	WindowRight   *time.Time
}

type RollingWindowGetElementOption interface {
	apply(getConfig) getConfig
}

type rollingWindowGetElementOptionFunc func(getConfig) getConfig

func (fn rollingWindowGetElementOptionFunc) apply(cfg getConfig) getConfig {
	return fn(cfg)
}

func WithCurrentWindow() RollingWindowGetElementOption {
	return rollingWindowGetElementOptionFunc(func(gc getConfig) getConfig {
		gc.CurrentWindow = true
		return gc
	})
}

func WithAfter(after time.Time) RollingWindowGetElementOption {
	return rollingWindowGetElementOptionFunc(func(gc getConfig) getConfig {
		gc.WindowLeft = &after
		return gc
	})
}

func WithBefore(before time.Time) RollingWindowGetElementOption {
	return rollingWindowGetElementOptionFunc(func(gc getConfig) getConfig {
		gc.WindowRight = &before
		return gc
	})
}

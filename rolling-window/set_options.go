package rollingwindow

import "time"

type setConfig struct {
	T     *time.Time
	SetNX bool
}

func DefaultSetOptions() setConfig {
	now := time.Now()
	return setConfig{
		T:     &now,
		SetNX: false,
	}
}

type RollingWindowSetElementOption interface {
	apply(setConfig) setConfig
}

type rollingWindowSetElementOptionFunc func(setConfig) setConfig

func (fn rollingWindowSetElementOptionFunc) apply(cfg setConfig) setConfig {
	return fn(cfg)
}

func WithTime(t time.Time) RollingWindowSetElementOption {
	return rollingWindowSetElementOptionFunc(func(gc setConfig) setConfig {
		gc.T = &t
		return gc
	})
}

func WithSetNX(b bool) RollingWindowSetElementOption {
	return rollingWindowSetElementOptionFunc(func(gc setConfig) setConfig {
		gc.SetNX = b
		return gc
	})
}

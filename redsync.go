package redsync

import (
	"github.com/applinskinner/redsync/redis"
	"time"
)

// Redsync provides a simple method for creating distributed mutexes using multiple Redis connection pools.
type Redsync struct {
	pools []redis.Pool
}

// New creates and returns a new Redsync instance from given Redis connection pools.
func New(pools []redis.Pool) *Redsync {
	return &Redsync{
		pools: pools,
	}
}

// NewMutex returns a new distributed mutex with given name.
func (r *Redsync) NewMutex(name string, options ...Option) *Mutex {
	m := &Mutex{
		name:         name,
		expiry:       8 * time.Second,
		tries:        32,
		delayFunc:    func(tries int) time.Duration { return 500 * time.Millisecond },
		genValueFunc: genValue,
		factor:       0.01,
		quorum:       len(r.pools)/2 + 1,
		pools:        r.pools,
	}
	for _, o := range options {
		o.Apply(m)
	}
	return m
}

// An Option configures a mutex.
type Option interface {
	Apply(*Mutex)
}

// OptionFunc is a function that configures a mutex.
type OptionFunc func(*Mutex)

// Apply calls f(mutex)
func (f OptionFunc) Apply(mutex *Mutex) {
	f(mutex)
}

// SetExpiry can be used to set the expiry of a mutex to the given value.
func SetExpiry(expiry time.Duration) Option {
	return OptionFunc(func(m *Mutex) {
		m.expiry = expiry
	})
}

// SetTries can be used to set the number of times lock acquire is attempted.
func SetTries(tries int) Option {
	return OptionFunc(func(m *Mutex) {
		m.tries = tries
	})
}

// SetRetryDelay can be used to set the amount of time to wait between retries.
func SetRetryDelay(delay time.Duration) Option {
	return OptionFunc(func(m *Mutex) {
		m.delayFunc = func(tries int) time.Duration {
			return delay
		}
	})
}

// SetRetryDelayFunc can be used to override default delay behavior.
func SetRetryDelayFunc(delayFunc DelayFunc) Option {
	return OptionFunc(func(m *Mutex) {
		m.delayFunc = delayFunc
	})
}

// SetDriftFactor can be used to set the clock drift factor.
func SetDriftFactor(factor float64) Option {
	return OptionFunc(func(m *Mutex) {
		m.factor = factor
	})
}

// SetGenValueFunc can be used to set the custom value generator.
func SetGenValueFunc(genValueFunc func() (string, error)) Option {
	return OptionFunc(func(m *Mutex) {
		m.genValueFunc = genValueFunc
	})
}

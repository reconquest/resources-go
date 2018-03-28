package resources

import (
	"sync"
	"time"
)

type loaderResult struct {
	result interface{}
	err    error
}

// ErrorHandler is a function that handles error that happened while updating
// a resource. error can't be nil.
type ErrorHandler func(error)

// Loader is a function that loads/updates resource, result and error returning
// values will be stored as the cache. If error is returned by loader, it will
// be also returned as error of Get() call.
type Loader func() (interface{}, error)

// Resources is a common structure that incapsulates work with cache and
// provides useful methods for synchronizing resources.
type Resources struct {
	loaders      *sync.Map
	items        *sync.Map
	interval     time.Duration
	errorHandler ErrorHandler
}

// NewResources creates a new instance for working with cachable resources.
func NewResources(interval time.Duration) *Resources {
	return &Resources{
		items:    &sync.Map{},
		loaders:  &sync.Map{},
		interval: interval,
	}
}

// SetLoader sets given Loader function for given key. By given key user will
// be able retrieve cache value.
func (resources *Resources) SetLoader(key string, loader Loader) {
	resources.loaders.Store(key, loader)
}

// SetErrorHandler is optional method, passed function will be called and error
// will be passed to it if some error happened while updating refresh caches.
// This function is not called if you specified waitOnce=true in Sync() method.
func (resources *Resources) SetErrorHandler(handler ErrorHandler) {
	resources.errorHandler = handler
}

// Sync given resources, if waitOnce is true then Sync() will block program
// execution until all resources synchronized or an error is occurred, after
// that goroutine with synchronization loop will be created. If waitOnce is
// false then error will be never returned (even if any). If you use
// waitOnce=false then you probably want to run it as goroutine.
func (resources *Resources) Sync(waitOnce bool) error {
	if waitOnce {
		err := resources.sync(true)

		go func() {
			time.Sleep(resources.interval)
			resources.syncLoop()
		}()

		return err
	} else {
		resources.syncLoop()

		return nil
	}
}

func (resources *Resources) syncLoop() {
	for {
		resources.sync(false)
		time.Sleep(resources.interval)
	}
}

func (resources *Resources) sync(must bool) error {
	var syncErr error
	resources.loaders.Range(func(key interface{}, value interface{}) bool {
		name, loader := key.(string), value.(Loader)

		result, err := loader()
		if err != nil {
			if must {
				syncErr = err
				return false
			}

			if resources.errorHandler != nil {
				resources.errorHandler(err)
			}
		}

		resources.items.Store(name, loaderResult{
			result: result,
			err:    err,
		})

		return true
	})

	return syncErr
}

// Get value by given key, if a resource with specified key is not defined,
// the function will return ErrorUndefinedResource error. The function returns
// error that happened while loading/updating a resource.
func (resources *Resources) Get(key string) (interface{}, error) {
	raw, ok := resources.items.Load(key)
	if !ok {
		return nil, ErrorUndefinedResource{key}
	}

	result := raw.(loaderResult)

	return result.result, result.err
}

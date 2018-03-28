package resources

import (
	"sync"
	"time"
)

type loaderResult struct {
	result interface{}
	err    error
}

type ErrorHandler func(error)

type Loader func() (interface{}, error)

type Resources struct {
	loaders      *sync.Map
	items        *sync.Map
	interval     time.Duration
	errorHandler ErrorHandler
}

func NewResources(interval time.Duration) *Resources {
	return &Resources{
		items:    &sync.Map{},
		loaders:  &sync.Map{},
		interval: interval,
	}
}

func (resources *Resources) SetLoader(key string, loader Loader) {
	resources.loaders.Store(key, loader)
}

func (resources *Resources) SetErrorHandler(handler ErrorHandler) {
	resources.errorHandler = handler
}

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

func (resources *Resources) Get(key string) (interface{}, error) {
	raw, ok := resources.items.Load(key)
	if !ok {
		return nil, ErrorUndefinedResource{key}
	}

	result := raw.(loaderResult)

	return result.result, result.err
}

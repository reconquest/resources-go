package resources

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResourcesSync_WaitOnce_ReturnsErrorIfAny(t *testing.T) {
	test := assert.New(t)

	resources := NewResources(time.Second * 2)
	resources.SetLoader("test-a", func() (interface{}, error) {
		return "test-a", nil
	})

	resources.SetLoader("test-b", func() (interface{}, error) {
		return "test-b-value", errors.New("test-b-error")
	})

	err := resources.Sync(true)
	test.Error(err)
	test.EqualError(err, "test-b-error")
}

func TestResourcesSync_WaitOnce_ReturnsNilIfNoError(t *testing.T) {
	test := assert.New(t)

	resources := NewResources(time.Second * 2)
	resources.SetLoader("test-a", func() (interface{}, error) {
		return "test-a", nil
	})

	resources.SetLoader("test-b", func() (interface{}, error) {
		return "test-b-value", nil
	})

	err := resources.Sync(true)
	test.NoError(err)
}

func TestResourcesGet_SyncWithWait_ReturnsValueAndNoError(t *testing.T) {
	test := assert.New(t)

	resources := NewResources(time.Second * 2)
	resources.SetLoader("test-a", func() (interface{}, error) {
		return "test-a-value", nil
	})
	resources.SetLoader("test-b", func() (interface{}, error) {
		return "test-b-value", nil
	})

	err := resources.Sync(true)
	test.NoError(err)

	value, err := resources.Get("test-a")
	test.NoError(err)
	test.Equal("test-a-value", value)
}

func TestResourcesGet_SyncWithoutWait_ReturnsValueAndNoError(t *testing.T) {
	test := assert.New(t)

	resources := NewResources(time.Second * 2)
	resources.SetLoader("test-a", func() (interface{}, error) {
		return "test-a-value", nil
	})
	resources.SetLoader("test-b", func() (interface{}, error) {
		return "test-b-value", nil
	})

	go resources.Sync(false)

	// no other way if we specified waitonce = false
	time.Sleep(time.Millisecond * 10)

	value, err := resources.Get("test-a")
	test.NoError(err)
	test.Equal("test-a-value", value)
}

func TestResourcesGet_SyncWithoutWait_ReturnsUndefinedResourceError(t *testing.T) {
	test := assert.New(t)

	resources := NewResources(time.Second * 2)
	resources.SetLoader("test-a", func() (interface{}, error) {
		time.Sleep(time.Second)
		return "test-a-value", nil
	})

	go resources.Sync(false)

	value, err := resources.Get("test-a")
	test.Error(err)
	test.EqualError(err, "undefined resource: test-a")
	test.Nil(value)
}

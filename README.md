# resources-go

[![](https://godoc.org/github.com/reconquest/resources-go?status.svg)](http://godoc.org/github.com/reconquest/resources-go)

This library helps to build cache-like storage that always has a value. The
basic problem of existing cache packages is that cache doesn't exist until some
code calls getter, but it's not suitable approach if human will need to wait
for this getter. This package offers another way - it refreshes caches every N
seconds/minutes and a user will get last value from cache (LRU).

Let's imagine that you build a web service that works with data from slow
backend service. You can use memcached-like solutions and set expiration time
for cached entries. But then caches will expire, and a http client will have
increased response time because your code will request value from some slow
backend again and cache it again. But there is a way to avoid increasing
response time - retrieve/compute data and cache it in background goroutine with
given interval of time, in this way cache storage will always have a value.

# how to build and use it

```go
func main() {
    resources := NewResources(time.Second*10)

    resources.SetLoader("kubernetes-pods", ListKubernetesPods)

    err := resources.Sync(true)
    if err != nil {
        log.Fatal(err)
    }

    value, err := resources.Get("kubernetes-pods")
    if err != nil {
        log.Println(err)
    } else {
        fmt.Println(err)
    }
}

func ListKubernetesPods() (interface{}, error) {
    // do something slow such as kubernetes API
    return something, nil
}
```

In this example we specify refresh interval as 10 seconds, but waitOnce
parameter of Sync() is set to true, it means that Sync will block until
synchronization is complete of failed, after completing synchronization it will
run goroutine with synchronization with given time interval. If synchronization
is failed it will return error, but goroutine will be started anyway too.

# when you should not use this package

If you are going to create cache entity dynamically and you don't control cache
identifier, then it would be better to use memcached-like solutions.

This library is not cache miss tolerant.

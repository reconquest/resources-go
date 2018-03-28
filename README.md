# resources-go

This library helps to build cache-like storage that always has a value. The
basic problem of existing cache packages is that cache doesn't exist until some
code calls getter, but it's not suitable approach if human will need to wait
for this getter. This package offers another way - it refreshes caches every N
seconds/minutes and a user will get last value from cache (LRU).

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

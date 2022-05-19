# Lock

1. Redis Cache
2. Redis Lock

### Lock Use

```go

lock := cache.Manage.AcquireLocker("example:100")
defer lock.Release()

ctx := log.NewContext()
log.Infofc(ctx, "begin locker")

got, err := lock.Apply(ctx)
if err != nil {
    log.Errorfc(ctx, "apply error: %v.", err)
    return
}

if !got {
    log.Infofc(ctx, "apply failed")
    return
}

// continue

```



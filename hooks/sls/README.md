sls
===


### example

```
hooks {
       expander {}

       sls {
            endpoint = "cn-beijing.log.aliyuncs.com"
            access-key-id = ""
            access-key-secret = ""

            project = gogap   # sls project name
            store = test      # default log store name

            levels = ["debug","info","error","warn"] # these level's logs will post to sls

            fields {
            	store = "store"    # you could use WithField("store","test2") to specfic the store, orelse use default
                topic = "err_ns"   # you could use WithField("err_ns","this-is-topic") to specfic the topic
                source = "source"  # you could use WithField("source","this-is-soure") to specfic the source
                tags = "tags"      # you could use WithField("tags", map[string]string{"k","v"}) to specfic the tags
                context ="err_ctx" # you could use WithField("err_ctx", string|interface{}) to specfic the context
            }
       }

}
```


```go
mikeLoger.WithError(e).WithField("source", "IAMSOURCE").WithField("store", "test2").Errorln("hello with gogap errors")
```
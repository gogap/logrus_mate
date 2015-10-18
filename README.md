# Logrus Mate <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

**Logrus mate** is a tool for [Logrus](https://github.com/Sirupsen/logrus), it will help you to initial logger by config, including `Formatter`, `Hook`，`Level` and `Environments`.

#### Example

**Example 1:**

Using internal default logrus mate:

```go
package main

import (
    "github.com/gogap/logrus_mate"
)

func main() {
    logrus_mate.Logger().Infoln("Using internal defualt logurs mate")
}
```

**Example 2:**

Create new logger:

```go
package main

import (
    "github.com/gogap/logrus_mate"
)

func main() {
    loggerConf := logrus_mate.LoggerConfig{
        Level: "info",
        Formatter: logrus_mate.FormatterConfig{
            Name: "json",
        },
    }
    
    // package level
    if jackLogger, err := logrus_mate.NewLogger("jack", loggerConf); err != nil {
        return
    } else {
        jackLogger.Infoln("hello logurs")
    }

    // so, the jack added into internal logurs mate
    logrus_mate.Logger("jack").Debugln("not print")
    logrus_mate.Logger("jack").Infoln("Hello, I am A Logger from jack")
}
```

**Example 3:**

Create logurs mate from config file (also you could fill config struct manually):

```bash
export RUN_MODE=production
```

`mate.conf`
```json
{
    "env_keys": {
        "run_env": "RUN_MODE",
        "env_json": ""
    },
    "loggers": [{
        "name": "mike",
        "config": {
            "production": {
                "level": "error",
                "formatter": {
                    "name": "json"
                },
                "hooks":{
                    "syslog":{
                        "network":"udp",
                        "address":"localhost:514",
                        "priority":"LOG_ERR",
                        "tag":""
                    }
                }
            }
        }
    }]
}
```

```go
package main

import (    
    "github.com/gogap/logrus_mate"
    _ "github.com/gogap/logrus_mate/hooks/syslog"
)

func main() {
    if mateConf, err := logrus_mate.LoadLogrusMateConfig("mate.conf"); err != nil {
        return
    } else {
        if newMate, err := logrus_mate.NewLogrusMate(mateConf); err != nil {
            return
        } else {
            newMate.Logger("mike").Errorln("I am mike in new logrus mate")
        }
    }
}
```

> In this example, we used the syslog hook, so we should import package of syslog

``` go
import _ "github.com/gogap/logrus_mate/hooks/syslog"
```

logrus mate support environments notion, same logger could have different environment config, the above `mate.conf` only have `production` config, so if `RUN_MODE` is `production`, it will use this section's options, or else, there have no loggers generate.

#### Environments

different environment could have own `level`, `hooks` and `formatters`, logrus mate have `Environments` config for create the instance, you can see the above config file of `mate.conf`

```go
type Environments struct {
    RunEnv  string `json:"run_env"`
    EnvJson string `json:"env_json"`
}
```

`run_env`: this filed is the key of run env, it will get actual value from environment by this key.
`env_json`: the env_json's env key, if this field is not empty, it will use `env_json` to compile config file while you use func `logrus_mate.LoadLogrusMateConfig`, please forward to the project of [env_json](https://github.com/gogap/env_json) to known more details.

#### Hooks
| Hook  | Options |
| ----- | ----------- |
| [Airbrake](https://github.com/gemnasium/logrus-airbrake-hook) | `project_id` `api_key` `env`|
| [Syslog](https://github.com/Sirupsen/logrus/blob/master/hooks/syslog/syslog.go) | `network` `address` `priority` `tag`|
| [BugSnag](https://github.com/Sirupsen/logrus/blob/master/hooks/bugsnag/bugsnag.go) | `api_key` |
| [Slackrus](https://github.com/johntdyer/slackrus) | `url` `levels` `channel` `emoji` `username`|
| [Graylog](https://github.com/gemnasium/logrus-graylog-hook) | `address` `facility` `extra`|
| [Mail](https://github.com/zbindenren/logrus_mail) | `app_name` `host` `port` `from` `to` `username` `password`|

When we need use above hooks, we need import these package as follow:

```go
import _ "github.com/gogap/logrus_mate/hooks/syslog"
import _ "github.com/gogap/logrus_mate/hooks/mail"
```

If you want write your own hook, you just need todo as follow:

```go
package myhook

import (
    "github.com/gogap/logrus_mate"
)

type MyHookConfig struct {
    Address  string `json:"address"`
}

func init() {
    logrus_mate.RegisterHook("myhook", NewMyHook)
}

func NewMyHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
    conf := MyHookConfig{}
    if err = options.ToObject(&conf); err != nil {
        return
    }

    // write your hook logic code here

    return
}
```


#### Formatters

**internal formatters:**

| Formatter  | Output Example |
| ----- | ----------- |
|null||
|text|DEBU[0000] Hello Default Logrus Mate|
|json|{"level":"info","msg":"Hello, I am A Logger from jack","time":"2015-10-18T21:24:19+08:00"}|

**3rd formatters:**

| Formatter  | Output Example |
| ----- | ----------- |
|logstash||

When we need use 3rd formatter, we need import these package as follow:

```go
import _ "github.com/gogap/logrus_mate/formatters/logstash"
```

If you want write your own formatter, you just need todo as follow:

```go
package myformatter

import (
    "github.com/gogap/logrus_mate"
)

type MyFormatterConfig struct {
    Address  string `json:"address"`
}

func init() {
    logrus_mate.RegisterFormatter("myformatter", NewMyFormatter)
}

func NewMyFormatter(options logrus_mate.Options) (formatter logrus.Formatter, err error) {
    conf := MyFormatterConfig{}
    if err = options.ToObject(&conf); err != nil {
        return
    }

    // write your formatter logic code here

    return
}
```
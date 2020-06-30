package main

import (
  "flag"
  "fmt"
  "github.com/kataras/iris/v12"
  subscriber "github.com/mingcheng/ssr-subscriber"
  "io/ioutil"
  "log"
  "os"
  "time"

  "gopkg.in/yaml.v2"
)

var (
  configFile string
  configure  = Configure{}

  httpMode bool
  err      error

  CheckedConfig []*subscriber.Config
  LastCheckTime time.Time
)

func init() {
  flag.StringVar(&configFile, "config", "/etc/ssr-subscriber.yml", "subscribe configure file")
  flag.BoolVar(&httpMode, "http", false, "start http listen mode")

  flag.Parse()
}

func main() {
  // read configure from config file
  if yamlFile, err := ioutil.ReadFile(configFile); err != nil {
    log.Fatalln(err)
  } else {
    if err := yaml.Unmarshal(yamlFile, &configure); err != nil {
      log.Fatalln(err)
    }
  }

  fmt.Printf("%v", configure)

  // check output directory
  if stat, err := os.Stat(configure.Output); err != nil || !stat.IsDir() {
    log.Fatal(err)
  }

  CheckedConfig, err = fetchAndCheck()
  if err != nil {
    log.Fatal(err)
  }

  if httpMode {
    ticker := time.NewTicker(time.Duration(configure.Interval) * time.Hour)
    go func() {
      for {
        select {
        case <-ticker.C:
          if CheckedConfig, err = fetchAndCheck(); err != nil {
            _, _ = fmt.Fprint(os.Stderr, err.Error())
          }
        }
      }
    }()

    // start web server
    app := iris.New()

    app.Handle("GET", "/", func(ctx iris.Context) {
      ctx.Header("Last-Check", LastCheckTime.String())
      if len(CheckedConfig) > 0 {
        _, _ = ctx.JSON(CheckedConfig)
      } else {
        ctx.NotFound()
      }
    })

    app.Handle("GET", "/config", func(ctx iris.Context) {
      _, _ = ctx.JSON(configure)
    })

    app.Handle("GET", "/last-check", func(ctx iris.Context) {
      _, _ = ctx.WriteString(LastCheckTime.String())
    })

    err = app.Run(iris.Addr(configure.Bind))
  }
}

// fetchAndCheck that fetch configs from subscriber url or file, then check its health
func fetchAndCheck() ([]*subscriber.Config, error) {
  // TODO sync.RWMutex{}
  var (
    configs []*subscriber.Config
    err     error
  )

  // do not bind listen address one-shot only
  configs, err = fetchNodes(append(configure.URL, configure.File...))
  if err != nil {
    return nil, err
  }

  configs, err = checkAndSaveConfigs(configs, configure.Check, configure.Output)
  if err != nil {
    return nil, err
  }

  if configure.AutoClean {
    defer cleanExceedConfig(configure.Output, time.Duration(configure.Exceed)*(time.Hour*24))
  }

  LastCheckTime = time.Now()
  return configs, err
}

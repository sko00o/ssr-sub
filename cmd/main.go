package main

import (
	"flag"
	"math/rand"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/judwhite/go-svc"
	iris "github.com/kataras/iris/v12"
	subscriber "github.com/mingcheng/ssr-subscriber"
	log "github.com/sirupsen/logrus"
)

var (
	BuildTime = "unknown"
	Version   = "unknown"
)

var (
	configFile string
	IrisApp    *iris.Application
)

type program struct {
	Config  *Configure
	Fetcher *subscriber.Fetcher
}

func (p *program) Init(_ svc.Environment) error {
	// start fetch and check
	p.Fetcher = &subscriber.Fetcher{
		Output:      p.Config.Output,
		Exceed:      time.Duration(p.Config.Exceed) * time.Hour * 24,
		AutoClean:   p.Config.AutoClean,
		Sources:     append(p.Config.File, p.Config.URL...),
		CheckConfig: p.Config.Check,
		Proxy:       p.Config.Proxy,
		Interval:    time.Duration(p.Config.Interval) * time.Hour,
	}
	log.Trace(p.Fetcher)

	return nil
}

func (p *program) Start() error {
	if err := p.Fetcher.Start(); err != nil {
		log.Error(err)
		return err
	}

	IrisApp = iris.New()

	IrisApp.Get("/random/{type:string}", func(ctx iris.Context) {
		configs := p.Fetcher.AllConfigs()
		log.Debugf("get all %d configs", len(configs))

		if len(configs) == 0 {
			ctx.StatusCode(http.StatusNotFound)
			return
		}

		// initialize global pseudo random generator
		rand.Seed(time.Now().Unix())
		config := configs[rand.Intn(len(configs))]
		log.Trace(config)

		ctx.StatusCode(http.StatusOK)
		switch ctx.Params().Get("type") {
		case "yaml", "yml":
			_, _ = ctx.YAML(config)
		default:
			_, _ = ctx.JSON(config)
		}
	})

	IrisApp.Get("/all", func(ctx iris.Context) {
		configs := p.Fetcher.AllConfigs()
		log.Debugf("get all %d configs", len(configs))

		ctx.StatusCode(http.StatusOK)
		_, _ = ctx.JSON(configs)
	})

	IrisApp.Handle("GET", "/last-check-time", func(ctx iris.Context) {
		_, _ = ctx.WriteString(p.Fetcher.LastCheckTime().String())
	})

	go IrisApp.Run(iris.Addr(p.Config.Bind))
	return nil
}

func (p *program) Stop() error {
	if err := p.Fetcher.Stop(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func init() {
	flag.StringVar(&configFile, "config", "/etc/ssr-subscriber.yml", "subscribe configure file")
	flag.Usage = func() {
		log.Printf("version v%s(%s)", Version, BuildTime)
		os.Exit(0)
	}
	flag.Parse()
}

func main() {
	// parse configure from local file
	config, err := ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := svc.Run(&program{
		Config: config,
	}, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

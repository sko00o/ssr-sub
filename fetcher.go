package subscriber

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type Fetcher struct {
	Output      string
	Proxy       string
	Exceed      time.Duration
	Sources     []string
	CheckConfig CheckConfig

	Interval  time.Duration
	AutoClean bool

	configs       map[string]*Config
	ticker        *time.Ticker
	lastFetchTime time.Time
}

// Start to restore configs then fetch and update save
func (f *Fetcher) Start() error {
	f.configs = make(map[string]*Config)

	if err := f.Restore(); err != nil {
		log.Error(err)
	}

	if f.Interval > 0 {
		f.ticker = time.NewTicker(f.Interval)
		log.Debug(f.Interval)

		go func() {
			for {
				if err := f.Fetch(); err != nil {
					log.Error(err)
				}

				if err := f.Save(); err != nil {
					log.Error(err)
				}

				<-f.ticker.C
			}
		}()
	}

	return nil
}

// Stop and save configs
func (f *Fetcher) Stop() error {
	if f.ticker != nil {
		f.ticker.Stop()
	}

	_ = f.Save()
	return nil
}

// Fetch to get ssr configure from file or via url
func (f *Fetcher) Fetch() error {
	for _, uri := range f.Sources {
		var (
			configs []*Config
			err     error
		)

		if path.IsAbs(uri) {
			configs, err = FromFile(uri)
		} else {
			configs, err = FromURL(uri, f.Proxy)
		}

		if err != nil {
			log.Error(err)
			return err
		}

		// can not get anything
		if len(configs) == 0 {
			return fmt.Errorf("can not get any configure nodes")
		}

		for _, v := range configs {
			if CheckNode(v, f.CheckConfig) {
				log.Infof("add config %s into configs", v.Server)
				f.configs[f.genKey(v)] = v
			}
		}
	}

	f.lastFetchTime = time.Now()
	return nil
}

// Save configs into filesystem
func (f *Fetcher) Save() error {
	for _, v := range f.configs {
		fileName := path.Join(f.Output, fmt.Sprintf("%s.json", f.genKey(v)))

		data, err := json.Marshal(v)
		if err != nil {
			log.Error(err)
			return err
		}

		err = ioutil.WriteFile(fileName, data, 0644)
		if err != nil {
			log.Error(err)
			return err
		}

		log.Infof("saved %s into file %s", v.Server, fileName)
	}

	return nil
}

// Restore for restore config from local filesystem
func (f *Fetcher) Restore() error {
	// clean first
	_ = f.Clean()

	return filepath.Walk(f.Output, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			log.Warnf("%s is directory, so ignore", path)
			return nil
		}

		if filepath.Ext(path) == ".json" {
			config, err := f.fromSingleFile(path)
			if err != nil {
				log.Error(err)
				return err
			}

			log.Infof("restored config %s:%d", config.Server, config.ServerPort)
			f.configs[f.genKey(config)] = config
		}

		return nil
	})
}

// Clean to check and delete unhealthy config
func (f *Fetcher) Clean() error {
	// remove config if not passed health check
	for k, v := range f.configs {
		if !CheckNode(v, f.CheckConfig) {
			delete(f.configs, k)
		}
	}

	// remove outdated local fs configs
	return filepath.Walk(f.Output, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return fmt.Errorf("%s is not directory", info.Name())
		}

		// if file is out date, delete it
		timeleft := time.Now().Sub(info.ModTime())
		if filepath.Ext(path) == ".json" && timeleft < f.Exceed {
			if f.AutoClean {
				log.Warnf("auto clean file %s", path)
				return syscall.Unlink(path)
			}
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			_, err := f.fromSingleFile(path)

			if err != nil {
				if f.AutoClean {
					log.Warnf("auto clean file %s", path)
					return syscall.Unlink(path)
				}

				return err
			}
		}

		return nil
	})
}

// fromSingleFile to read and parse single config from file
func (f *Fetcher) fromSingleFile(path string) (*Config, error) {
	var config Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if !CheckNode(&config, f.CheckConfig) {
		err := fmt.Errorf("healty check failed, %s:%d is not avaiable", config.Server, config.ServerPort)
		log.Warn(err)
		return nil, err
	}

	return &config, nil
}

// genKey for generate a string key for single config
func (f *Fetcher) genKey(c *Config) string {
	return fmt.Sprintf("_%s:%d", c.Server, c.ServerPort)
}

// AllConfigs to return all stored configs
func (f *Fetcher) AllConfigs() []Config {
	var configs []Config

	for _, v := range f.configs {
		configs = append(configs, *v)
	}

	return configs
}

// LastCheckTime to return last check time
func (f *Fetcher) LastCheckTime() time.Time {
	return f.lastFetchTime
}

package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	subscriber "github.com/mingcheng/ssr-subscriber.go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"time"
)

var configFile string

type Configure struct {
	Url    []string `yaml:"url"`
	File   []string `yaml:"file"`
	Output string   `yaml:"output"`
	Check  check    `yaml:"check"`
}

type check struct {
	Timeout string `yaml:"timeout"`
	Not     string `yaml:"not"`
}

func checkNode(node *subscriber.Config, check check) bool {
	if matched, _ := regexp.MatchString(check.Not, node.Remarks); matched {
		log.Printf("remarks %s not allowed, ignore", node.Remarks)
		return false
	}

	if duration, err := time.ParseDuration(check.Timeout); err != nil {
		log.Fatalln(err)
		return true
	} else {
		_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", node.Server, node.ServerPort), duration);
		if err != nil {
			log.Printf("check %s:%d, failed", node.Server, node.ServerPort)
			return false
		} else {
			log.Printf("check %s:%d, passed", node.Server, node.ServerPort)
		}
	}

	return true
}

func saveConfigToFile(config *subscriber.Config, dir string) error {
	address := fmt.Sprintf("%s:%d", config.Server, config.ServerPort)
	name := fmt.Sprintf("%x", md5.Sum([]byte(address)))
	fileName := fmt.Sprintf("%s%c%s.json", dir, os.PathSeparator, name)
	bs, _ := json.MarshalIndent(config, "", " ")
	if err := ioutil.WriteFile(fileName, bs, 0644); err != nil {
		return err
	} else {
		log.Printf("saved to %s\n", fileName)
		return nil
	}
}

func init() {
	flag.StringVar(&configFile, "c", "", "subscribe configure file")
	flag.Parse()
}

func main() {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	configure := Configure{}
	if err := yaml.Unmarshal(yamlFile, &configure); err != nil {
		log.Fatalln(err)
	}

	if stat, err := os.Stat(configure.Output); err != nil || !stat.IsDir() {
		log.Fatal()
	}

	var configNodes []*subscriber.Config

	// patch parse urls
	if configure.Url != nil && len(configure.Url) > 0 {
		for _, u := range configure.Url {
			if nodes, err := subscriber.FromUrl(u); err != nil {
				log.Printf(err.Error())
			} else {
				for _, node := range nodes {
					configNodes = append(configNodes, node)
				}
			}
		}
	}

	// patch parse files
	if configure.File != nil && len(configure.File) > 0 {
		for _, c := range configure.File {
			if nodes, err := subscriber.FromFile(c); err != nil {
				log.Printf(err.Error())
			} else {
				for _, node := range nodes {
					configNodes = append(configNodes, node)
				}
			}
		}
	}

	if len(configNodes) <= 0 {
		log.Fatalln("can not get any configure nodes")
	}

	for _, node := range configNodes {
		if checkNode(node, configure.Check) {
			_ = saveConfigToFile(node, configure.Output)
		}
	}
}

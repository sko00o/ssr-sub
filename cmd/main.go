package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	subscriber "github.com/mingcheng/ssr-subscriber.go"
	"gopkg.in/yaml.v2"
)

var configFile string

// Configure struct, for more information see config-example.yaml file
type Configure struct {
	URL    []string `yaml:"url"`
	File   []string `yaml:"file"`
	Output string   `yaml:"output"`
	Check  check    `yaml:"check"`
}

type check struct {
	Timeout string `yaml:"timeout"`
	Not     string `yaml:"not"`
}

// checkNode for check ssr config server available
func checkNode(node *subscriber.Config, c check) bool {
	if matched, _ := regexp.MatchString(c.Not, node.Remarks); matched {
		log.Printf("remarks %s not allowed, ignore", node.Remarks)
		return false
	}

	if duration, err := time.ParseDuration(c.Timeout); err != nil {
		log.Fatalln(err)
	} else {
		_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", node.Server, node.ServerPort), duration)
		if err != nil {
			log.Printf("check %s:%d, failed", node.Server, node.ServerPort)
			return false
		}
	}

	log.Printf("check %s:%d, passed", node.Server, node.ServerPort)
	return true
}

// saveConfigToFile for save config to JSON format file
func saveConfigToFile(config *subscriber.Config, dir string) error {
	address := fmt.Sprintf("%s:%d", config.Server, config.ServerPort)
	name := fmt.Sprintf("%x", md5.Sum([]byte(address)))
	fileName := fmt.Sprintf("%s%c%s.json", dir, os.PathSeparator, name)
	bs, _ := json.MarshalIndent(config, "", " ")
	if err := ioutil.WriteFile(fileName, bs, 0644); err != nil {
		return err
	}

	log.Printf("saved to %s\n", fileName)
	return nil
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

	for _, o := range append(configure.URL, configure.File...) {
		var nodes []*subscriber.Config
		if _, err := os.Stat(o); os.IsExist(err) {
			nodes, _ = subscriber.FromFile(o)
		} else {
			nodes, _ = subscriber.FromURL(o)
		}

		configNodes = append(configNodes, nodes...)
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

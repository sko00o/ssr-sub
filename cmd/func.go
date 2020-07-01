package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	subscriber "github.com/mingcheng/ssr-subscriber"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// cleanExceedConfig to clean outdated configure file in specified directory
func cleanExceedConfig(dir string, duration time.Duration) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		var returnError error

		//&& time.Now().Sub(info.ModTime()) < duration
		leftTime := time.Now().Sub(info.ModTime())
		if !info.IsDir() && filepath.Ext(path) == ".json" && leftTime > duration {
			// remove(unlink) configure file if outdated
			returnError = syscall.Unlink(path)
		}

		return returnError
	})
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

// fetchNodes that fetch sources and parse node config once
func fetchNodes(sources []string) ([]*subscriber.Config, error) {
	var (
		configNodes []*subscriber.Config
		err         error
	)

	for _, o := range sources {
		var nodes []*subscriber.Config
		if _, err = os.Stat(o); os.IsExist(err) {
			nodes, err = subscriber.FromFile(o)
		} else {
			nodes, err = subscriber.FromURL(o, configure.Proxy)
		}

		configNodes = append(configNodes, nodes...)
	}

	if len(configNodes) <= 0 {
		err = fmt.Errorf("can not get any configure nodes")
	}

	return configNodes, err
}

// checkAndSaveConfigs to check and save the node config if health check is fine
func checkAndSaveConfigs(nodes []*subscriber.Config, config subscriber.CheckConfig, dir string) ([]*subscriber.Config, error) {
	var (
		err          error
		checkedNodes []*subscriber.Config
	)

	for _, node := range nodes {
		if subscriber.CheckNode(node, config) {
			// put to checked nodes when pass-though the health check
			checkedNodes = append(checkedNodes, node)
			err = saveConfigToFile(node, dir)
		}
	}

	return checkedNodes, err
}

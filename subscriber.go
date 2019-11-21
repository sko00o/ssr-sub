package subscriber

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Config set for ssr config
type Config struct {
	Server        string `json:"server"`
	ServerPort    int    `json:"server_port"`
	Method        string `json:"method"`
	Protocol      string `json:"protocol"`
	ProtocolParam string `json:"protocol_param"`
	OBFS          string `json:"obfs"`
	OBFSParam     string `json:"obfs_param"`
	Password      string `json:"password"`
	Remarks       string `json:"remarks"`
	Group         string `json:"group"`
}

// FromURL fetch and parse configs from url
func FromURL(url string) ([]*Config, error) {
	response, err := http.Get(url)

	if err != nil || response.StatusCode != http.StatusOK {
		return nil, errors.New("request subscribe url error")
	}

	return FromReader(response.Body)
}

// FromFile parse configs from local base64-hashed file
func FromFile(path string) ([]*Config, error) {
	stat, err := os.Stat(path)
	if err != nil || !stat.Mode().IsRegular() {
		return nil, errors.New("not a regular file")
	}

	fd, err := os.OpenFile(path, os.O_RDONLY, os.ModeTemporary)

	if err != nil {
		return nil, err
	}

	defer fd.Close()
	return FromReader(bufio.NewReader(fd))
}

// FromString parse from string
func FromString(data string) ([]*Config, error) {
	return FromReader(strings.NewReader(data))
}

// FromReader from steam reader
func FromReader(r io.Reader) ([]*Config, error) {
	reader := base64.NewDecoder(base64.RawStdEncoding, r)
	data, err := ioutil.ReadAll(reader)

	if err != nil || len(data) <= 0 {
		return nil, err
	}

	return Decode(string(data))
}

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

func FromUrl(url string) ([]*Config, error) {
	if response, err := http.Get(url); err != nil || response.StatusCode != http.StatusOK {
		return nil, errors.New("request subscribe url error")
	} else {
		return FromReader(response.Body)
	}
}

func FromFile(path string) ([]*Config, error) {
	stat, err := os.Stat(path)
	if err != nil || !stat.Mode().IsRegular() {
		return nil, errors.New("not a regular file")
	}

	if f, err := os.OpenFile(path, os.O_RDONLY, os.ModeTemporary); err != nil {
		return nil, err
	} else {
		defer f.Close()
		return FromReader(bufio.NewReader(f))
	}
}

func FromReader(r io.Reader) ([]*Config, error) {
	reader := base64.NewDecoder(base64.RawStdEncoding, r)
	if data, err := ioutil.ReadAll(reader); err != nil || len(data) <= 0 {
		return nil, err
	} else {
		return Decode(string(data))
	}
}

func FromString(data string) ([]*Config, error) {
	return FromReader(strings.NewReader(data))
}

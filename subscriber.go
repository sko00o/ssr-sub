package subscriber

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
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
func FromURL(url string, proxyString string) ([]*Config, error) {
	httpTransport := &http.Transport{}
	client := &http.Client{Transport: httpTransport, Timeout: 5 * time.Second}

	if len(proxyString) > 0 {
		log.Printf("Using socks5 proxy address %s", proxyString)
		dialer, err := proxy.SOCKS5("tcp", proxyString, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}

		httpTransport.DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			return dialer.Dial(network, addr)
		}
	}

	// Do requests
	response, err := client.Get(url)

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

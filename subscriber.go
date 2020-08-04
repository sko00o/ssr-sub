package subscriber

import (
  "bufio"
  "context"
  "encoding/base64"
  "errors"
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "net"
  "net/http"
  "os"
  "regexp"
  "time"

  "golang.org/x/net/proxy"
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

// CheckConfig for CheckNode usage
type CheckConfig struct {
  TCPTimeout string `yaml:"timeout"`
  Not        string `yaml:"not"`
}

// CheckNode for check ssr config server available
func CheckNode(node *Config, c CheckConfig) bool {
  if matched, _ := regexp.MatchString(c.Not, node.Remarks); matched {
    log.Printf("remarks %s not allowed, ignore", node.Remarks)
    return false
  }

  if duration, err := time.ParseDuration(c.TCPTimeout); err != nil {
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

  if data, err := ioutil.ReadAll(response.Body); err != nil {
    return nil, err
  } else {
    return FromBytes(data)
  }
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
  var (
    err     error
    decoded []byte
  )

  if decoded, err = base64.StdEncoding.DecodeString(data); err != nil {
    decoded, err = base64.RawStdEncoding.DecodeString(data)
  }

  if err != nil || len(decoded) <= 0 {
    return nil, err
  }

  return Decode(string(decoded))
}

// FromReader from steam reader
func FromReader(r io.Reader) ([]*Config, error) {
  data, err := ioutil.ReadAll(r)
  if err != nil {
    return nil, err
  }

  return FromBytes(data)
}

// FromBytes get configs from bytes array
func FromBytes(data []byte) ([]*Config, error) {
  return FromString(string(data))
}

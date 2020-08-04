package subscriber

import (
  "encoding/base64"
  "errors"
  "net/url"
  "strconv"
  "strings"
)

const ssrPrefix = "ssr://"

// Decode batched ssr links
func Decode(data string) ([]*Config, error) {
  var result []*Config
  uris := strings.Split(string(data), "\n")

  for _, uri := range uris {
    if uri == "" {
      continue
    }

    if node, err := DecodeURI(uri); err == nil {
      result = append(result, node)
    }
  }

  if len(result) > 0 {
    return result, nil
  } else {
    return nil, errors.New("result is empty")
  }
}

// decode for decoding url from base64 strings
func decode(enc string) ([]byte, error) {
  if result, errs := base64.RawURLEncoding.DecodeString(enc); errs != nil {
    return base64.StdEncoding.DecodeString(enc)
  } else {
    return result, nil
  }
}

// forceDecode for force decoding strings without any errors
func forceDecode(in string) string {
  if in != "" {
    if b, err := decode(in); err == nil {
      in = string(b)
    }
  }
  return in
}

// DecodeURI for decode URI params once
func DecodeURI(uri string) (*Config, error) {
  if !strings.HasPrefix(uri, ssrPrefix) {
    return nil, errors.New("not a valid ssr string")
  } else {
    uri = uri[len(ssrPrefix):]
  }

  b, err := decode(uri)
  if err != nil {
    return nil, err
  }

  s := string(b)
  c := &Config{}

  i := strings.Index(s, ":")
  if i > -1 {
    c.Server = strings.TrimSpace(s[:i])
    s = s[i+1:]
  }
  i = strings.Index(s, ":")
  if i > -1 {
    c.ServerPort, _ = strconv.Atoi(s[:i])
    s = s[i+1:]
  }
  i = strings.Index(s, ":")
  if i > -1 {
    c.Protocol = s[:i]
    s = s[i+1:]
  }
  i = strings.Index(s, ":")
  if i > -1 {
    c.Method = s[:i]
    s = s[i+1:]
  }
  i = strings.Index(s, ":")
  if i > -1 {
    c.OBFS = s[:i]
    s = s[i+1:]
  }
  i = strings.Index(s, "/")
  if i > -1 {
    c.Password = strings.TrimSpace(forceDecode(s[:i]))
    s = s[i+1:]
  }

  u, err := url.Parse(s)
  if err != nil {
    return nil, err
  }
  c.OBFSParam = forceDecode(u.Query().Get("obfsparam"))
  c.ProtocolParam = forceDecode(u.Query().Get("protoparam"))
  c.Remarks = forceDecode(u.Query().Get("remarks"))
  c.Group = forceDecode(u.Query().Get("group"))

  return c, nil
}

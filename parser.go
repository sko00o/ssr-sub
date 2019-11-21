package subscriber

import (
    "encoding/base64"
    "errors"
    "net/url"
    "strconv"
    "strings"
)

const ssrPrefix = "ssr://"

func base64Decode(in string) string {
    if in != "" {
        b, err := base64.RawURLEncoding.DecodeString(in)
        if err == nil {
            in = string(b)
        }
    }
    return in
}

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

// DecodeURI for decode URI params once
func DecodeURI(uri string) (*Config, error) {
    if !strings.HasPrefix(uri, ssrPrefix) {
        return nil, errors.New("not a valid ssr string")
    } else {
        uri = uri[len(ssrPrefix):]
    }

    b, err := base64.RawURLEncoding.DecodeString(uri)
    if err != nil {
        return nil, err
    }

    s := string(b)
    c := &Config{}

    i := strings.Index(s, ":")
    if i > -1 {
        c.Server = s[:i]
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
        c.Password = base64Decode(s[:i])
        s = s[i+1:]
    }

    u, err := url.Parse(s)
    if err != nil {
        return nil, err
    }
    c.OBFSParam = base64Decode(u.Query().Get("obfsparam"))
    c.ProtocolParam = base64Decode(u.Query().Get("protoparam"))
    c.Remarks = base64Decode(u.Query().Get("remarks"))
    c.Group = base64Decode(u.Query().Get("group"))

    return c, nil
}

package subscriber

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/mingcheng/ssr-subscriber/node"
)

const ssrPrefix = "ssr://"

// decode batched ssr links
func decode(data string) (map[string]node.Config, error) {
	result := make(map[string]node.Config)

	uris := strings.Split(data, "\n")

	for _, uri := range uris {
		if uri == "" {
			continue
		}

		if config, err := decodeURI(strings.TrimSpace(uri)); err == nil {
			result[config.ID] = *config
		}
	}

	if len(result) > 0 {
		return result, nil
	}

	return nil, errors.New("sorry, result is empty")
}

// base64decode for decoding url from base64 strings
func base64decode(enc string) ([]byte, error) {
	if result, errs := base64.RawURLEncoding.DecodeString(enc); errs != nil {
		return base64.StdEncoding.DecodeString(enc)
	} else {
		return result, nil
	}
}

// forceDecode for force decoding strings without any errors
func forceDecode(in string) string {
	if in != "" {
		if b, err := base64decode(in); err == nil {
			in = string(b)
		}
	}
	return in
}

// decodeURI for decode URI params once
func decodeURI(uri string) (*node.Config, error) {
	if !strings.HasPrefix(uri, ssrPrefix) {
		return nil, errors.New("not a valid ssr string")
	} else {
		uri = uri[len(ssrPrefix):]
	}

	b, err := base64decode(uri)
	if err != nil {
		return nil, err
	}

	s := string(b)
	c := &node.Config{}

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

	c.ID = genKey(fmt.Sprintf("%s:%d", c.Server, c.ServerPort))

	return c, nil
}

func genKey(suffix string) string {
	return base64.StdEncoding.EncodeToString([]byte(suffix))
}

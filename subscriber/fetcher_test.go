package subscriber

import (
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/mingcheng/ssr-subscriber/node"
	"github.com/stretchr/testify/assert"
)

var (
	RedisClient *redis.Client
	fetcher     *Fetcher
	err         error
)

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_SERVER"),
	})
}

func TestInit(t *testing.T) {
	fetcher = &Fetcher{
		RedisClient: RedisClient,
		Proxy:       os.Getenv("FETCH_PROXY"),
		CheckConfig: node.CheckConfig{
			TCPTimeout: "3s",
			Not:        "免费|普通|回国|过期|剩余",
		},
	}

	err := fetcher.Init()
	assert.Nil(t, err)
}

// func TestFromBytes(t *testing.T) {
//   if fetcher == nil {
//     TestInit(t)
//   }
//
//
//   const filePath = "./test/1.txt"
//
//   if _, err := os.Stat(filePath); os.IsExist(err) {
//     test, err := ioutil.ReadFile(filePath)
//     assert.Nil(t, err)
//
//     got, err := fetcher.FromBytes(test)
//     assert.Nil(t, err)
//     assert.NotEmpty(t, got)
//   }
// }

//
// func TestFromFile(t *testing.T) {
//   const filePath = "./test/1.txt"
//
//   if _, err := os.Stat(filePath); os.IsExist(err) {
//     got, err := FromFile(filePath)
//     assert.Nil(t, err)
//     assert.NotEmpty(t, got)
//   }
// }

func TestFromURL(t *testing.T) {
	if fetcher == nil {
		TestInit(t)
	}

	err = fetcher.FromURL("https://www.taobao.com/")
	assert.NotNil(t, err)

	err = fetcher.FromURL("https://724.subadd.xyz/link/KD4OLAyyHEzUOzVr?sub=1")
	assert.Nil(t, err)
}

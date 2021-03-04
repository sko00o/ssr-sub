package subscriber

import (
  "io/ioutil"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
  test, err := ioutil.ReadFile("./test/2.txt")
  assert.Nil(t, err)

  got, err := FromString(string(test))
  assert.Nil(t, err)
  assert.NotEmpty(t, got)
}

func TestFromBytes(t *testing.T) {
  test, err := ioutil.ReadFile("./test/2.txt")
  assert.Nil(t, err)

  got, err := FromBytes(test)
  assert.Nil(t, err)
  assert.NotEmpty(t, got)
}

func TestFromFile(t *testing.T) {
  got, err := FromFile("./test/1.txt")
  assert.Nil(t, err)
  assert.NotEmpty(t, got)
}

func TestFromURL(t *testing.T) {
  _, err := FromURL("https://www.taobao.com/", "")
  assert.NotNil(t, err)
}

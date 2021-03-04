package subscriber

import (
  "io/ioutil"
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
  const filePath = "./test/1.txt"

  if _, err := os.Stat(filePath); os.IsExist(err) {
    test, err := ioutil.ReadFile(filePath)
    assert.Nil(t, err)

    got, err := FromString(string(test))
    assert.Nil(t, err)
    assert.NotEmpty(t, got)
  }
}

func TestFromBytes(t *testing.T) {
  const filePath = "./test/1.txt"

  if _, err := os.Stat(filePath); os.IsExist(err) {
    test, err := ioutil.ReadFile(filePath)
    assert.Nil(t, err)

    got, err := FromBytes(test)
    assert.Nil(t, err)
    assert.NotEmpty(t, got)
  }
}

func TestFromFile(t *testing.T) {
  const filePath = "./test/1.txt"

  if _, err := os.Stat(filePath); os.IsExist(err) {
    got, err := FromFile(filePath)
    assert.Nil(t, err)
    assert.NotEmpty(t, got)
  }
}

func TestFromURL(t *testing.T) {
  _, err := FromURL("https://www.taobao.com/", "")
  assert.NotNil(t, err)
}

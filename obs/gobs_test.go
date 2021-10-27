package obs_test

import (
	"fmt"
	"testing"

	"github.com/daqiancode/gobs"
	"github.com/stretchr/testify/assert"
)

var accessKey = "#"
var secretKey = "#"
var endPoint = "#"
var bucket = "ad-dev"

func TestGet(t *testing.T) {
	cli, err := gobs.NewOBS(accessKey, secretKey, endPoint, bucket)
	assert.Nil(t, err)
	r, err := cli.ListFile("test", 0)
	assert.Nil(t, err)
	fmt.Printf("%#v\n", r)
}

func TestRead(t *testing.T) {
	cli, err := gobs.NewOBS(accessKey, secretKey, endPoint, bucket)
	assert.Nil(t, err)
	for i := 0; i < 10; i++ {
		r, err := cli.Read("test/drum/master.m3u8")
		assert.Nil(t, err)
		fmt.Println(r.String())
	}
}

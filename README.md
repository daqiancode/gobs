# gobs
Go obs client for huawei cloud OBS

### Example
```go
func TestRead(t *testing.T) {
	cli, err := gobs.NewOBS(accessKey, secretKey, endPoint, bucket)
	assert.Nil(t, err)
	r, err := cli.Read("texts/a.txt")
	assert.Nil(t, err)
	fmt.Println(r.String())
}
```
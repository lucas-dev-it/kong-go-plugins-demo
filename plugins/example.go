package main

import (
	"fmt"

	"github.com/Kong/go-pdk"
)

type Conf struct {
}

func New() interface{} {
	return &Conf{}
}

func (conf Conf) Access(kong *pdk.PDK) {
	path, err := kong.Request.GetPath()
	if err != nil {
		kong.Log.Err(err.Error())
	}

	err = kong.Response.SetHeader("x-go-example-path", fmt.Sprintf("path %s", path))
	if err != nil {
		kong.Log.Err(err.Error())
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	commonHttp "github.com/kulycloud/common/http"
	"github.com/kulycloud/common/logging"
)

var logger = logging.GetForComponent("service")

func main() {
	srv, err := commonHttp.NewServer(30000, echoHandler)
	if err != nil {
		logger.Panicw("could not create server", "error", err)
	}

	err = srv.Serve()

	if err != nil {
		logger.Panicw("could not serve", "error", err)
	}
}

type config struct {
	Data string `json:"data"`
	ContentType *string `json:"contentType,omitempty"`
}

func echoHandler(_ context.Context, request *commonHttp.Request) *commonHttp.Response {
	conf := &config{}
	err := json.Unmarshal([]byte(request.KulyData.Step.Config), conf)
	resp := commonHttp.NewResponse()

	if err != nil {
		resp.Status = 500
		resp.Body.Write([]byte(fmt.Sprintf("Invalid config specified: %s", err.Error())))
		return resp
	}

	if conf.ContentType != nil {
		resp.Headers.Set("Content-Type", *conf.ContentType)
	} else {
		resp.Headers.Set("Content-Type", "text/plain")
	}

	resp.Status = 200
	resp.Body.Write([]byte(conf.Data))

	return resp
}

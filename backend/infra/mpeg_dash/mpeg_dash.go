package mpegdash

import (
	"net/url"

	"github.com/walnuts1018/mucaron/backend/config"
)

type Controller struct {
	serverEndpoint url.URL
	adminToken     string
}

func NewController(cfg config.Config) (*Controller, error) {
	endpoint, err := url.Parse(cfg.MpegDashServerEndpoint)
	if err != nil {
		return nil, err
	}

	return &Controller{
		serverEndpoint: *endpoint,
		adminToken:     cfg.MpegDashAdminToken,
	}, nil
}

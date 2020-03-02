package es

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/kairen/rotate-elasticsearch-index/pkg/util"
	"github.com/olivere/elastic/v7"
)

const timeout = time.Second * 5

type Client struct {
	raw *elastic.Client
	ctx context.Context
}

func NewClient(cfg Config) (*Client, error) {
	if len(cfg.Servers) < 1 {
		return nil, errors.New("No servers specified")
	}

	options, err := cfg.getConfigOptions()
	if err != nil {
		return nil, err
	}

	rawClient, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	if cfg.Version == 0 {
		// Determine ElasticSearch Version
		pingResult, _, err := rawClient.Ping(cfg.Servers[0]).Do(context.Background())
		if err != nil {
			return nil, err
		}
		esVersion, err := strconv.Atoi(string(pingResult.Version.Number[0]))
		if err != nil {
			return nil, err
		}
		glog.Infof("Elasticsearch version detected: %d", esVersion)
		cfg.Version = uint(esVersion)
	}
	return &Client{
		raw: rawClient,
		ctx: context.Background(),
	}, nil
}

func (c *Client) RotateIndex(regx string, day int, dateFormat string) error {
	indices, err := c.catIndices(regx)
	if err != nil {
		return err
	}

	for _, i := range indices {
		date := util.ParseDate(i.Index)
		expired, err := util.IsExpired(day, date, dateFormat)
		if err != nil {
			return err
		}

		if expired {
			glog.V(3).Infof("Index %s has expired.", i.Index)
			if _, err := c.deleteIndex(i.Index); err != nil {
				return err
			}
			glog.Infof("Successfully removed index %s.", i.Index)
		}
	}
	return nil
}

func (c *Client) catIndices(index string) (elastic.CatIndicesResponse, error) {
	return c.raw.CatIndices().Index(index).Do(c.ctx)
}

func (c *Client) deleteIndex(index string) (*elastic.IndicesDeleteResponse, error) {
	return c.raw.DeleteIndex(index).Do(c.ctx)
}

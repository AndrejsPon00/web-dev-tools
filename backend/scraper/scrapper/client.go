package scrapper

import (
	"context"
	"sync"

	"github.com/AndrejsPon00/web-dev-tools/backend/module"
	"github.com/AndrejsPon00/web-dev-tools/backend/scrapper/banknote"
	"github.com/AndrejsPon00/web-dev-tools/backend/scrapper/pp.lv"
	"github.com/AndrejsPon00/web-dev-tools/backend/scrapper/ss.lv"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	PaginationChan chan *module.Pagination
	ResultChan     chan *module.PreviewPost
	ErrorChan      chan error
	WG             *sync.WaitGroup
	Params         *module.URLParams
	Done           context.CancelFunc
	Context        context.Context
	SearchedItem   string
}

func (c *Client) ScrapPosts() {
	for _, source := range c.Params.Sources {
		c.WG.Add(1)
		switch source {
		case module.SOURCE_SS:
			collector := colly.NewCollector()
			go ss.ScrapPosts(c.Params.SearchedItem, c.Params.SSCurrentPage, c.Params.Filter, c.WG, collector, c.PaginationChan, c.ResultChan, c.ErrorChan)
		case module.SOURCE_BANKNOTE:
			go banknote.ScrapPosts(c.Params.SearchedItem, c.Params.BanknoteCurrentPage, c.Params.Filter, c.WG, c.PaginationChan, c.ResultChan, c.ErrorChan)
		case module.SOURCE_FACEBOOK:
			//add scrap facebook
			c.WG.Done()
		case module.SOURCE_GELIOS:
			//add scrap gelios
			c.WG.Done()
		case module.SOURCE_PP:
			go pp.ScrapPosts(c.Params.SearchedItem, c.Params.PPCurrentPage, c.Params.Filter, c.WG, c.PaginationChan, c.ResultChan, c.ErrorChan)
		default:
			log.Errorln("Unknown source: ", source)
			c.WG.Done()
		}
	}

	c.WG.Done()
	c.WG.Wait()
	c.Done()
}

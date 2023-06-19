package crawler

import (
	"context"
	"sync"

	"github.com/cortze/ragno/crawler/db"
	models "github.com/cortze/ragno/pkg"
	"github.com/sirupsen/logrus"
)

type Crawler struct {
	ctx context.Context

	// host
	host *Host

	// database
	db *db.Database

	// discovery

	// peer connections

	// ip_locator

	// prometheus

}

func NewCrawler(ctx context.Context, conf CrawlerRunConf) (*Crawler, error) {
	// create a private key

	// create metrics module

	// create db crawler
	db, err := db.New(ctx, conf.DbEndpoint, 10, 2)
	if err != nil {
		logrus.Error("Couldn't init DB")
		return nil, err
	}

	// create a host
	host, err := NewHost(
		ctx,
		conf.HostIP,
		conf.HostPort,
		// default configuration so far
	)
	if err != nil {
		logrus.Error("failed to create host:")
		return nil, err
	}

	// set the file to read the enrs from if provided
	if conf.File != "" {
		ctx = context.WithValue(ctx, "File", conf.File)
	}

	// set the enr to connect to if provided
	if conf.Enr != "" {
		ctx = context.WithValue(ctx, "Enr", conf.Enr)
	}

	// set the number of workers if provided
	if conf.WorkerNum != 0 {
		ctx = context.WithValue(ctx, "Workers", conf.WorkerNum)
	}

	// create the discovery modules

	crwl := &Crawler{
		ctx:  ctx,
		host: host,
		db:   db,
	}

	// add all the metrics for each module to the prometheus endp

	return crwl, nil
}

func (c *Crawler) Run() error {
	// init list of peers to connect to
	peers, err := GetListELNodeInfo(c.ctx)
	if err != nil {
		logrus.Error("Couldn't get list of peers")
		return err
	}

	// channel for the saving of the peers
	savingChan := make(chan *models.ELNodeInfo, 100)
	// channel for the peers to connect to
	connChan := make(chan *models.ELNodeInfo, len(peers))

	// fill the channel with the peers
	go func() {
		for _, peer := range peers {
			connChan <- peer
		}
	}()

	// init the peer connections
	workersAmount := c.ctx.Value("Workers").(int)

	var wg sync.WaitGroup

	for i := 0; i < workersAmount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case peer := <-connChan:
					// try to connect to the peer
					Connect(&c.ctx, peer, c.host, savingChan)
					// save the peer
					c.db.InsertNode(peer)
				case <-c.ctx.Done():
					return
				}
			}
		}()
	}

	// init IP locator

	// init host

	// init discoveries

	return nil
}

func (c *Crawler) Close() {
	// finish discovery

	// stop host

	// stop IP locator

	// stop db

	logrus.Info("Ragno closing routine done! See you!")
}
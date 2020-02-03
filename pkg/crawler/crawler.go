package crawler

import (
  "net/url"
  "log"
  "sync"
)

type SiteData struct {
  Link *url.URL
  Body []byte
}

type urlCache struct {
  sync.Mutex
  data map[string]struct{}
}

func newUrlCache() *urlCache {
  return &urlCache{
    data: make(map[string]struct{}),
  }
}

func (c *urlCache) atomicSet(link *url.URL) bool {
  c.Lock()
  _, ok := c.data[link.String()]

  if !ok {
    c.data[link.String()] = struct{}{}
  }
  
  c.Unlock()
  return !ok
}

type Crawler struct {
	instances  []*url.URL
	numWorkers int
	cache      *urlCache
}

func NewCrawler(instances []*url.URL, numWorkers int) *Crawler {
	return &Crawler{
		instances:  instances,
		numWorkers: numWorkers,
		cache:      newUrlCache(),
	}
}

func (c *Crawler) crawl(link *url.URL, sig <-chan bool, jobs chan<- int, data chan<- SiteData, errs chan<- error) {
  <-sig
  body, links, err := fetch(link)

  if err != nil {
    errs <- err
    return
  }

  for _, link := range links {
    if c.cache.atomicSet(link) {
      jobs <- 1

      go c.crawl(link, sig, jobs, data, errs)
    }
  }

  data <- SiteData{
    link,
    body,
  }
}

func (c *Crawler) Start() {
	sig := make(chan bool, c.numWorkers)
  jobs := make(chan int)
  data := make(chan SiteData)
  errs := make(chan error)
  defer close(sig)

  for _, link := range c.instances {
	  go c.crawl(link, sig, jobs, data, errs)
  }

  for i := 0; i < c.numWorkers; i++ {
    sig <- true
  }

  toFetch := 1
  for toFetch > 0 {
    select {
    case d := <-data:
    	// TODO: Extract information
      log.Println(d.Link.String())

      toFetch--
      sig <- true
    case j := <-jobs:
      toFetch += j
    case e := <-errs:
      log.Println(e.Error())

      toFetch--
      sig <- true
    }
  }
}

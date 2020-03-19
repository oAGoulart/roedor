package crawler

import (
  "net/url"
  "log"
  "os"
  "sync"

  ."github.com/oagoulart/roedor/pkg/util"
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
  numWorkers uint
  cache      *urlCache
  tokens     Tokens
  output     string
}

func NewCrawler(instances []*url.URL, numWorkers uint, tokens Tokens, output string) *Crawler {
  return &Crawler{
    instances:  instances,
    numWorkers: numWorkers,
    cache:      newUrlCache(),
    tokens:     tokens,
    output:     output,
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

  f, err := os.Create(c.output)
  PanicErr(err)
  defer f.Close()

  for _, link := range c.instances {
    go c.crawl(link, sig, jobs, data, errs)
  }

  for i := 0; uint(i) < c.numWorkers; i++ {
    sig <- true
  }

  toFetch := 1
  for toFetch > 0 {
    select {
    case d := <-data:
      body, err := extract(d.Body, c.tokens)
      LogErr(err)

      // Create CSV row
      row := make([]byte, 1)
      row = []byte("\"")
      row = append(row, []byte(d.Link.String())...)
      row = append(row, []byte("\",\"")...)
      row = append(row, body...)
      row = append(row, []byte("\"\n")...)

      _, err = f.Write(row)
      if !LogErr(err) {
        f.Sync()
      }

      log.Println(d.Link.String(), string(body))

      toFetch--
      sig <- true
    case j := <-jobs:
      toFetch += j
    case e := <-errs:
      LogErr(e)

      toFetch--
      sig <- true
    }
  }
}

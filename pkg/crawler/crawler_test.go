package crawler

import (
  "testing"
  "net/url"
  "log"
)

func TestCrawler(t *testing.T) {
  const numWorkers = 4
  link, err := url.Parse("https://gobyexample.com/")

  if err != nil {
    log.Panicln(err.Error())
  }

  c := NewCrawler(
    []*url.URL{
      link,
    },
    numWorkers,
  )
  c.Start()
}

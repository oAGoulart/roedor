package crawler

import (
  "net/url"
  "testing"

  "github.com/stretchr/testify/assert"
)

// TestCrawler tests crawling job
func TestCrawler(t *testing.T) {
  const numWorkers = 4

  link, err := url.Parse("https://gobyexample.com/")
  assert.Nil(t, err)

  tokens := make(Tokens)
  tokens["p"] = "\n{}"

  c := NewCrawler(
    []*url.URL{
      link,
    },
    numWorkers,
    tokens,
    "./roedor.json",
  )

  // Test will run while links are found
  c.Start()
}

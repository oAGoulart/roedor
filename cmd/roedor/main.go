package main

import (
  "encoding/json"
  "net/url"
  "flag"
  "os"
  "errors"

  "github.com/oagoulart/roedor/pkg/crawler"
  "github.com/oagoulart/roedor/pkg/util"
)

func main() {
	const (
		maxWorkers = 4
	)

  var numWorkers uint
  var siteURL, tokensJSON, outputPath string
  var showHelp bool

  flag.UintVar(&numWorkers, "workers", maxWorkers, "number of parallel `workers` at the same time")
  flag.StringVar(&siteURL, "url", "", "`link` to crawl onto")
  flag.StringVar(&tokensJSON, "tokens", "", "JSON string with `tokens`")
  flag.StringVar(&outputPath, "output", "./roedor.csv", "`path` to output CSV file")
  flag.BoolVar(&showHelp, "help", false, "show help message")
  flag.Parse()

  if showHelp {
  	flag.PrintDefaults()
  	os.Exit(1)
  } else if siteURL == "" {
  	util.FatalErr(errors.New("you need to specify --url. use --help for usage"))
  } else if tokensJSON == "" {
  	util.FatalErr(errors.New("you need to specifty --tokens. use --help for usage"))
  }

  link, err := url.Parse(siteURL)
  util.PanicErr(err)

  var tokens crawler.Tokens
  err = json.Unmarshal([]byte(tokensJSON), &tokens)
  util.PanicErr(err)

  c := crawler.NewCrawler(
    []*url.URL{
      link,
    },
    numWorkers,
    tokens,
    outputPath,
  )
  c.Start()
}

package main

import (
  "encoding/json"
  "net/url"
  "flag"
  "os"
  "errors"

  ."github.com/oagoulart/roedor/pkg/crawler"
  ."github.com/oagoulart/roedor/pkg/util"
)

func main() {
	const (
		NUMWORKERS = 4
	)

  var numWorkers uint
  var siteUrl, tokensJson, outputPath string
  var showHelp bool

  flag.UintVar(&numWorkers, "workers", NUMWORKERS, "number of parallel `workers` per link")
  flag.StringVar(&siteUrl, "url", "", "`link` to crawl into")
  flag.StringVar(&tokensJson, "tokens", "", "JSON string with `tokens`")
  flag.StringVar(&outputPath, "output", "./roedor.csv", "`path` to output CSV file")
  flag.BoolVar(&showHelp, "help", false, "show help message")
  flag.Parse()

  if showHelp {
  	flag.PrintDefaults()
  	os.Exit(1)
  } else if siteUrl == "" {
  	FatalErr(errors.New("You need to specify --url. Use --help for usage."))
  } else if tokensJson == "" {
  	FatalErr(errors.New("You need to specifty --tokens. Use --help for usage."))
  }

  link, err := url.Parse(siteUrl)
  PanicErr(err)

  var tokens Tokens
  err = json.Unmarshal([]byte(tokensJson), &tokens)
  PanicErr(err)

  c := NewCrawler(
    []*url.URL{
      link,
    },
    numWorkers,
    tokens,
    outputPath,
  )
  c.Start()
}

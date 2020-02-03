package crawler

import (
  "net/http"
  "net/url"
  "io/ioutil"
  "html"
  "regexp"
)

func getHtml(link *url.URL) ([]byte, error) {
  resp, err := http.Get(link.String())

  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  return ioutil.ReadAll(resp.Body)
}

func getLinks(link *url.URL) ([]*url.URL, error) {
  link.RawQuery = ""
  body, err := getHtml(link)

  if err != nil {
    return nil, err
  }

  re := regexp.MustCompile("href=\"(?:(https?:\\/\\/[^\"]*))?\"")
  matches := re.FindAllStringSubmatch(string(body), -1)
  links := make([]*url.URL, len(matches))

  for i, match := range matches {
    parsed, err := url.Parse(html.EscapeString(match[1]))

    if err != nil {
      return links, err
    }

    link.RawQuery = ""
    links[i] = parsed
  }

  return links, nil
}

func fetch(link *url.URL) ([]byte, []*url.URL, error) {
  links, err := getLinks(link)

  if err != nil {
    return nil, nil, err
  }

  body, err := getHtml(link)

  if err != nil {
    return nil, nil, err
  }

  return body, links, nil
}

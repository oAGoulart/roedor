package crawler

import (
  "net/http"
  "net/url"
  "io/ioutil"
  "os/exec"
  "encoding/json"
  "html"
  "regexp"
  "fmt"
)

const (
  MARKOUT = "markout_html"
)

type Tokens map[string]string

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
    return nil, links, err
  }

  return body, links, nil
}

func extract(body []byte, tokens Tokens) ([]byte, error) {
  jsonTokens, err := json.Marshal(tokens)
  if err != nil {
    return nil, err
  }
  
  arg1 := fmt.Sprintf("-e %s", body)
  arg2 := fmt.Sprintf("-t %s", jsonTokens)
  arg3 := fmt.Sprintf("-o %s", "html")
  cmd := exec.Command(MARKOUT, arg1, arg2, arg3)

  out, err := cmd.CombinedOutput()
  if err != nil {
    return out, err
  }

  return out, nil
}

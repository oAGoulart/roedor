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
  markout = "markout_html" // markout's cli
)

// Tokens is used to convert JSON tokens into a map type
type Tokens map[string]string

// getHTML uses `link` to make HTTP GET request and return response's body
func getHTML(link *url.URL) ([]byte, error) {
  resp, err := http.Get(link.String())
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  return ioutil.ReadAll(resp.Body)
}

// getLinks extracts every link from `body`
func getLinks(body []byte) ([]*url.URL, error) {
  re := regexp.MustCompile("href=\"(?:(https?:\\/\\/[^\"]*))?\"")
  matches := re.FindAllStringSubmatch(string(body), -1)
  links := make([]*url.URL, len(matches))

  for i, match := range matches {
    parsed, err := url.Parse(html.EscapeString(match[1]))
    if err != nil {
      return links, err
    }

    parsed.RawQuery = ""
    links[i] = parsed
  }

  return links, nil
}

// fetch extracts HTML and links from `link` web page
func fetch(link *url.URL) ([]byte, []*url.URL, error) {
  body, err := getHTML(link)
  if err != nil {
    return nil, nil, err
  }

  links, err := getLinks(body)
  if err != nil {
    return body, nil, err
  }

  return body, links, nil
}

// extract uses Markout to extract content from HTML `body` using `tokens`
func extract(body []byte, tokens Tokens) ([]byte, error) {
  jsonTokens, err := json.Marshal(tokens)
  if err != nil {
    return nil, err
  }
  
  arg1 := fmt.Sprintf("-e %s", body)
  arg2 := fmt.Sprintf("-t %s", jsonTokens)
  arg3 := fmt.Sprintf("-o %s", "html")
  cmd := exec.Command(markout, arg1, arg2, arg3)

  out, err := cmd.CombinedOutput()
  if err != nil {
    return out, err
  }

  return out, nil
}

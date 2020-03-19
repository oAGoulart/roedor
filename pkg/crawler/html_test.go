package crawler

import (
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestHtml_Extract(t *testing.T) {
  body := `
  	<html>
  		<head></head>
  		<body>
  			<h1>best</h1>
  			<p>potato</p>
  		</body>
  	</html>
  `

  tokens := make(Tokens)
  tokens["p"] = "{}"
  tokens["h1"] = "{}"
  
  ext, err := extract([]byte(body), tokens)
  assert.Nil(t, err)

  assert.Equal(t, []byte("best\npotato"), ext, "The extraction is not correct.")
}

# Roedor

[![License](https://img.shields.io/badge/license-MIT-informational.svg)](https://opensource.org/licenses/MIT)
![Go](https://github.com/oAGoulart/roedor/workflows/Go/badge.svg)

A modular web crawler Go module. The purpose of this module is to crawl through web sites and extract data using the python package [**Markout**](https://github.com/oAGoulart/markout). All extracted data is then stored into a CSV file to be analyzed later.

## Usage

Firstly, make sure you fulfill the requirements below:

+ **Python** >= 3.5
+ **Pip** (*Python module*) >= 18.1
+ **Go** >= 1.13.6

To install this package, run the `installing` script on `./scripts`:

```sh
./scripts/installing
```

**NOTE:** This package uses an external Python package, if you want to use (*Roedor*) on your code you must also make sure that package (*Markout*) is also installed (the `installing` script should take care of this).

### Using the CLI

If you want to use the CLI, all you have to do is to call `markout_html` with the following flags (`$GOBIN`must be set for this command to work):

`--workers`: Number of parallel workers at the same time.

`--url`: link to crawl onto.

`--tokens`: JSON string with tokens to be used (see [Markout](https://github.com/oAGoulart/markout) for details).

`--output`: filename of output CSV file (optional).

You may also use the `--help` flag to list all the flags above with help messages.

### Using on your code

If you this want to use this package on your code, you can just import it!
But remember that this package requires external Python packages!

Here's an example of use:

```go
link, err := url.Parse("https://gobyexample.com/")
if err != nil {
  panic("hot potatoes")
}

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

// This will run while links are found
c.Start()
```

---

## Contributions

Feel free to leave your contribution here, I would really appreciate it!
Also, if you have any doubts or troubles using this package just contact me or leave an issue.

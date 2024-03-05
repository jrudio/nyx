Nyx is a Go library intended to help parse the favorited artists listed on the shareable link generated from the Insomniac app. My personal use-case is to keep track of which artists I intend to watch vs artists I ended up watching, but I'm sure there are better, more productive use-cases.

> https://insom.app/\<unique-id\>

or

> https://www.insomniac.com/s/\<id\>/lineup/\<unique-id\>

The following example is also under `cmd/`:

```
package main

import (
  "fmt"
  "github.com/jrudio/nyx"
)

func main() {
  // *inputFile is path to a html file
  data, _ := os.ReadFile(*inputFile)

  lineup, err := nyx.ParseLineup(data)

  fmt.Println(lineup.Artists)

  ...

  insomniacURL := "https://www.insomniac.com/s/<id>/lineup/<unique-id>/"

  lineup, err := nyx.Get(insomniacURL)

  fmt.Println(lineup.Artists)
}

```
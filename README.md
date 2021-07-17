# Args


### Usage

```go
package main

import (
  "fmt"

  "github.com/taybart/args"
)

func main() {
	// Set up app
	app := args.App{
		Name:    "My App",
		Version: "v0.0.1",
		Author:  "Person <real@email.com>",
		About:   "Really cool app for accomplishing stuff",
		Args: map[string]*args.Arg{
			"port": {
				Short:   "p",
				Long:    "port",
				Help:    "Port to listen on",
				Default: 8080,
			},
		},
	}
  app.Parse()
  fmt.Println(app.Get("port").Int())
}
```

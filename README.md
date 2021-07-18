# Args


### Usage

```go
package main

import (
  "fmt"
  "os"

  "github.com/taybart/args"
)


var (
	app = args.App{
		Name:    "My App",
		Version: "v0.0.1",
		Author:  "Tay Bart <taybart@email.com>",
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
)

func main() {
	// Set up app
  app.Parse()
  if err := run(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func run() error {
  fmt.Println(app.Get("port").Int())
  return nil
}
```

```go
package main

import (
  "fmt"
  "os"

  "github.com/taybart/args"
)

const (
  port = "port"
)

var (
	app = args.App{
		Name:    "My App",
		Version: "v0.0.1",
		Author:  "Tay Bart <taybart@email.com>",
		About:   "Really cool app for accomplishing stuff",
		Args: map[string]*args.Arg{
			port: {
				Short:   "p",
				Long:    "port",
				Help:    "Port to listen on",
				Default: 8080,
			},
		},
	}
)

func main() {
	// Set up app
  app.Parse()
  if err := run(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func run() error {
  fmt.Println(app.Get(port).Int())
  return nil
}
```


*Note:* if the variable is supposed to be treated as a boolean, `Default: false` is required 

Reserved flags:

`-h,-help,--help`

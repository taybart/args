# Args

Replicating [clap](https://github.com/clap-rs/clap) in go.

I started learning rust and really liked [clap](https://github.com/clap-rs/clap). I am sure this is a derivitave of something else, blah blah. I thought it would be nice for me to use.

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
  if err := run(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func run() error {
  // Set up app
  err := app.Parse()
  if err != nil {
    return err
  }
  fmt.Println(app.Int("port"))
  return nil
}
```

*Note:* if the variable is supposed to be treated as a boolean, `Default: false` is required 

Reserved flags:

`-h,-help,--help`

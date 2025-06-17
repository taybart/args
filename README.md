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
      "zero": {
        Short:   "z",
        Help:    "If divide by zero is ok",
        Default: false,
      },
      "port": {
        Short:   "p",
        Help:    "Port to listen on",
        Default: 8080,
      },
    },
    // optional usage override function
    UsageFunc: usageOverride,
  }
)

func usageOverride(u args.Usage) {
    // ordered flag slice, default is alphabetical (passed in u)
	cli := []string{
        "zero",
		"port",
	}
	var usage strings.Builder
	usage.WriteString("Here is my usage help override!\n")
	u.BuildFlagString(&usage, cli)
	fmt.Println(usage.String())
}

func main() {
  if err := run(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func run() error {
  // Set up app
  if err := app.Parse(); err != nil {
    // user passed -h or --help
    if errors.Is(err, args.ErrUsageRequested) {
        return nil
    }
    return err
  }
  fmt.Println(app.Int("port"))

  // use go struct
  config := struct {
    Port    int    `arg:"port"`
    Zero    bool    `arg:"zero"`
  }{}
  if err = app.Marshal(&config); err != nil {
    return err
  }
  fmt.Println(config.Port)
  return nil
}
```

_Note:_ if the variable is supposed to be treated as a boolean, `Default: false` is required

If `-h,-help,--help` are not specified in the app definition, one is provided automatically. It will return the error `args.ErrUsageRequested` when the app is parsed. Make sure to add `Default: false` if you do override the help flag

package args_test

import (
	"errors"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/taybart/args"
)

func TestNewApp(t *testing.T) {
	is := is.New(t)

	// Set up app
	app := args.App{
		Name:    "My App",
		Version: "v0.0.1",
		Author:  "tester mctestyface <tmct@email.com>",
		About:   "Really cool app for accomplishing stuff",
		Args: map[string]*args.Arg{
			"print": {
				Short:   "p",
				Long:    "print",
				Help:    "Prints a cool message",
				Default: false,
			},
			"message": {
				Short:   "m",
				Long:    "message",
				Help:    "Sets a cool message",
				Default: "COOL!",
			},
		},
	}

	// Add cli args
	os.Args = []string{"./test", "--message=test", "-p"}

	err := app.Parse()
	is.NoErr(err)

	is.True(app.Get("print").Bool())               // should default to false, boolflag to true
	is.True(app.Get("message").String() == "test") // set through os.Args
	app.Usage()

}

func TestDuplicateFlags(t *testing.T) {
	is := is.New(t)

	// Set up app
	app := args.App{
		Name:    "My App",
		Version: "v0.0.1",
		Author:  "tester mctestyface <tmct@email.com>",
		About:   "Really cool app for accomplishing stuff",
		Args: map[string]*args.Arg{
			"print": {
				Short:   "p",
				Long:    "print",
				Help:    "Prints a cool message",
				Default: false,
			},
			"port": {
				Short:   "p",
				Long:    "port",
				Help:    "Port to listen on",
				Default: 8080,
			},
		},
	}

	// Add cli args
	os.Args = []string{"./test", "--message=test"}

	err := app.Parse()
	is.True(errors.Is(err, args.ErrDuplicateKey))
}

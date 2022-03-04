package args

import (
	"errors"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestNewApp(t *testing.T) {
	is := is.New(t)
	// log.SetLevel(log.VERBOSE)

	// Add cli args
	os.Args = []string{"./test", "--message=test", "-p", "-n=69"}

	// Set up app
	app := App{
		Name:    "My App",
		Version: "v0.0.1",
		Author:  "tester mctestyface <tmct@email.com>",
		About:   "Really cool app for accomplishing stuff",
		Args: map[string]*Arg{
			"print": {
				Short:   "p",
				Help:    "Allows the cool message to be printed",
				Default: false,
			},
			"message": {
				Short:   "m",
				Help:    "Sets a cool message",
				Default: "COOL!",
			},
			"nums": {
				Short:   "n",
				Help:    "A really fun number",
				Default: 0,
			},
		},
	}

	err := app.Parse()
	is.NoErr(err)

	is.True(app.Get("print").Bool())
	is.True(app.Get("message").String() == "test")
	is.True(app.Get("nums").Int() == 69)

}

func TestDuplicateFlags(t *testing.T) {
	is := is.New(t)

	// Set up app
	app := App{
		Name:    "My App",
		Version: "v0.0.1",
		Author:  "tester mctestyface <tmct@email.com>",
		About:   "Really cool app for accomplishing stuff",
		Args: map[string]*Arg{
			"print": {
				Short:   "p",
				Help:    "Prints a cool message",
				Default: false,
			},
			"port": {
				Short:   "p",
				Help:    "Port to listen on",
				Default: 8080,
			},
		},
	}

	// Add cli args
	os.Args = []string{"./test", "--message=test"}

	err := app.Parse()
	is.True(errors.Is(err, ErrDuplicateKey))
}

func TestRequiredFlags(t *testing.T) {
	is := is.New(t)
	// log.SetLevel(log.DEBUG)

	// Set up app
	app := App{
		Name:    "My App",
		Version: "v0.0.1",
		Author:  "tester mctestyface <tmct@email.com>",
		About:   "Really cool app for accomplishing stuff",
		Args: map[string]*Arg{
			"cool": {
				Short:    "c",
				Help:     "Makes your program cool af",
				Required: true,
			},
			"port": {
				Short:   "p",
				Help:    "Port to listen on",
				Default: 8080,
			},
		},
	}

	// Add cli args
	os.Args = []string{"./test", "-p=8080"}

	err := app.Parse()
	is.True(errors.Is(err, ErrMissingRequired))
}

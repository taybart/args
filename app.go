package args

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/taybart/log"
)

var (
	ErrDuplicateKey = errors.New("duplicate keys")
)

type Semver string

/*********************************
 *************** App *************
 *********************************/
type App struct {
	Name    string          `json:"name,omitempty"`
	Version Semver          `json:"version,omitempty"`
	Author  string          `json:"author,omitempty"`
	About   string          `json:"about,omitempty"`
	Args    map[string]*Arg `json:"args,omitempty"`
}

func (a *App) Import(app App) App {
	if a.Args == nil {
		a.Args = make(map[string]*Arg)
	}
	// prefer Name,Version,Author,About from original
	for k, v := range app.Args {
		a.Args[k] = v
	}
	return *a
}

func (a *App) Parse() error {
	err := a.Validate()
	if err != nil {
		return err
	}

	fs := flag.NewFlagSet(a.Name, flag.ContinueOnError)
	for k, arg := range a.Args {
		arg.Name = k
		arg.init(fs)
	}

	a.CreateConfig()

	// test comment
	buf := bytes.NewBuffer([]byte{}) // suppress default output
	fs.SetOutput(buf)
	err = fs.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	return nil
}

func (a *App) Validate() error {
	defined := map[string]string{}
	issues := []string{}
	// check duplicate flags
	for k, arg := range a.Args {
		if defined[arg.Short] != "" {
			issue := fmt.Sprintf("flag %s already used in arg %s", arg.Short, defined[arg.Short])
			issues = append(issues, issue)
			continue
		}
		if defined[arg.Long] != "" {
			issue := fmt.Sprintf("flag %s already used in arg %s", arg.Long, defined[arg.Long])
			issues = append(issues, issue)
			continue
		}
		arg.Name = k
		arg.validate()
		if arg.Long != "" {
			defined[arg.Long] = k
		}
		if arg.Short != "" {
			defined[arg.Short] = k
		}
	}
	if len(issues) > 0 {
		err := ErrDuplicateKey
		return fmt.Errorf("%v %w", issues, err)
	}
	return nil
}

func (a *App) Usage() {
	var usage strings.Builder
	fmt.Fprintf(&usage, "%s%s%s %s [option]\n  %s\n", log.Blue, a.Name, log.Reset, os.Args[0], a.About)
	for k, arg := range a.Args {
		fmt.Fprintf(&usage, "    %s%s%s: -%s, --%s\n\t%s\n", log.Blue, k, log.Reset, arg.Short, arg.Long, arg.Help)
	}
	fmt.Println(usage.String())
}

func (a *App) Get(key string) *Arg {
	return a.Args[key]
}

func (a *App) String(key string) string {
	return a.Args[key].String()
}

func (a *App) CreateConfig() (string, error) {
	b, err := json.Marshal(a.Args)
	return string(b), err
}

func (a *App) ConfigRead() {}

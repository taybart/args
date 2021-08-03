package args

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/taybart/log"
)

var (
	ErrDuplicateKey    = errors.New("duplicate keys")
	ErrMissingRequired = errors.New("missing required keys")
)

var (
	flagRx = regexp.MustCompile(`(?:-+)([[:alnum:]-_]+)(?:=| )?(.*)?`)
	// use CSL in arg to split and parse, this is stupid
	cslRx = regexp.MustCompile("(.*),(.*),(.*),(.*),(.*),(.*),(.*),(.*),(.*)")
)

type Semver string

func ToSemver(in string) Semver {
	// TODO maybe use https://github.com/Masterminds/semver
	// or be a bad boy and write one
	return Semver(in)
}

/*********************************
 *************** App *************
 *********************************/
type App struct {
	Name          string          `json:"name,omitempty"`
	Version       Semver          `json:"version,omitempty"`
	Author        string          `json:"author,omitempty"`
	About         string          `json:"about,omitempty"`
	ExitOnFailure bool            `json:"exit_on_failure"`
	Args          map[string]*Arg `json:"args,omitempty"`
	// ArgsRequired bool            `json:"args_required,omitempty"`
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

	// make names and bool values
	for k, arg := range a.Args {
		arg.Name = k
		arg.init()
	}

	// a.CreateConfig()

	log.Verbosef("os.Args: %v\n", os.Args)
	for _, v := range os.Args {
		log.Debugf("trying to match: %v\n", v)
		matches := flagRx.FindAllStringSubmatch(v, -1) // regexp each flag
		log.Debugf("result: %v\n", matches)
		for _, arg := range a.Args {
			log.Debugf("arg: %s %v\n", arg.Name, matches)
			if len(matches) > 0 { // arg exists
				name := matches[0][1]
				if arg.Short != name && arg.Long != name {
					log.Debug("didn't match", name)
					continue
				}
				if arg.isBool {
					if name != "" {
						arg.SetBool(true)
					}
					continue
				}
				if len(matches[0]) > 0 { // arg was set
					value := matches[0][2]
					if arg.Short == name || arg.Long == name {
						err = arg.Set(value)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	for _, arg := range a.Args {
		if arg.Required && !arg.IsSet() {
			a.Usage()
			return ErrMissingRequired
		}
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

// func (a *App) Set(key string, value interface{}) *Arg {
// 	arg := a.Args[key]
// 	arg.Set(value)
// 	a.Args[key] = arg
// 	return a.Args[key]
// }

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

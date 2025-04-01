package args

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/taybart/log"
)

var (
	ErrDuplicateKey    = errors.New("duplicate keys")
	ErrMissingRequired = errors.New("missing required keys")
	ErrUsageRequested  = errors.New("usage requested")
)

var (
	flagRx = regexp.MustCompile(`(?:-+)([[:alnum:]-_]+)(?:=| )?(.*)?`)
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
	App           interface{}     // marshal result
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

	log.Verbosef("os.Args: %v\n", os.Args)
	for i, v := range os.Args {
		log.Debugf("trying to match: %v\n", v)
		matches := flagRx.FindAllStringSubmatch(v, -1) // regexp each flag
		log.Debugf("result: %v\n", matches)
		for _, arg := range a.Args {
			log.Debugf("arg: %s %v\n", arg.Name, matches)
			if len(matches) > 0 { // arg exists
				name := matches[0][1]
				if name == "h" || name == "help" {
					a.Usage()
					return ErrUsageRequested
				}
				if arg.Short != name && arg.Name != name {
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
					if arg.Short == name || arg.Name == name {
						value := matches[0][2]
						if value != "" {
							err = arg.Set(value)
							if err != nil {
								return err
							}
						} else {
							i++
							next := os.Args[i]
							if next[0] == '-' {
								return fmt.Errorf("flag given but argument (%s) not set", arg.Name)
							}
							err = arg.Set(next)
							if err != nil {
								return err
							}
							continue
						}
					}
				}

			}
		}
	}
	var req string
	for _, arg := range a.Args {
		if arg.Required && !arg.IsSet() {
			if req == "" {
				req = arg.Name
			} else {
				req = fmt.Sprintf("%s,%s", req, arg.Name)
			}
		}
	}

	if len(req) > 0 {
		fmt.Printf("Missing required arguments: %s%s%s\n", log.Red, req, log.Reset)
		a.Usage()
		return ErrMissingRequired
	}

	if err := a.MarshalInternal(); err != nil {
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
		arg.Name = k
		if defined[arg.Name] != "" {
			issue := fmt.Sprintf("flag %s already used in arg %s", arg.Name, defined[arg.Name])
			issues = append(issues, issue)
			continue
		}
		arg.validate()
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
	// fmt.Fprintf(&usage, "%s%s%s: %s%s%s\n%s [options]\n",
	// log.Blue, a.Name, log.Reset, log.Gray, a.About, log.Reset, os.Args[0])
	l := len(a.Args)
	for _, arg := range a.Args {
		l--
		fmt.Fprintf(&usage, "    --%s, -%s:\n\t%s", arg.Name, arg.Short, arg.Help)
		if l > 0 {
			fmt.Fprintf(&usage, "\n")
		}

	}
	fmt.Println(usage.String())
}

func (a *App) MarshalInternal() error {
	if a.App != nil {
		if err := a.Marshal(a.App); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Marshal(i interface{}) error {

	v := reflect.ValueOf(i).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}

	// TODO: return better casting errors
	for index := 0; index < v.NumField(); index++ {
		field := v.Type().Field(index)

		tag := field.Tag.Get("arg")
		typ := field.Type.Name()

		log.Debugf("%v (%v), tag: %v\n", field.Name, typ, a.Get(tag))

		f := v.Field(index)
		switch typ {
		case "int":
			f.Set(reflect.ValueOf(a.Int(tag)))
		case "bool":
			f.Set(reflect.ValueOf(a.Is(tag)))
		case "string":
			f.Set(reflect.ValueOf(a.String(tag)))
		default:
			return fmt.Errorf("unknown type %s", typ)
		}

	}
	return nil
}

func (a *App) Get(key string) *Arg {
	return a.Args[key]
}

func (a *App) String(key string) string {
	return a.Args[key].String()
}

func (a *App) Int(key string) int {
	return a.Args[key].Int()
}

func (a *App) Bool(key string) bool {
	return a.Args[key].Bool()
}

func (a *App) Is(key string) bool {
	return a.Args[key].Bool()
}

func (a *App) True(key string) bool {
	return a.Args[key].Bool()
}

func (a *App) File(key string) []byte {
	return a.Args[key].File()
}

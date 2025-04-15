package args

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
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
	overrideHelp  bool
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

	if err := a.Validate(); err != nil {
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
				if (name == "h" || name == "help") && !a.overrideHelp {
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
							if i > len(os.Args)-1 {
								// this should also take into account if the user just didn't provide an arg
								// and it was the last in the list
								return fmt.Errorf("argument was not provided with a value, this might mean that it is a boolean and Default was not specified")
							}
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
		if arg.Short == "h" || arg.Name == "help" {
			a.overrideHelp = true
		}
	}
	if len(issues) > 0 {
		err := ErrDuplicateKey
		return fmt.Errorf("%v %w", issues, err)
	}
	return nil
}

func (a *App) Usage() {

	// Sort args in alphabetical order
	keys := make([]string, 0, len(a.Args))
	for key := range a.Args {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var usage strings.Builder
	for i, key := range keys {
		arg := a.Args[key]
		if arg == nil {
			continue
		}
		fmt.Fprintf(&usage, "    --%s", arg.Name)
		if arg.Short != "" {
			fmt.Fprintf(&usage, ", -%s", arg.Short)
		}
		fmt.Fprintf(&usage, ":\n\t%s", arg.Help)
		if i < len(a.Args)-1 {
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
		return fmt.Errorf("marshal interface must be a pointer in order to assign")
	}

	// TODO: return better casting errors
	for index := 0; index < v.NumField(); index++ {
		field := v.Type().Field(index)

		tag := field.Tag.Get("arg")
		typ := field.Type.Name()

		log.Debugf("%v (%v), tag: %v\n", field.Name, typ, a.Get(tag))
		if a.Get(tag) == nil { // we don't have that tag
			continue
		}

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

package args

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
)

/*********************************
 *************** Arg *************
 *********************************/

type Arg struct {
	// in config if this is not defined long is used
	// could also be embedded (ex. logs.verbose)
	Name     string      `json:"name,omitempty"`
	Short    string      `json:"short,omitempty"`
	Long     string      `json:"long,omitempty"`
	Help     string      `json:"help,omitempty"`
	Required bool        `json:"required,omitempty"`
	Default  interface{} `json:"default,omitempty"`
	Type     string      `json:"type,omitempty"`
	value    *string
	wasSet   *bool
	isBool   bool
}

func (arg Arg) IsBoolFlag() bool {
	return arg.isBool
}

func (arg Arg) Set(s string) error {
	*arg.wasSet = true
	*arg.value = s
	return nil
}

func (arg Arg) IsSet(s string) bool {
	return *arg.wasSet
}

func (arg *Arg) Bool() bool {
	if !*arg.wasSet {
		if arg.Default == nil {
			return false
		}
		return arg.Default.(bool)
	}
	return *arg.value == "true"
}

func (arg *Arg) Print() {
	fmt.Printf("%s=%v wasSet=%v\n", arg.Name, *arg.value, *arg.wasSet)
}

func (arg *Arg) Int() int {
	if !*arg.wasSet {
		if arg.Default == nil {
			return 0
		}
		return arg.Default.(int)
	}
	i, err := strconv.Atoi(*arg.value)
	if err != nil {
		panic(fmt.Sprintf("flag provided for %s could not be converted to int", arg.Name))
	}
	return i
}
func (arg *Arg) String() string {
	if arg.isBool {
		if arg.Bool() {
			return "true"
		}
		return "false"
	}
	if !*arg.wasSet {
		if arg.Default == nil {
			return ""
		}
		return arg.Default.(string)
	}
	return *arg.value
}
func (arg *Arg) validate() error {
	if arg.Long == "" && arg.Short == "" {
		return fmt.Errorf("arg requires a flag name")
	}
	return nil
}
func (arg *Arg) init(fs *flag.FlagSet) error {
	arg.isBool = reflect.TypeOf(arg.Default).String() == "bool"

	// init pointers
	str := ""
	arg.value = &str
	ws := false
	arg.wasSet = &ws

	if arg.Long != "" {
		fs.Var(arg, arg.Long, arg.Help)
	}
	if arg.Short != "" {
		fs.Var(arg, arg.Short, arg.Help)
	}
	return nil
}

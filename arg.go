package args

import (
	"fmt"
	"reflect"

	"github.com/taybart/log"
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
	CSL      bool        `json:"CSL,omitempty"` // comma seperated
	Default  interface{} `json:"default,omitempty"`
	Type     string      `json:"type,omitempty"`
	value    interface{}
	wasSet   bool
	isBool   bool
}

func (arg Arg) IsBoolFlag() bool {
	return arg.isBool
}

func (arg *Arg) Set(value interface{}) error {
	log.Verbosef("setting %s%+v%s => %s%+v%s\n",
		log.BoldGreen, arg.Name, log.Reset,
		log.BoldBlue, value, log.Reset,
	)
	arg.wasSet = true
	arg.value = value

	log.Verbosef("set => %+v\n", arg.value)
	return nil
}

func (arg *Arg) SetBool(value bool) error {
	arg.wasSet = true
	// if v, ok := value.(string); ok {
	log.Verbosef("setting %sbool%s %s%+v%s => %s%+v%s\n",
		log.Yellow, log.Reset,
		log.BoldGreen, arg.Name, log.Reset,
		log.BoldBlue, value, log.Reset,
	)
	arg.value = value
	return nil
	// }
	// panic("setting string that is not a string")
}

func (arg *Arg) SetString(value string) error {
	arg.wasSet = true
	// if v, ok := value.(string); ok {
	log.Verbosef("setting %sstring%s %s%+v%s => %s%+v%s\n",
		log.Yellow, log.Reset,
		log.BoldGreen, arg.Name, log.Reset,
		log.BoldBlue, value, log.Reset,
	)
	arg.value = value
	return nil
	// }
	// panic("setting string that is not a string")
}

func (arg *Arg) SetInt(value int) error {
	arg.wasSet = true
	// if v, ok := value.(string); ok {
	log.Verbosef("setting %sint%s %s%+v%s => %s%+v%s\n",
		log.Yellow, log.Reset,
		log.BoldGreen, arg.Name, log.Reset,
		log.BoldBlue, value, log.Reset,
	)
	arg.value = value
	return nil
	// }
	// panic("setting string that is not a string")
}

func (arg Arg) IsSet() bool {
	if arg.Default != nil {
		return true
	}
	return arg.wasSet
}

func (arg *Arg) Bool() bool {
	if !arg.wasSet {
		if arg.Default == nil {
			return false
		}
		return arg.Default.(bool)
	}
	return arg.value.(bool)
}

func (arg *Arg) Print() {
	fmt.Printf("%s=%v wasSet=%v\n", arg.Name, arg.value, arg.wasSet)
}

func (arg *Arg) Int() int {
	if !arg.wasSet {
		if arg.Default == nil {
			return 0
		}
		return arg.Default.(int)
	}
	// i, err := strconv.Atoi(arg.value)
	// if err != nil {
	// 	panic(fmt.Sprintf("flag provided for %s could not be converted to int", arg.Name))
	// }
	return arg.value.(int)
}

func (arg *Arg) String() string {
	if arg.isBool {
		if arg.Bool() {
			return "true"
		}
		return "false"
	}
	if !arg.wasSet {
		if arg.Default == nil {
			return ""
		}
		return fmt.Sprintf("%v", arg.Default)
	}
	switch arg.value.(type) {
	case bool:
		return fmt.Sprintf("%t", arg.Bool())
	case int:
		return fmt.Sprintf("%d", arg.Int())
	case string:
		if s, ok := arg.value.(string); ok {
			return s
		}
	default:
		return arg.Default.(string)
	}
	return ""
}

func (arg *Arg) validate() error {
	if arg.Long == "" && arg.Short == "" {
		return fmt.Errorf("arg requires a flag name")
	}
	return nil
}

func (arg *Arg) init() error {
	arg.isBool = false
	if arg.Default != nil {
		log.Debugf("value %s is type %s\n", arg.Name, reflect.TypeOf(arg.Default).String())
		switch arg.Default.(type) {
		case bool:
			arg.isBool = true
			arg.value = arg.Default.(bool)
		case string:
			arg.value = arg.Default.(string)
		case int:
			arg.value = arg.Default.(int)
		}
		// log.Debug(arg.isBool)
	}

	// // init pointers
	// // str := ""
	// // arg.value = &str
	// // ws := false
	// // arg.wasSet = ws

	// if arg.Long != "" {
	// 	fs.Var(arg, arg.Long, arg.Help)
	// }
	// if arg.Short != "" {
	// 	fs.Var(arg, arg.Short, arg.Help)
	// }
	return nil
}

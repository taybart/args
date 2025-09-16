package args

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/taybart/log"
)

/*********************************
 *************** Arg *************
 *********************************/

type Arg struct {
	// in config if this is not defined long is used
	// could also be embedded (ex. logs.verbose)
	Name             string      `json:"name,omitempty"`
	Short            string      `json:"short,omitempty"`
	Help             string      `json:"help,omitempty"`
	Required         bool        `json:"required,omitempty"`
	Default          interface{} `json:"default,omitempty"`
	Provided         bool        `json:"provided,omitempty"`
	DoesNotNeedValue bool        `json:"doesNotNeedValue,omitempty"`
	Type             string      `json:"type,omitempty"`
	value            interface{}
	wasSet           bool
	isBool           bool
	// new
	Before func()
	After  func()
}

func (arg Arg) IsBoolFlag() bool {
	return arg.isBool
}

func (arg *Arg) Set(value interface{}) error {
	log.Verbosef("arg.Set(%s%+v%s => %s%+v%s)\n",
		log.BoldGreen, arg.Name, log.Reset,
		log.BoldBlue, value, log.Reset,
	)
	if arg.Type != "" {
		// TODO: move to enum
		switch arg.Type {
		case "int":
			arg.SetInt(value.(int))
		case "bool":
			arg.SetBool(value.(bool))
		case "string":
			arg.SetString(value.(string))
		default:
			arg.wasSet = true
			arg.value = value
			log.Verbosef("setting %s%+v%s => %s%+v%s\n",
				log.BoldGreen, arg.Name, log.Reset,
				log.BoldBlue, value, log.Reset,
			)
		}
	}
	switch arg.Default.(type) {
	case int:
		arg.SetInt(value)
	case bool:
		arg.SetBool(value)
	case string:
		arg.SetString(value)
	default:
		arg.wasSet = true
		arg.value = value
		log.Verbosef("setting %s%+v%s => %s%+v%s\n",
			log.BoldGreen, arg.Name, log.Reset,
			log.BoldBlue, value, log.Reset,
		)
	}

	return nil
}

func (arg *Arg) SetBool(value interface{}) error {
	arg.wasSet = true
	log.Verbosef("setting %s%+v%s =>  %sbool%s %s%+v%s\n",
		log.BoldGreen, arg.Name, log.Reset,
		log.Yellow, log.Reset,
		log.BoldBlue, value, log.Reset,
	)
	if v, ok := value.(bool); ok {
		arg.value = v
		return nil
	}
	// assume string
	arg.value = value.(string) == "true"
	return nil
}

func (arg *Arg) SetString(value interface{}) error {
	arg.wasSet = true
	log.Verbosef("setting %s%+v%s =>  %sstring%s %s%+v%s\n",
		log.BoldGreen, arg.Name, log.Reset,
		log.Yellow, log.Reset,
		log.BoldBlue, value, log.Reset,
	)
	arg.value = value.(string)
	return nil
}

func (arg *Arg) SetInt(value interface{}) error {
	arg.wasSet = true
	log.Verbosef("setting %s%+v%s =>  %sint%s %s%+v%s\n",
		log.BoldGreen, arg.Name, log.Reset,
		log.Yellow, log.Reset,
		log.BoldBlue, value, log.Reset,
	)

	if v, ok := value.(int); ok {
		arg.value = v
		return nil
	}
	// assume string
	v, err := strconv.Atoi(value.(string))
	if err != nil {
		return err
	}
	arg.value = v
	return nil
}

// UserSet check if the user provided a value
func (arg Arg) UserSet() bool {
	if arg.Default != nil {
		return arg.Default != arg.value
	}
	return arg.wasSet
}

// IsSet check if the arg has a value (including being set by default)
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

func (arg *Arg) File() []byte {
	fn := arg.String()
	b, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	return b
}

func (arg *Arg) validate() error {
	if arg.Name == "" && arg.Short == "" {
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
	}
	return nil
}

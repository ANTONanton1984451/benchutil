package cli

import "flag"

// CmdFlag интерфейс для корректной работы help команды и корректного парса флагов командной строки
type CmdFlag interface {
	bind(fs *flag.FlagSet)

	name() string
	usage() string
	defaultVal() interface{}
}

type IntFlag struct {
	Name        string
	Destination *int
	Default     int
	Usage       string
}

func (f IntFlag) bind(fs *flag.FlagSet) {
	fs.IntVar(f.Destination, f.Name, f.Default, f.Usage)
}
func (f IntFlag) defaultVal() interface{} {
	return f.Default
}

func (f IntFlag) usage() string {
	return f.Usage
}

func (f IntFlag) name() string {
	return f.Name
}

type StringFlag struct {
	Name        string
	Destination *string
	Default     string
	Usage       string
}

func (f StringFlag) bind(fs *flag.FlagSet) {
	fs.StringVar(f.Destination, f.Name, f.Default, f.Usage)
}

func (f StringFlag) defaultVal() interface{} {
	return f.Default
}

func (f StringFlag) usage() string {
	return f.Usage
}

func (f StringFlag) name() string {
	return f.Name
}

type BoolFlag struct {
	Name        string
	Destination *bool
	Default     bool
	Usage       string
}

func (f BoolFlag) bind(fs *flag.FlagSet) {
	fs.BoolVar(f.Destination, f.Name, f.Default, f.Usage)
}

func (f BoolFlag) defaultVal() interface{} {
	return f.Default
}

func (f BoolFlag) usage() string {
	return f.Usage
}

func (f BoolFlag) name() string {
	return f.Name
}

package args

import (
	"fmt"
	"strings"

	"github.com/alexflint/go-arg"
)

// GitCommit for use when compiling
var GitCommit string

// GitLastTag for use when compiling
var GitLastTag string

// GitExactTag for use when compiling
var GitExactTag string

// Date for use when compiling
var Date string

// Args commandline arguments
type args struct {
	Files   []string `arg:"positional"`
	Number  bool     `arg:"-n,--number" default:"false"`
	Lines   int      `arg:"-l,--lines" default:"0"`
	Supress bool     `arg:"-s,--supress" help:"suppress output of non-printable characters" default:"false"`
	Pretty  bool     `arg:"-p,--pretty" default:"false"`
}

// Args command line args
var Args args

func init() {
	arg.MustParse(&Args)
	arg.MustParse(&Args)
}

func (args) Description() string {
	return `This is an implementation of the Plan 9 p  utility.
It is basically a simple less command.
`
}

// Version version information
func (args) Version() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("commit: %8s\n", GitCommit))
	sb.WriteString(fmt.Sprintf("tag: %10s\n", GitExactTag))
	sb.WriteString(fmt.Sprintf("date: %23s\n", Date))

	return sb.String()
}

package logging

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/pquerna/ffjson/ffjson"
)

func init() {
	fmt := &logrus.TextFormatter{
		ForceColors:      true,
		DisableSorting:   true,
		DisableTimestamp: true,
	}

	logrus.SetFormatter(fmt)
}

func SetDebug(enabled bool) {
	if enabled {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func Printf(format string, a ...interface{}) {
	fmt.Fprintf(logrus.StandardLogger().Out, format, a...)
}

func Println(a ...interface{}) {
	fmt.Fprintln(logrus.StandardLogger().Out, a...)
}

func Print(a ...interface{}) {
	fmt.Fprint(logrus.StandardLogger().Out, a...)
}

func DumpUnmarshaled(descr string, in []byte) {
	if logrus.GetLevel() < logrus.DebugLevel {
		return
	}

	var res interface{}
	if err := ffjson.Unmarshal(in, &res); err != nil {
		panic("DumpUnmarshaled: unable to unmarshal input")
	}

	Dump(descr, res)
}

func DumpJSON(descr string, in interface{}) {
	if logrus.GetLevel() < logrus.DebugLevel {
		return
	}

	out, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		panic("DumpJSON: unable to marshal input")
	}

	Printf("%s ------------------------- dump start ---------------------------------------\n", descr)
	Println(out)
	Printf("%s -------------------------  dump end  ---------------------------------------\n\n", descr)
}

func Dump(descr string, in interface{}) {
	if logrus.GetLevel() < logrus.DebugLevel {
		return
	}

	Printf("%s ------------------------- dump start ---------------------------------------\n", descr)
	spew.Fdump(logrus.StandardLogger().Out, in)
	Printf("%s -------------------------  dump end  ---------------------------------------\n\n", descr)
}

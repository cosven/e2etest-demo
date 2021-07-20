package ginkgo_helper

import (
	"flag"
	"os"
)

func ParameterizedGinkgoContainer(name string, flagSetInit func(flagSet *flag.FlagSet)) bool {
	// The name should not be conflict with other command.
	// TODO(cosven): add validation for name.
	flagSet := flag.NewFlagSet(name, flag.ExitOnError)
	flagSetInit(flagSet)

	// TODO(cosven): I think there exists a library doing the follow parsing.
	var args []string
	var flagSetArgs []string
	foundDelimiter := false
	for _, arg := range os.Args[1:] {
		if !foundDelimiter {
			if arg == name {
				foundDelimiter = true
				continue
			}
		}
		if foundDelimiter {
			flagSetArgs = append(flagSetArgs, arg)
		} else {
			args = append(args, arg)
		}
	}
	flagSet.Parse(flagSetArgs)
	return true
}

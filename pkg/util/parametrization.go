package util

import (
	"flag"
	"fmt"
	"os"
)

func ParameterizedTestCase(name string, flagSetInit func(flagSet *flag.FlagSet)) bool {
	flags := flag.NewFlagSet(name, flag.ExitOnError)
	flagSetInit(flags)

	Parse(name, flag.Args(), flags)

	return true
}

func Parse(name string, args []string, flags *flag.FlagSet) {
	// TODO(cosven): The name should not be conflict with other command. add validation for name.
	if len(args) < 1 {
		fmt.Println("you should specify the name of the test in command line arguments. make sure you properly --focus and --regexScansFilePath on this test as well!")
		os.Exit(1)
	}

	if args[0] == name {
		err := flags.Parse(args[1:])
		if err != nil {
			panic(err)
		}
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

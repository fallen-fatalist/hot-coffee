package flag

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Global flags
var (
	StoragePath = "data"
	Port        = 4000
)

func Parse(args []string) (err error) {
	for _, arg := range args {
		if arg == "--help" {
			PrintHelp()
			os.Exit(0)
		}
	}

	for flagIdx := 0; flagIdx < len(args); flagIdx += 2 {
		flagName, flagValue := args[flagIdx], args[flagIdx+1]
		switch strings.TrimPrefix(flagName, "--") {
		case "port":
			Port, err = strconv.Atoi(flagValue)
			if err != nil {
				return fmt.Errorf("error while parsing the port: %w", err)
			} else if Port < 1024 || Port > 65535 {
				return fmt.Errorf("incorrect range port, port must me between 1024 and 65535")
			}
		case "dir":
			StoragePath = flagValue
		default:
			return fmt.Errorf("unknown flag: %s", flagName)
		}
	}

	return nil
}

func PrintHelp() {
	fmt.Println("Simple Storage Service.")
	fmt.Println("")
	fmt.Println("**Usage:**")
	fmt.Println("\ttriple-s [--port <N>] [--dir <S>]")
	fmt.Println("\ttriple-s --help")
	fmt.Println("")
	fmt.Println("**Options:**")
	fmt.Println("- --help     Show this screen.")
	fmt.Println("- --port N   Port number")
	fmt.Println("- --dir S    Path to the directory")
}

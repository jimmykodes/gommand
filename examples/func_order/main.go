package main

import (
	"fmt"
	"os"

	"github.com/jimmykodes/gommand"
)

var (
	rootCmd = &gommand.Command{
		Name:              "preruns",
		PersistentPreRun:  func(*gommand.Context) error { fmt.Println("root persistent pre run"); return nil },
		PersistentPostRun: func(*gommand.Context) error { fmt.Println("root persistent post run"); return nil },
	}
	subCmd = &gommand.Command{
		Name:              "sub",
		PersistentPreRun:  func(*gommand.Context) error { fmt.Println("sub persistent pre run"); return nil },
		PersistentPostRun: func(*gommand.Context) error { fmt.Println("sub persistent post run"); return nil },
	}
	finalCmd = &gommand.Command{
		Name:   "final",
		PreRun: func(*gommand.Context) error { fmt.Println("final pre run"); return nil },
		Run: func(*gommand.Context) error {
			fmt.Println("final running")
			return nil
		},
		PostRun: func(context *gommand.Context) error { fmt.Println("final post run"); return nil },
	}
	finalErrCmd = &gommand.Command{
		Name:   "final-err",
		PreRun: func(*gommand.Context) error { fmt.Println("final pre run"); return nil },
		Run: func(*gommand.Context) error {
			fmt.Println("final running")
			return fmt.Errorf("some error")
		},
		PostRun: func(context *gommand.Context) error { fmt.Println("final post run"); return nil },
	}
	finalErrDeferredCmd = &gommand.Command{
		Name:      "final-err-deferred",
		DeferPost: true,
		PreRun:    func(*gommand.Context) error { fmt.Println("final pre run"); return nil },
		Run: func(*gommand.Context) error {
			fmt.Println("final running")
			return fmt.Errorf("some error")
		},
		PostRun: func(context *gommand.Context) error { fmt.Println("final post run"); return nil },
	}
	finalPanicCmd = &gommand.Command{
		Name:   "final-panic",
		PreRun: func(*gommand.Context) error { fmt.Println("final pre run"); return nil },
		Run: func(*gommand.Context) error {
			fmt.Println("final running")
			panic("final-panic")
		},
		PostRun: func(context *gommand.Context) error { fmt.Println("final post run"); return nil },
	}
	finalPanicDeferredCmd = &gommand.Command{
		Name:      "final-panic-deferred",
		DeferPost: true,
		PreRun:    func(*gommand.Context) error { fmt.Println("final pre run"); return nil },
		Run: func(*gommand.Context) error {
			fmt.Println("final running")
			panic("final-panic")
		},
		PostRun: func(context *gommand.Context) error { fmt.Println("final post run"); return nil },
	}
)

func init() {
	rootCmd.SubCommand(subCmd)
	subCmd.SubCommand(finalCmd)
	subCmd.SubCommand(finalErrCmd)
	subCmd.SubCommand(finalErrDeferredCmd)
	subCmd.SubCommand(finalPanicCmd)
	subCmd.SubCommand(finalPanicDeferredCmd)
}

// `./func_order sub final` returns
//
// root persistent pre run
// sub persistent pre run
// final pre run
// final running
// final post run
// sub persistent post run
// root persistent post run
//
// `./func_order sub final-err` returns
// root persistent pre run
// sub persistent pre run
// final pre run
// final running
// some error
// exit status 1
//
// `./func_order sub final-err-deferred` returns
// root persistent pre run
// sub persistent pre run
// final pre run
// final running
// final post run
// sub persistent post run
// root persistent post run
// some error
// exit status 1
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

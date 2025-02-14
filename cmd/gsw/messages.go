package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vertgenlab/gonomics/fileio"
	"github.com/vertgenlab/gonomics/genomeGraph"
	"github.com/vertgenlab/gonomics/giraf"
)

// TODO: View feature needs some work.
type ViewExe struct {
	Cmd       *flag.FlagSet
	GirafFile string
	SamFile   string
}

var extendHelpMsg *flag.FlagSet = flag.NewFlagSet("help", flag.ExitOnError)

func helpMessage() {
	fmt.Printf(
		"  help\t\tDetailed help message for any command\n")
}

// TODO: finish implementing view function.
// it works, but we need to figure out an easy way to run the command.
func viewUsage() {
	fmt.Printf(
		"  view\t\tVisualize graph generated alignment\n")
}
func flagsPrint() {
	fmt.Print(
		"\nFlags:\n" +
			"  -h, --help\t\tEnter gsw --help [align/ggtools/view] for detailed information\n" +
			"  -o, --out\t\tFilename[.gg/.vcf/.gz/.sam]  (default: /dev/stdout)\n\n")
	//"  -t, --threads\t\tNumber of CPUs for goroutines  (default: 4)\n\n")
}

func errorMessage() {
	log.Fatalf("Error: Apologies, your command prompt was not recognized...\n\n-xoxo GG\n")
}

func moreHelp(cmdFlag string) {
	if strings.Contains("align", cmdFlag) {
		extendedAlignUsage()
	} else if strings.Contains("ggtools", cmdFlag) {
		ggtoolsExtend()
	} else if strings.Contains("view", cmdFlag) {
		viewExtend()
	} else {
		errorMessage()
	}
}

func viewExtend() {
	fmt.Print(
		"Usage:\n" +
			"  gsw view [options] ref\n\n" +
			"Options:\n" +
			"  -g, --giraf\t\tProvide a giraf alignment with its reference to visualize the alignment\n" +
			"  -s, --sam\t\tProvide a sam alignment with its reference to visualize the alignment\n\n")
}

func ViewArgs() *ViewExe {
	view := &ViewExe{Cmd: flag.NewFlagSet("view", flag.ExitOnError)}
	view.Cmd.StringVar(&view.GirafFile, "giraf", "", "Visualize sequences aligned to graphs in giraf format\n")
	view.Cmd.StringVar(&view.SamFile, "sam", "", "Visualize sequences aligned to graphs in sam format\n\n")
	return view
}

func RunViewExe() error {
	view := ViewArgs()
	view.Cmd.Parse(os.Args[2:])

	tail := view.Cmd.Args()
	if len(tail) > 1 || !strings.HasSuffix(tail[0], ".gg") {
		flag.PrintDefaults()
		return fmt.Errorf("Error: Apologies, your command prompt was not recognized...\n\n-xoxo GG\n")
	}
	viewAlignmentStdOut(tail[0], view.GirafFile)
	return nil
}

func viewAlignmentStdOut(ref string, g string) {
	if !strings.HasSuffix(ref, ".gg") {
		log.Fatalf("Error: Apologies, your command prompt was not recognized...\n\n-xoxo GG\n")
	}
	gg := genomeGraph.Read(ref)
	if strings.HasSuffix(g, ".giraf") {
		file := fileio.EasyOpen(g)
		defer file.Close()
		for curr, done := giraf.NextGiraf(file); !done; curr, done = giraf.NextGiraf(file) {
			log.Printf("%s\n", genomeGraph.ViewGraphAlignment(curr, gg))
		}
	}
}

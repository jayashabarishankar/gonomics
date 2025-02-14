package main

import (
	"github.com/vertgenlab/gonomics/exception"
	"github.com/vertgenlab/gonomics/fileio"
	"os"
	"testing"
)

var SubsequenceSwapTestCases = []struct {
	InFile       string
	Background   int
	Foreground   int
	SwapRegions  string
	OutFile      string
	OutSeqName   string
	ExpectedFile string
}{
	{
		InFile:       "testdata/test.fa",
		Background:   3,
		Foreground:   2,
		SwapRegions:  "testdata/swapRegionsTest.bed",
		OutFile:      "testdata/outputFile.fa",
		OutSeqName:   "Seq5",
		ExpectedFile: "testdata/expectedSwap.fa",
	},
}

func TestSubsequenceSwap(t *testing.T) {
	var err error
	for _, tc := range SubsequenceSwapTestCases {
		//loop through all test cases
		SubsequenceSwap(tc.InFile, tc.SwapRegions, tc.Background, tc.Foreground, tc.OutFile, tc.OutSeqName)

		if !fileio.AreEqual(tc.OutFile, tc.ExpectedFile) {
			t.Errorf("Error: Output was not as expected. \n")
		} else {
			err = os.Remove(tc.OutFile)
			exception.PanicOnErr(err)
		}
	}
}

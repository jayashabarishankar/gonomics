package main

import (
	"log"
	"os"
	"testing"
)

var (
	genome = "testdata/testGenome.fa"
)

func TestBinGenome(t *testing.T) {
	faBin(genome, "testdata", 0, 6, "")

	err1 := os.Remove("testdata/chr1.fa")
	err2 := os.Remove("testdata/testGenome.bin1.fa")

	if err1 != nil || err2 != nil {
		log.Panic(err1, err2)
	}

	faBin(genome, "testdata", 2, -1, "")

	err4 := os.Remove("testdata/testGenome.bin0.fa")
	err5 := os.Remove("testdata/testGenome.bin1.fa")

	if err4 != nil || err5 != nil {
		log.Panic(err4, err5)
	}
}

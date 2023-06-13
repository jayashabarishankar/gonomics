package bedpe

import (
	"github.com/vertgenlab/gonomics/bed"
	"github.com/vertgenlab/gonomics/chromInfo"
	"github.com/vertgenlab/gonomics/interval"
	"log"
	//TODO: Reintroduce when bug is fixed "log"
	//"fmt"
)

// Fill3dSpace takes in a set of BedPe format contact points, single-base pair []bed.Bed representing
// trancription start sites, and a reference genome as a map[string]chromInfo.ChromInfo. This program
// returns a []bed.Bed representing the closest TSS in 3D space for each position in the reference genome,
// using the BedPe contacts to determine the closest TSS in 3D space.
// Note that the input TSS must have 0 in score fields.
func Fill3dSpace(contacts []BedPe, tss []bed.Bed, sizes map[string]chromInfo.ChromInfo) []bed.Bed {
	var chromInFile bool
	geneChroms := makeChromSlice(tss) //makes a list of the chroms in the gene file
	answer := make([]bed.Bed, len(tss))
	copy(answer, tss)
	var currNearestBed, currAnswerA, currAnswerB bed.Bed
	var closest2dGenesIntervals, currNearest []interval.Interval
	closest2dGene := bed.FillSpaceNoHiddenValue(tss, sizes)

	for i := range closest2dGene {
		closest2dGenesIntervals = append(closest2dGenesIntervals, closest2dGene[i])
	}
	closest2dGeneTree := interval.BuildTree(closest2dGenesIntervals)

	midpointBedpe := contactsToMidpoints(contacts)

	// this for loop finds the nearest gene and hidden value for each bedpe foot midpoint
	for j := range midpointBedpe {
		//this is just a check to make sure that for any new chromosomes we encounter that they exist in the gene file, so we don't throw an error.
		if midpointBedpe[j].A.Chrom != midpointBedpe[j-1].A.Chrom {
			chromInFile = checkGeneFileForChrom(midpointBedpe[j], geneChroms)
			if !chromInFile {
				continue
			}
		}
		currNearest = interval.Query(closest2dGeneTree, midpointBedpe[j].A, "any")
		if len(currNearest) > 1 || len(currNearest) == 0 {
			log.Fatal("Space Filled bed should return one nearest bed entry.")
		}
		currNearestBed = currNearest[0].(bed.Bed)
		currAnswerA = bed.Bed{Chrom: midpointBedpe[j].A.Chrom,
			ChromStart:        midpointBedpe[j].A.ChromStart,
			ChromEnd:          midpointBedpe[j].A.ChromEnd,
			Name:              currNearestBed.Name,
			Score:             0,
			FieldsInitialized: 5}
		if currNearestBed.Score < midpointBedpe[j].A.ChromStart {
			currAnswerA.Score = midpointBedpe[j].A.ChromStart - currNearestBed.Score
		} else {
			currAnswerA.Score = currNearestBed.Score - midpointBedpe[j].A.ChromStart
		}

		currNearest = interval.Query(closest2dGeneTree, midpointBedpe[j].B, "any")
		if len(currNearest) > 1 {
			log.Fatal("Space Filled bed should only return one nearest bed entry.")
		}
		currNearestBed = currNearest[0].(bed.Bed)
		currAnswerB = bed.Bed{Chrom: midpointBedpe[j].B.Chrom,
			ChromStart:        midpointBedpe[j].B.ChromStart,
			ChromEnd:          midpointBedpe[j].B.ChromEnd,
			Name:              currNearestBed.Name,
			Score:             0,
			FieldsInitialized: 5}
		if currNearestBed.Score < midpointBedpe[j].B.ChromStart { // currNearestBed.Score is the abolute position of the TSS
			currAnswerB.Score = midpointBedpe[j].B.ChromStart - currNearestBed.Score
		} else {
			currAnswerB.Score = currNearestBed.Score - midpointBedpe[j].B.ChromStart
		}

		if currAnswerA.Score < currAnswerB.Score {
			currAnswerB.Score = currAnswerA.Score
			currAnswerB.Name = currAnswerA.Name
		} else if currAnswerB.Score < currAnswerA.Score {
			currAnswerA.Score = currAnswerB.Score
			currAnswerA.Name = currAnswerB.Name
		} //if scores are equal to each other they both retain their original closest gene assignment
		answer = append(answer, currAnswerA)
		answer = append(answer, currAnswerB)
	}

	return bed.FillSpaceHiddenValue(answer, sizes)
}

func makeChromSlice(records []bed.Bed) []string {
	var chroms []string
	var inChroms bool

	for r := range records {
		if r == 0 {
			chroms = append(chroms, records[r].Chrom)
		} else {
			inChroms = false
			for n := range chroms {
				if records[r].Chrom == chroms[n] {
					inChroms = true
				}
			}
			if !inChroms {
				chroms = append(chroms, records[r].Chrom)
			}
		}
	}
	return chroms
}

func checkGeneFileForChrom(a BedPe, b []string) bool {
	chromInFile := false

	for r := range b {
		if b[r] == a.A.Chrom {
			chromInFile = true
		}
	}
	return chromInFile
}

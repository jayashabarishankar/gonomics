package sam

import (
	"math"
	"testing"
)

func TestSeekBamRegion(t *testing.T) {
	br, _ := OpenBam("testdata/rand.bam")
	bai := ReadBai("testdata/rand.bam.bai")
	reads := SeekBamRegion(br, bai, "chr7", 45000000, 45200000)
	for i := range reads {
		if reads[i].RName != "chr7" || !(reads[i].GetChromStart() < 45200000 && reads[i].GetChromEnd() > 45000000) {
			t.Error("problem with SeekBamRegion")
		}
	}

	reads = SeekBamRegion(br, bai, "chr9", 130590067, 130591448)
	if len(reads) > 0 {
		t.Error("problem with SeekBamRegion")
	}

	reads = SeekBamRegion(br, bai, "chr9", 130591894, 130592016)
	if len(reads) != 1 {
		t.Error("problem with SeekBamRegion")
	}

	reads = SeekBamRegion(br, bai, "chr9", 130592026, 130592027)
	if len(reads) != 2 {
		t.Error("problem with SeekBamRegion")
	}

	reads = SeekBamRegion(br, bai, "chr9", 0, math.MaxUint32)
	if len(reads) != 12 {
		t.Error("problem with SeekBamRegion")
	}

	reads = SeekBamRegion(br, bai, "chrX", 0, 0)
	if len(reads) > 0 {
		t.Error("problem with SeekBamRegion")
	}
}

func TestSeekBamRegionRecycle(t *testing.T) {
	br, _ := OpenBam("testdata/rand.bam")
	bai := ReadBai("testdata/rand.bam.bai")
	reads := make([]Sam, 10)
	reads = SeekBamRegionRecycle(br, bai, "chr7", 45000000, 45200000, reads)
	for i := range reads {
		if reads[i].RName != "chr7" || !(reads[i].GetChromStart() < 45200000 && reads[i].GetChromEnd() > 45000000) {
			t.Error("problem with SeekBamRegionRecycle")
		}
	}

	reads = SeekBamRegionRecycle(br, bai, "chr9", 130590067, 130591448, reads)
	if len(reads) > 0 {
		t.Error("problem with SeekBamRegionRecycle")
	}

	reads = SeekBamRegionRecycle(br, bai, "chr9", 130591894, 130592016, reads)
	if len(reads) != 1 {
		t.Error("problem with SeekBamRegionRecycle")
	}

	reads = SeekBamRegionRecycle(br, bai, "chr9", 130592026, 130592027, reads)
	if len(reads) != 2 {
		t.Error("problem with SeekBamRegionRecycle")
	}
	for i := range reads {
		if reads[i].RName != "chr9" || !(reads[i].GetChromStart() < 130592026 && reads[i].GetChromEnd() > 130592027) {
			t.Error("problem with SeekBamRegionRecycle")
		}
	}

	reads = SeekBamRegionRecycle(br, bai, "chr9", 0, math.MaxUint32, reads)
	if len(reads) != 12 {
		t.Error("problem with SeekBamRegionRecycle")
	}

	reads = SeekBamRegionRecycle(br, bai, "chrX", 0, 0, reads)
	if len(reads) > 0 {
		t.Error("problem with SeekBamRegionRecycle")
	}
}

func TestSeekManyReads(t *testing.T) {
	br, _ := OpenBam("testdata/peak.bam")
	bai := ReadBai("testdata/peak.bam.bai")
	reads := SeekBamRegion(br, bai, "chr9", 130591103, 130592987)
	if len(reads) != 561 {
		t.Error("problem with SeekBamRegion")
	}
}

package tworker

import (
	"fmt"
	"regexp"
	"testing"
)

func TestWorkerTb(t *testing.T) {
	expr := `1dcj/file_hwrsr" type="video/mp4" awgr47ar7e8 1dcj/file_ghwurawn" type="video/mp4"`
	reg, err := regexp.Compile(`.*?file_(.*?)" type="video/mp4"`)

	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	result := reg.FindAllStringSubmatch(expr, -1)

	fmt.Printf("%d : %q\n", len(result), result)
}

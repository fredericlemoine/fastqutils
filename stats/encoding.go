package stats

import (
	"errors"
	"fmt"
	"github.com/fredericlemoine/fastqutils/error"
)

const (
	SANGER = iota
	SOLEXA
	ILLUMINA_1_3
	ILLUMINA_1_5
	ILLUMINA_1_8
	UNKOWN
)

func DetectEncoding(min, max int) int {
	if min >= 33 && max <= 73 {
		return SANGER
	} else if min >= 33 && max <= 74 {
		return ILLUMINA_1_8
	} else if min >= 67 && max <= 104 {
		return ILLUMINA_1_5
	} else if min >= 64 && max <= 104 {
		return ILLUMINA_1_3
	} else if min >= 59 && max <= 104 {
		return SOLEXA
	} else {
		return UNKOWN
	}
}

func EncodingToString(encod int) string {
	switch encod {
	case SANGER:
		return "Sanger"
	case SOLEXA:
		return "Solexa"
	case ILLUMINA_1_3:
		return "Illumina 1.3"
	case ILLUMINA_1_5:
		return "Illumina 1.5"
	case ILLUMINA_1_8:
		return "Illumina 1.8"
	case UNKOWN:
		return "Unknown"
	default:
		error.ExitWithMessage(errors.New(fmt.Sprintf("This encoding Code does not exist : %d", encod)))
	}
	return ""
}

func EncodingFromString(encod string) int {
	switch encod {
	case "sanger":
		return SANGER
	case "solexa":
		return SOLEXA
	case "illumina1.3":
		return ILLUMINA_1_3
	case "illumina1.5":
		return ILLUMINA_1_5
	case "illumina1.8":
		return ILLUMINA_1_8
	case "unknown":
		return UNKOWN
	default:
		error.ExitWithMessage(errors.New(fmt.Sprintf("This encoding Code does not exist : %s", encod)))
	}
	return UNKOWN
}

func EncodingOffset(encod int) int {
	switch encod {
	case SANGER:
		return 33
	case SOLEXA:
		return 64
	case ILLUMINA_1_3:
		return 64
	case ILLUMINA_1_5:
		return 64
	case ILLUMINA_1_8:
		return 33
	case UNKOWN:
		return 0
	default:
		error.ExitWithMessage(errors.New(fmt.Sprintf("This encoding Code does not exist : %d", encod)))
	}
	return 0
}

func MinQual(encod int) int {
	switch encod {
	case SANGER:
		return 33
	case SOLEXA:
		return 59
	case ILLUMINA_1_3:
		return 64
	case ILLUMINA_1_5:
		return 67
	case ILLUMINA_1_8:
		return 33
	case UNKOWN:
		return 0
	default:
		error.ExitWithMessage(errors.New(fmt.Sprintf("This encoding Code does not exist : %d", encod)))
	}
	return 0
}

func MaxQual(encod int) int {
	switch encod {
	case SANGER:
		return 73
	case SOLEXA:
		return 104
	case ILLUMINA_1_3:
		return 104
	case ILLUMINA_1_5:
		return 104
	case ILLUMINA_1_8:
		return 74
	case UNKOWN:
		return 126
	default:
		error.ExitWithMessage(errors.New(fmt.Sprintf("This encoding Code does not exist : %d", encod)))
	}
	return 0
}

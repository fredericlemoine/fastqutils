package stats

import (
	"fmt"
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

func EncodingToString(encod int) (enc string, err error) {
	switch encod {
	case SANGER:
		enc = "Sanger"
	case SOLEXA:
		enc = "Solexa"
	case ILLUMINA_1_3:
		enc = "Illumina 1.3"
	case ILLUMINA_1_5:
		enc = "Illumina 1.5"
	case ILLUMINA_1_8:
		enc = "Illumina 1.8"
	case UNKOWN:
		enc = "Unknown"
	default:
		err = fmt.Errorf("this encoding Code does not exist : %d", encod)
	}
	return
}

func EncodingFromString(encod string) (enc int, err error) {
	switch encod {
	case "sanger":
		enc = SANGER
	case "solexa":
		enc = SOLEXA
	case "illumina1.3":
		enc = ILLUMINA_1_3
	case "illumina1.5":
		enc = ILLUMINA_1_5
	case "illumina1.8":
		enc = ILLUMINA_1_8
	case "unknown":
		enc = UNKOWN
	default:
		err = fmt.Errorf("this encoding Code does not exist : %s, possible values are : sanger, solexa, illumina1.3, illumina1.5, illumina1.8", encod)
	}
	return
}

func EncodingOffset(encod int) (off int, err error) {
	switch encod {
	case SANGER:
		off = 33
	case SOLEXA:
		off = 64
	case ILLUMINA_1_3:
		off = 64
	case ILLUMINA_1_5:
		off = 64
	case ILLUMINA_1_8:
		off = 33
	case UNKOWN:
		off = 0
	default:
		err = fmt.Errorf(" : %d", encod)
	}
	return
}

func MinQual(encod int) (minq int, err error) {
	switch encod {
	case SANGER:
		minq = 33
	case SOLEXA:
		minq = 59
	case ILLUMINA_1_3:
		minq = 64
	case ILLUMINA_1_5:
		minq = 67
	case ILLUMINA_1_8:
		minq = 33
	case UNKOWN:
		minq = 0
	default:
		err = fmt.Errorf("this encoding Code does not exist : %d", encod)
	}
	return
}

func MaxQual(encod int) (maxq int, err error) {
	switch encod {
	case SANGER:
		maxq = 73
	case SOLEXA:
		maxq = 104
	case ILLUMINA_1_3:
		maxq = 104
	case ILLUMINA_1_5:
		maxq = 104
	case ILLUMINA_1_8:
		maxq = 74
	case UNKOWN:
		maxq = 126
	default:
		err = fmt.Errorf("this encoding Code does not exist : %d", encod)
	}
	return
}

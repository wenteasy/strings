package strings

import (
	"fmt"
	"strconv"
	"strings"
)

//Binary Byte
type BByte uint64

type BByteUnit int

// EiB            1,152,921,504,606,846,976 B
// uint64        18,446,744,073,709,551,615
// ZiB        1,180,591,620,717,411,303,424 B
// YiB    1,208,925,819,614,629,174,706,176 B

const (
	BByteB BByteUnit = iota
	BByteK
	BByteM
	BByteG
	BByteT
	BByteP
	BByteE
	BByteZ
	BByteY
	BByteUndefined
	BByteError
)

var bbyteUnitName = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

func (u BByteUnit) String() string {
	if u == BByteError {
		return "Error"
	} else if u == BByteUndefined {
		return "YiB Over"
	}
	return bbyteUnitName[u]
}

func (u BByteUnit) KMGT() string {
	if u == BByteB {
		return ""
	} else if u == BByteError {
		return "Error"
	} else if u == BByteUndefined {
		return "YiB Over"
	}
	return string(bbyteUnitName[u][0])
}

const unit1K = 1 << 10

// Value() is
func (v BByte) Value() (float64, BByteUnit) {
	val := float64(v)
	u := BByteB
	for ; val >= unit1K; u++ {
		val /= unit1K
	}
	return val, u
}

func (v BByte) Format(short bool, fms ...string) string {

	fm := "%.02f"
	if len(fms) > 0 {
		fm = fms[0]
	}

	val, u := v.Value()

	if u == BByteB {
		fm = "%.0f"
		if len(fms) > 1 {
			fm = fms[1]
		}
	}

	if short {
		return fmt.Sprintf(fm+"%s", val, u.KMGT())
	}
	return fmt.Sprintf(fm+"%s", val, u)
}

func (v BByte) Spec(bu BByteUnit) float64 {
	val := float64(v)
	for u := BByteB; u < bu; u++ {
		val /= unit1K
	}
	return val
}

func (v BByte) String() string {
	return v.Format(false)
}

func ParseBByte(v string) BByte {

	t := v
	u := ""
	fmt.Println(v)

	for idx, b := range []byte(t) {
		if b >= '0' && b <= '9' || b == '.' {
			continue
		}
		t = v[:idx]
		u = v[idx:]
		break
	}

	iv, err := strconv.ParseFloat(t, 64)
	if err == nil && u != "" {

		tl := strings.ToLower(u)
		tu := ""
		if len(tl) >= 1 {
			tu = tl[:1]
		}

		for idx, elm := range bbyteUnitName {

			l := strings.ToLower(elm)

			if tl == l || tu == l[:1] {
				for jdx := 0; jdx < idx; jdx++ {
					iv = iv * unit1K
				}
			}
		}
	}

	return BByte(uint64(iv))
}

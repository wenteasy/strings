package strings_test

import (
	"fmt"
	"testing"

	"github.com/wenteasy/strings"
)

func TestBByte_Value(t *testing.T) {
}

func TestBByte_Format(t *testing.T) {
}

func TestBByte_Spec(t *testing.T) {
}

func TestBByte_Parse(t *testing.T) {

	wants := []struct {
		Arg string
		Ans strings.BByte
	}{
		{"0", strings.BByte(0)},
		{"3024", strings.BByte(3_024)},
		{"0.5M", strings.BByte(524_288)},
		{"1.2k", strings.BByte(1_228)},
		{"1kibibyte", strings.BByte(1_024)},
		{"200.0m", strings.BByte(209_715_200)},
		{"45024.", strings.BByte(45_024)},
		{"12..21", strings.BByte(0)}, //ParseError
	}

	for _, want := range wants {
		got := strings.ParseBByte(want.Arg)
		if got != want.Ans {
			t.Errorf("BByte.Parse(%s): want[%v], got[%v]", want.Arg, uint64(want.Ans), uint64(got))
		}
	}
}

func ExampleBByteUnit() {

	fmt.Println(strings.BByteB)
	fmt.Println(strings.BByteB.KMGT())
	fmt.Println(strings.BByteK)
	fmt.Println(strings.BByteK.KMGT())
	fmt.Println(strings.BByteM)
	fmt.Println(strings.BByteM.KMGT())
	fmt.Println(strings.BByteG)
	fmt.Println(strings.BByteG.KMGT())
	fmt.Println(strings.BByteT)
	fmt.Println(strings.BByteT.KMGT())
	fmt.Println(strings.BByteP)
	fmt.Println(strings.BByteP.KMGT())
	fmt.Println(strings.BByteE)
	fmt.Println(strings.BByteE.KMGT())
	fmt.Println(strings.BByteZ)
	fmt.Println(strings.BByteZ.KMGT())
	fmt.Println(strings.BByteY)
	fmt.Println(strings.BByteY.KMGT())
	fmt.Println(strings.BByteUndefined)
	fmt.Println(strings.BByteUndefined.KMGT())
	fmt.Println(strings.BByteError)
	fmt.Println(strings.BByteError.KMGT())

	// Output:
	// B
	//
	// KiB
	// K
	// MiB
	// M
	// GiB
	// G
	// TiB
	// T
	// PiB
	// P
	// EiB
	// E
	// ZiB
	// Z
	// YiB
	// Y
	// YiB Over
	// YiB Over
	// Error
	// Error
	//
}
func ExampleBByte() {

	b := strings.BByte(1024)
	fmt.Println(b)
	v, u := b.Value()
	fmt.Printf("%0.0f %s\n", v, u)

	b = strings.BByte(1024 * 1024)
	fmt.Println(b)
	v, u = b.Value()
	fmt.Printf("%0.0f %s\n", v, u)

	b -= 1
	fmt.Println(b)
	v, u = b.Value()
	fmt.Printf("%0.5f %s\n", v, u)
	fmt.Printf("%0.9f M\n", b.Spec(strings.BByteM))

	// Output:
	// 1.00KiB
	// 1 KiB
	// 1.00MiB
	// 1 MiB
	// 1024.00KiB
	// 1023.99902 KiB
	// 0.999999046 M
	//
}

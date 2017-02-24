// Verify how the analysis handles aliasing through fields
// X and Y different types
//        x----------->{U.u:T.s.......}
//                       |
//                       |
//                       V
//                      {taint}
//                       ^
//                       |
//                       |
//      y------------>{V.u:T.s-------]
package main

type T struct {
	s string
}

type U struct {
	u T
}

type V struct {
	v T
}

func main() {
	t := T{s: source()}
	x := new(U)
	y := new(V)

	x.u = t
	y.v = t

	// @ExpectFlow: true
	sink(x.u.s)
	// @ExpectFlow: true
	sink(y.v.s)
}

func sink(s string) {
}

func source() string {
	return "secret"
}

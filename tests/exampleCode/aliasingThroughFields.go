// Verify how the analysis handles aliasing through fields
// X and Y of the same type (U)
//        x----------->{U.u:T.s.......}
//                       |
//                       |
//                       V
//                      {taint}
//                       ^
//                       |
//                       |
//      y------------>{U.u:T.s-------]
package main

type T struct {
	s string
}

type U struct {
	u T
}

func main() {
	t := T{s: source()}
	x := new(U)
	y := new(U)

	x.u = t
	y.u = t

	// @ExpectFlow: true
	sink(x.u.s)
	// @ExpectFlow: true
	sink(y.u.s)
}

func sink(s string) {
}

func source() string {
	return "secret"
}

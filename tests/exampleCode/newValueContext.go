package main

import (
	//"fmt"
	"os"
)

func main() {
	f, _ := os.OpenFile("./dummy.txt", os.O_RDWR, os.ModeAppend)
	//f, err := os.OpenFile("./dummy.txt", os.O_RDWR, os.ModeAppend)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	var b []byte
	b = make([]byte, 5, 5)
	_, _ = f.Read(b)
	//_, err = f.Read(b)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	funF(b, true)
}

func funF(a []byte, c bool) {
	if c {
		f2, _ := os.OpenFile("./dummy2.txt", os.O_RDWR, os.ModeAppend)
		//		f2, err := os.OpenFile("../examplecode/dummy2.txt", os.O_RDWR, os.ModeAppend)
		//	if err != nil {
		//		fmt.Println(err.Error())
		//	}
		_, _ = f2.Write(a)
		//	_, err = f3.Write(a)
		//	if err != nil {
		//		fmt.Println(err.Error())
		//	}
	} else {
		return
	}
}

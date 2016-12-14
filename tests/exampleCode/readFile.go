package main

import "os"

func main() {
	//	f, _ := os.OpenFile("../examplecode/dummy.txt", os.O_RDWR, os.ModeAppend)
	f, err := os.OpenFile("./dummy.txt", os.O_RDWR, os.ModeAppend)
	if err != nil {
		return
	}
	var b []byte
	b = make([]byte, 5, 5)
	_, err = f.Read(b)
	if err != nil {
		return
	}
	var f2 *os.File
	f2, err = os.OpenFile("./dummy2.txt", os.O_RDWR, os.ModeAppend)
	if err != nil {
		return
	}
	_, err = f2.Write(b)
	if err != nil {
		return
	}
}

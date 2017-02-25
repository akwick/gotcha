package taint

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

// Data is holds the signature and the callee of a source or sink.
type Data struct {
	*taintData
}

type taintData struct {
	sig         string
	callee      string
	isInterface bool
}

// IsInterface returns true if t contains the signature of an interface and
// not of an concrete type
func (t *Data) IsInterface() bool {
	return t.isInterface
}

func (t *Data) String() string {
	return t.sig + " " + t.callee
}

func (t *Data) GetSig() string {
	return t.sig
}

// Sinks contains a slice of TaintData which are sinks.
var Sinks = make([]*Data, 0)

// Sources contains a slice of TaintData which are sources.
var Sources = make([]*Data, 0)

// Read the file with the sources and sinks
func Read(fileName string) error {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Print(err)
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		currentLine := scanner.Text()
		// check whether its an interface (starts with I)
		isInterface := false
		if strings.HasPrefix(currentLine, "I") {
			isInterface = true
			currentLine = strings.TrimPrefix(currentLine, "I")
		}

		// source or sinks
		if strings.HasPrefix(currentLine, "<") {
			currentLine = strings.TrimPrefix(currentLine, "<")
			fields := strings.Fields(currentLine)
			if len(fields) < 5 {
				return errors.New("The string " + currentLine + " is too short for the parser (expect at minimum 5 fields)")
			}
			var td *Data
			if fields[len(fields)-1] == "_SOURCE_" {
				td, err = getTaintData(fields)
				if err != nil {
					return err
				}
				td.isInterface = isInterface
				Sources = append(Sources, td)
			} else {
				if fields[len(fields)-1] == "_SINK_" {
					td, err = getTaintData(fields)
					if err != nil {
						return err
					}
					td.isInterface = isInterface
					Sinks = append(Sinks, td)
				} else {
					return errors.New("Unknown ending of " + currentLine + ". Expected _SOURCE_ or _SINK")
				}
			}
		} else {
			// comment
			if strings.HasPrefix(currentLine, "%") {
				// do nothing it's only a comment
			} else {
				return errors.New("Invalid line start " + currentLine)
			}
		}
	}
	// Returning an error if either no source or sink is defined
	if len(Sources) == 0 {
		return errors.New("the provided file does not contain a valid source signature")
	}
	if len(Sinks) == 0 {
		return errors.New("the provided file does not contain a valid sink signature")
	}
	return nil
}

func getTaintData(s []string) (*Data, error) {
	var sig string
	if strings.Contains(s[0], "func(") {
	} else {
		return nil, errors.New("Expected a string starting with func(, but get: " + s[0])
	}
	//	s[0] = strings.TrimPrefix(s[0], "func(")

	// Creating the signature with the form
	// func(something) (return value);
	retValI := 1
	// Handling func(something)
	for i, strng := range s {
		//		fmt.Print(strng + " ")
		sig = sig + " " + strng
		if strings.HasSuffix(strng, ")") {
			retValI = i
			break
		}
	}

	// Handling the return values
	for i := retValI + 1; i < len(s); i++ {
		sig = sig + " " + s[i]
		//		fmt.Print(" retVal: " + s[i] + " ")
		if strings.HasSuffix(s[i], ";") {
			retValI = i
			break
		}
	}
	sig = strings.TrimSuffix(sig, ";")
	sig = strings.TrimPrefix(sig, " ")
	sig = strings.TrimSuffix(sig, " ")

	// Handling the Callee of the function
	var callee string
	for i := retValI + 1; i < len(s); i++ {
		callee = callee + s[i]
		if strings.HasSuffix(s[i], ">") {
			retValI = i
			break
		}
	}
	callee = strings.TrimSuffix(callee, ">")
	//	fmt.Println(" callee: " + callee)

	td := &Data{&taintData{sig: sig, callee: callee}}
	if s[len(s)-2] != "->" {
		return nil, errors.New("Expected the string -> as part of the line: " + s[3])
	}
	if len(s)-retValI != 3 {
		return nil, errors.New("The string is wrong formatted")
	}
	return td, nil
}

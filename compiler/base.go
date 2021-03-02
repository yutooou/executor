package compiler

var compileCommands = struct {
	GNUC   string
	GNUCPP string
	Java   string
}{
	GNUC:   "gcc %s -o %s -std=c11",
	GNUCPP: "g++ %s -o %s -std=c++11",
}

package main

import "fmt"

func main() {
	NewBenString("4:spam")
	NewBenString("0:")
	NewBenString("4:spams")
	NewBenString("0gg")
	NewBenString("0:gg")

	fmt.Println(NewBenStringFromValue(""))
	fmt.Println(NewBenStringFromValue("spam"))
}

package main

import (
	"fmt"

	"github.com/Qubitopia/QuantumScholar/server/mail"
)

func main() {
	fmt.Println("Hello, World!")

	var intArray [3]int32 = [3]int32{1, 2, 3}
	fmt.Println(&intArray[0])
	fmt.Println(&intArray[1])
	fmt.Println(&intArray[2])

	mail.SendEmailTo("chetaningale@acpce.ac.in", "Chetan Ingale", "https://quantumscholar.pages.dev/login?token=123456789")
}

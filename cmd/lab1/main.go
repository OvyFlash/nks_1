package main

import (
	"fmt"
	"nks/configs"
	"nks/pkg/lab1"
)

func main() {
	c := configs.NewConfig()
	l := lab1.NewLab1(c.Lab1Config)
	l.Start()
	//Print result
	fmt.Println(l)
}

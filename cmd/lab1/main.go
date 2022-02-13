package main

import (
	"fmt"
	"lab1/configs"
	"lab1/pkg/lab1"
)

func main() {
	c := configs.NewConfig()
	l := lab1.NewLab1(c)
	l.Start()
	//Print result
	fmt.Println(l)
}

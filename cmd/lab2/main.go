package main

import (
	"fmt"
	"nks/configs"
	"nks/pkg/lab2"
)

func main() {
	c := configs.NewConfig()
	l := lab2.NewLab2(c.Lab2Config)
	l.Start()
	//Print result
	fmt.Println(l)
}

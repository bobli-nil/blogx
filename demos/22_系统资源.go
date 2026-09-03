package main

import (
	"blogx_server/utils/computer"
	"fmt"
)

func main() {
	r1 := computer.GetCpuPercent()
	fmt.Println("CPU使用率", r1)

	r2 := computer.GetMemPercent()
	fmt.Println("内存使用率", r2)

	r3 := computer.GetDiskPercent()
	fmt.Println("磁盘使用率", r3)

}

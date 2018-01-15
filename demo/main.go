package main

import (
	"stock/source/tencent/marketCenter"
	"fmt"
)

func main() {
	class , err := marketCenter.ReaderClass()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _ , v := range class {
		fmt.Println(v.Name,"@",v.Id)
		for _ , val := range v.SubClass{
			fmt.Println("\t",val.Name,"@@",val.Id)
			for key , value := range val.LastSubClass{
				fmt.Println("\t\t",value,"*",key)
			}
		}
	}
}
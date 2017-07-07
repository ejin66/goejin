package main

import (
	"fmt"
)
type A struct {
	i int
	*B
}
type B struct {
	j int
}
func main() {
	//result := db.Query("user_info", nil)
	//fmt.Println(result)
	//ok := db.Delete("user_info",map[string]interface{}{ "id" : 1 }," or id = 4", " or id = 5", " or id = 6")
	//fmt.Println(ok)
	//ok := db.Update("user_info",map[string]interface{}{"user_name" : "ejin666"},map[string]interface{}{"id" : 2})
	//fmt.Println(ok)
	//ok := db.Insert("user_info",db.Ipt{"user_name" : "hahaha" })
	//fmt.Println(ok)


	a := A{B:&B{}}
	var b A
	c := new(A)

	fmt.Printf("a:%p\n" , &a )
	fmt.Printf("b:%p\n" , &b )
	fmt.Printf("c:%p\n" , c )

	fmt.Println(a,a.i,a.B.j)
	fmt.Println(b,b.i,b.B.j)
	fmt.Println(c,c.i,c.B.j)
}

func (this *A) do() {
	fmt.Println("sssssssssssssssssssss")
}
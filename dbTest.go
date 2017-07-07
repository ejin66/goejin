package main

import "fmt"

func main() {
	//result := db.Query("user_info", nil)
	//fmt.Println(result)
	//ok := db.Delete("user_info",map[string]interface{}{ "id" : 1 }," or id = 4", " or id = 5", " or id = 6")
	//fmt.Println(ok)
	//ok := db.Update("user_info",map[string]interface{}{"user_name" : "ejin666"},map[string]interface{}{"id" : 2})
	//fmt.Println(ok)
	//ok := db.Insert("user_info",db.Ipt{"user_name" : "hahaha" })
	//fmt.Println(ok)

	var i = 0

	for {

		if i == 10 {
			break
		}

		i++
		fmt.Println("i:",i)
	}

}
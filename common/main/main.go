package main

import (
	"fmt"

	"../mysql"
)

type User struct {
	Id     int
	Name   string
	Sex    string
	Age    int
	IdCard string
}

func main() {
	mysql.MySqlDb = mysql.NewMySqlDb("user", "passwd", "ip", "db", 23)
	mysql.MySqlDb.DbConn()
	defer mysql.MySqlDb.Close()

	user := &User{Name: "hanyan", Sex: "man", Age: 27, IdCard: "12345"}
	mysql.MySqlDb.Insert("user", user)

	//condsMap := make(map[string]interface{})
	//condsMap["id_card"] = "12345"
	//condsMap["id"] = 1
	//condsMap := map[string]interface{}{"name": "hanyan", "id in": []int{4, 6, 1}, "age !=": 27}
	condsMap := map[string]interface{}{"id": 22}
	//var users []User
	//sqlMap := map[string]interface{}{"age": 80, "name": "suannaithllll"}
	//sqlMap := &User{Name: "hy thl suannai", Sex: "man", Age: 99999, IdCard: "12345xxxxxxx"}
	//mysql.MySqlDb.UpdateByConds("user", sqlMap, condsMap)
	//fmt.Println(users)
	condsMap1 := map[string]interface{}{"name": "hanyan", "id in": []int{28, 29, 27}, "age !=": 27}
	var users []User
	//mysql.MySqlDb.SelectBySqlString("user", &users, condsMap1)
	fmt.Println(users)
	fmt.Println(condsMap)
	fmt.Println(condsMap1)
	mysql.MySqlDb.DeleteByConds("user", users, condsMap)
	mysql.MySqlDb.DeleteBySqlString("user", users, condsMap1)

}

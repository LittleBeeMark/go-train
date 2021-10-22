package test_json

import (
	"encoding/json"
	"fmt"
	"testing"

	"go-train/gorm_test/dao"
)

func init() {
	dao.Init()
}

func TestTestJson_Insert(t *testing.T) {
	mm := []*People{
		//&People{
		//	Name: "mark",
		//	Age: 18,
		//	Sex: "男",
		//},
	}
	mmRaw, err := json.Marshal(mm)
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	fmt.Println(string(mmRaw))

	data := TestJson2{}
	//err = data.MarshalPeople()
	//if err != nil {
	//	fmt.Println("err : ", err)
	//	return
	//}

	err = data.Insert()
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
}

func TestTestJson_Update(t *testing.T) {
	mm := []*People{
		&People{
			Name: "mark",
			Age:  18,
			Sex:  "男",
		},
	}
	mmRaw, err := json.Marshal(mm)
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	fmt.Println(mmRaw)

	data := &TestJson2{
		//ID:     "OhLN4Q0K",
		//People: mmRaw,
	}

	err = data.Update()
	fmt.Println("err : ", err)
}


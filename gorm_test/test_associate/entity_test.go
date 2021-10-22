package test_associate

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"

	"go-train/gorm_test/dao"
)

func init() {
	dao.Init()
	Init()
}

func TestUser_Create(t *testing.T) {
	user := User{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "jinzhu",
		//BillingAddress:  Address{Address1: "Billing Address - Address 1", AddressNo: 1},
		ShippingAddress: Address{Address1: "Shipping Address - Address 2", AddressNo: 110},
		//HomeAddress:     Address{Address1: "Shipping Address - Address 3", AddressNo: 3},
		//RAddress:        Address{Address1: "Shipping Address - Address 4", AddressNo: 4},
		Emails: []Email{
			{Email: "jinzhu@example.com"},
			{Email: "jinzhu-2@example.com"},
		},
		Languages: []Language{
			{Name: "EDG"},
			{Name: "FPX"},
		},
	}
	fmt.Println("err : ", user.Create())
}

func TestUser_UpdateEmail(t *testing.T) {
	user := User{
		Model: gorm.Model{ID: 1},
	}
	err := user.UpdateEmail([]Email{
		{Email: "a", UserID: 1},
		{Email: "b", UserID: 1},
	})

	if err != nil {
		t.Log(err)
	}
}

func TestUser_UpdateLanguage(t *testing.T) {
	user := User{
		Model: gorm.Model{ID: 1},
	}
	err := user.UpdateLanguage([]Language{
		{Name: "BM"},
		{Name: "Bf"},
	})

	if err != nil {
		t.Log(err)
	}
}

func TestUser_Query(t *testing.T) {
	user := &User{}
	err := user.Query()
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
	fmt.Printf("%+v ..", user)
}

func TestFriend_FirstOrCreate(t *testing.T) {
	p := People{
		Name: "mark",
		FriendList: []Friend{
			{Name: "rose"},
			{Name: "tidy"},
		},
	}
	fmt.Println("err : ", p.Create())
}

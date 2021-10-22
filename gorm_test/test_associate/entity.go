package test_associate

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"go-train/gorm_test/dao"
)

func Init() {
	dao.GetDB().AutoMigrate(
		&User{},
		&Address{},
		&Email{},
		&Language{},
		&People{},
		&Friend{},
	)

}

// User has and belongs to many languages, `user_languages` is the join table
type User struct {
	gorm.Model
	Name string
	//BillingAddress  Address `gorm:"foreignKey:ID"`
	//HomeAddress     Address
	//RAddress        Address    `gorm:"foreignKey:Address1;references:addr"`
	ShippingAddress Address    `gorm:"foreignKey:AddressNo"`
	Emails          []Email    `gorm:"foreignKey:UserID"`
	Languages       []Language `gorm:"many2many:user_languages;save_associations:false"`
}

type Address struct {
	gorm.Model
	Address1  string
	AddressNo int
}

type Email struct {
	gorm.Model
	Email  string
	UserID uint
}

type Language struct {
	gorm.Model
	Name string `gorm:"unique_index:idx_unique_common_name"`
}

// Create doc
func (u *User) Create() error {
	var err error
	return dao.GetDB().Transaction(func(tx *gorm.DB) error {
		err = dao.GetDB().Debug().Save(u).Error
		if err != nil {
			return err
		}

		return u.UpdateLanguage(u.Languages)
	})
}

// UpdateEmail doc
func (u *User) UpdateEmail(param []Email) error {
	return dao.GetDB().Model(&User{Model: gorm.Model{ID: 1}}).Debug().
		Association("Emails").Replace(param).Error
}

// UpdateLanguage doc
func (u *User) UpdateLanguage(param []Language) error {
	var err error
	return dao.GetDB().Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(param); i++ {
			err = param[i].FirstOrCreate()
			if err != nil {
				return err
			}

			fmt.Printf("language : %+v \n", param[i])
		}

		return dao.GetDB().Model(&User{Model: gorm.Model{ID: u.ID}}).Debug().
			Association("Languages").Replace(param).Error
	})
}

func (l *Language) FirstOrCreate() error {
	return dao.GetDB().Where("name=?", l.Name).FirstOrCreate(l).Error
}

// Query doc
func (u *User) Query() error {
	return dao.GetDB().
		Joins("JOIN user_languages ON user_languages.user_id=users.id").
		Joins("JOIN languages ON user_languages.language_id=languages.id").
		Where("languages.name=?", "S").
		Where("users.id=?", 1).
		Preload("Languages").
		Preload("ShippingAddress").
		Preload("Emails").
		Find(u).Error
}

type People struct {
	gorm.Model
	Name       string
	FriendList Friends `gorm:"many2many:people_friends;save_associations:false"`
}

type Friend struct {
	gorm.Model
	Name string `gorm:"unique_index:idx_unique_friend_name"`
}

func (f *Friend) FirstOrCreate() error {
	return dao.GetDB().Where("name=?", f.Name).FirstOrCreate(f).Error
}

func (f Friend) TableName() string {
	return "people_mark_friends"
}

func (p *People) UpdateFriends(param Friends) error {

	for i := 0; i < len(p.FriendList); i++ {
		err := p.FriendList[i].FirstOrCreate()
		if err != nil {
			return err
		}

		fmt.Printf("friend : %+v \n", p.FriendList[i])

	}

	return dao.GetDB().Model(&People{Model: gorm.Model{ID: p.ID}}).Debug().
		Association("").Replace(param).Error

}

type Friends []Friend

// Create doc
func (p *People) Create() error {

	var err error
	return dao.GetDB().Transaction(func(tx *gorm.DB) error {
		err = dao.GetDB().Debug().Create(p).Error
		if err != nil {
			return err
		}
		return p.UpdateFriends(p.FriendList)
	})
}

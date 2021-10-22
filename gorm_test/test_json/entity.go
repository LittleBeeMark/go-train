package test_json

import (
	"encoding/json"
	"fmt"

	"go-train/gorm_test/dao"
)

func Init() {
	dao.GetDB().AutoMigrate(
		// 缓存域名验证值
		&TestJson{},
		&TestJson2{},
	)
}

// TestJson doc
type TestJson struct {
	ID     string          `gorm:"column:id;PRIMARY_KEY"`
	People json.RawMessage `gorm:"json_message"`

	PeopleInfos Peoples `gorm:"-"`
}

type People struct {
	Name string
	Age  int
	Sex  string
}
type Peoples []*People

func (t *TestJson) TableName() string {
	return "test_jsons"
}

// Insert doc
func (t *TestJson) Insert() error {
	return dao.GetDB().Model(&TestJson{}).Debug().Create(t).Error
}

// Update doc
func (t *TestJson) Update() error {
	return dao.GetDB().Model(&TestJson{}).Debug().Save(t).Error
}

// List doc
func (t *TestJson) List() ([]*TestJson, error) {
	test := make([]*TestJson, 0)
	err := dao.GetDB().Debug().Find(&test).Error
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (t *TestJson) UnMarshalPeople() error {
	return json.Unmarshal(t.People, &t.PeopleInfos)
}

func (t *TestJson) MarshalPeople() error {
	raw, err := json.Marshal(t.People)
	if err != nil {
		fmt.Println("marshal err ", err)
		return err
	}

	t.People = raw
	return nil

}

// TestJson2 doc
type TestJson2 struct {
	ID     string `gorm:"column:id;PRIMARY_KEY"`
	People []byte `gorm:"json_message;type:jsonb;default:null"`

	PeopleInfos Peoples `gorm:"-"`
}

func (t *TestJson2) TableName() string {
	return "test_jsons2"
}

// Insert doc
func (t *TestJson2) Insert() error {
	return dao.GetDB().Model(&TestJson{}).Debug().Create(t).Error
}

// Update doc
func (t *TestJson2) Update() error {
	return dao.GetDB().Model(&TestJson{}).Debug().Save(t).Error
}

// List doc
func (t *TestJson2) List() ([]*TestJson2, error) {
	test := make([]*TestJson2, 0)
	err := dao.GetDB().Debug().Find(&test).Error
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (t *TestJson2) UnMarshalPeople() error {
	return json.Unmarshal(t.People, &t.PeopleInfos)
}

func (t *TestJson2) MarshalPeople() error {
	raw, err := json.Marshal(t.People)
	if err != nil {
		fmt.Println("marshal err ", err)
		return err
	}

	t.People = raw
	return nil

}

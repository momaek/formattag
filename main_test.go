package main

//import "time"

// TestStruct this is a test struct
/*
type TestStruct struct {
	ID            string `json:"id"              xml:"id"              bson:"id"`
	IfNotModified string `json:"if_not_modified" xml:"if_not_modified"`
	Name          string `json:"name"            xml:"name"            bson:"name"`

	ThisIsAStructWodeTianNa struct {
		FieldFromThisIsAStructWodeTianNa string `json:"field_from_this_is_a_struct_wode_tian_na" xml:"field_from_this_is_a_struct_wode_tian_na"`
		TianName                         string `json:"tian_name"                                xml:"tian_name"`
	} `json:"this_is_a_struct_wode_tian_na" xml:"this_is_a_struct_wode_tian_na"`

	T    time.Time `json:"t"    xml:"t"    bson:"t"`
	Fset Fset      `json:"fset" xml:"fset" bson:"fset"`
}

type Fset struct{}
*/

// Location 位置，分部，办公室
type Location struct {
	ID             uint    `json:"id,omitempty"    gorm:"primaryKey"`                                           // 主建
	OrganizationID uint    `json:"organization_id" gorm:"uniqueIndex:uix_org_name"`                             // 组织 ID
	Name           string  `json:"name"            gorm:"size:255;uniqueIndex:uix_org_name" binding:"required"` // 名称
	Description    string  `json:"description"     gorm:"type:text"`                                            // 描述
	Avatar         string  `json:"avatar"`                                                                      // LOGO
	Address        string  `json:"address"`                                                                     // 地址
	ZipCode        string  `json:"zip_code"`                                                                    // 邮政编码
	Phone          string  `json:"phone"`                                                                       // 电话号码
	Fax            string  `json:"fax"`                                                                         // 传真
	Position       float64 `json:"position"`                                                                    // 排序位置

	CreatedAt int64 `json:"created_at,omitempty"` // 创建时间
	UpdatedAt int64 `json:"updated_at,omitempty"` // 更新时间

	CountMembers int `json:"count_members" sql:"-"`
}

package main

import "time"

// TestStruct this is a test struct
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

	CreatedAt int64 `json:"created_at,omitempty" xml:"created_at"           bson:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty" xml:"updated_at,omitempty" bson:"created_at"`
}

type Fset struct{}

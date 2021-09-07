package tpl

import (
	"fmt"
)

type Attribute map[string]interface{}

func (a Attribute) Exists(key string) bool {
	if a == nil {
		return false
	}

	_, ok := a[key]
	return ok
}

func (a Attribute) GetString(key string) string {
	if a == nil {
		return ""
	}

	val, ok := a[key]
	if !ok {
		return ""
	}

	s, ok := val.(string)
	if !ok {
		return ""
	}
	return s
}

func (a Attribute) GetBool(key string) bool {
	if a == nil {
		return false
	}

	val, ok := a[key]
	if !ok {
		return false
	}

	s, ok := val.(bool)
	if !ok {
		return false
	}
	return s
}

const (
	GormFieldTag_Type          = "type"
	GormFieldTag_Size          = "size"
	GormFieldTag_Primary       = "primary"
	GormFieldTag_Unique        = "unique"
	GormFieldTag_Default       = "default"
	GormFieldTag_Precision     = "precision"
	GormFieldTag_NotNull       = "not_null"
	GormFieldTag_AutoIncrement = "auto_increment"
	GormFieldTag_Index         = "index"
	GormFieldTag_UniqueIndex   = "unique_index"
)

var GormFieldTagEnabled = []string{
	GormFieldTag_Type,
	GormFieldTag_Size,
	GormFieldTag_Primary,
	GormFieldTag_Unique,
	GormFieldTag_Default,
	GormFieldTag_Precision,
	GormFieldTag_NotNull,
	GormFieldTag_AutoIncrement,
	GormFieldTag_Index,
	GormFieldTag_UniqueIndex,
}

func (a Attribute) GetAllGormTag() (res string) {
	for _, v := range GormFieldTagEnabled {
		res += a.GetGormTag(v)
	}
	return res
}

func (a Attribute) GetGormTag(key string) string {
	if a == nil {
		return ""
	}

	val, ok := a[key]
	if !ok {
		return ""
	}

	switch key {
	case GormFieldTag_Type:
		return fmt.Sprintf(";type:%v", val)
	case GormFieldTag_Size:
		return fmt.Sprintf(";size:%v", val)
	case GormFieldTag_Primary:
		return ";primary_key"
	case GormFieldTag_Unique:
		return ";unique"
	case GormFieldTag_Default:
		return fmt.Sprintf(";default:%v", val)
	case GormFieldTag_Precision:
		return fmt.Sprintf(";precision:%v", val)
	case GormFieldTag_NotNull:
		return ";not null"
	case GormFieldTag_AutoIncrement:
		str, ok2 := val.(bool)
		if !ok2 {
			return ";auto_increment"
		}
		if str {
			return ";auto_increment"
		} else {
			return ";auto_increment:false"
		}
	case GormFieldTag_Index:
		fallthrough
	case GormFieldTag_UniqueIndex:
		str, ok2 := val.(string)
		if !ok2 {
			return ";" + key
		}
		if len(str) == 0 {
			return ";" + key
		} else {
			return fmt.Sprintf(";%v:%v", key, str)
		}
	default:
		return ""
	}

	/*
		Column	Specifies column name

		Type	Specifies column data type
		Size	Specifies column size, default 255
		PRIMARY_KEY	Specifies column as primary key
		UNIQUE	Specifies column as unique
		DEFAULT	Specifies column default value
		PRECISION	Specifies column precision
		NOT NULL	Specifies column as NOT NULL
		AUTO_INCREMENT	Specifies column auto incrementable or not
		INDEX	Create index with or without name, same name creates composite indexes
		UNIQUE_INDEX	Like INDEX, create unique index
	*/
}

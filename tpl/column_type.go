package tpl

type ColumnType string

// Constants for database(mysql) types
const (
	ColumnTypeInt8    = "int8"
	ColumnTypeUInt8   = "uint8"
	ColumnTypeInt16   = "int16"
	ColumnTypeUInt16  = "uint16"
	ColumnTypeInt32   = "int32"
	ColumnTypeUInt32  = "uint32"
	ColumnTypeInt64   = "int64"
	ColumnTypeUInt64  = "uint64"
	ColumnTypeFloat32 = "float32"
	ColumnTypeFloat64 = "float64"

	ColumnTypeBool    = "bool"
	ColumnTypeVarchar = "varchar"
	ColumnTypeBinary  = "blob"
	ColumnTypeDecimal = "decimal"

	ColumnTypeTime     = "time"
	ColumnTypeDate     = "date"
	ColumnTypeDateTime = "datetime"
)

// Constants for return types of golang
const (
	golangByteArray = "[]byte"
	sqlNullInt      = "sql.NullInt64"
	golangInt       = "int"
	golangInt64     = "int64"
	sqlNullFloat    = "sql.NullFloat64"
	golangFloat     = "float"
	golangFloat32   = "float32"
	golangFloat64   = "float64"
	sqlNullString   = "sql.NullString"
	golangTime      = "time.Time"
	sqlNullTime     = "sql.NullTime"
)

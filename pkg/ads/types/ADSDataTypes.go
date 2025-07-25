package types

// ADSDataType is the ADS data type
type ADSDataType uint32

const (
	ADST_VOID     ADSDataType = 0  // ADST_VOID is a void type
	ADST_INT16    ADSDataType = 2  // ADST_INT16 is a 16-bit integer
	ADST_INT32    ADSDataType = 3  // ADST_INT32 is a 32-bit integer
	ADST_REAL32   ADSDataType = 4  // ADST_REAL32 is a 32-bit real
	ADST_REAL64   ADSDataType = 5  // ADST_REAL64 is a 64-bit real
	ADST_INT8     ADSDataType = 16 // ADST_INT8 is an 8-bit integer
	ADST_UINT8    ADSDataType = 17 // ADST_UINT8 is an 8-bit unsigned integer
	ADST_UINT16   ADSDataType = 18 // ADST_UINT16 is a 16-bit unsigned integer
	ADST_UINT32   ADSDataType = 19 // ADST_UINT32 is a 32-bit unsigned integer
	ADST_INT64    ADSDataType = 20 // ADST_INT64 is a 64-bit integer
	ADST_UINT64   ADSDataType = 21 // ADST_UINT64 is a 64-bit unsigned integer
	ADST_STRING   ADSDataType = 30 // ADST_STRING is a string
	ADST_WSTRING  ADSDataType = 31 // ADST_WSTRING is a wide string
	ADST_REAL80   ADSDataType = 32 // ADST_REAL80 is an 80-bit real
	ADST_BIT      ADSDataType = 33 // ADST_BIT is a bit
	ADST_MAXTYPES ADSDataType = 34 // ADST_MAXTYPES is the maximum number of types
	ADST_BIGTYPE  ADSDataType = 65 // ADST_BIGTYPE is a big type
)

// ADSDataTypeToString converts an ADSDataType to a string.
func ADSDataTypeToString(dataType ADSDataType) string {
	switch dataType {
	case ADST_VOID:
		return "VOID"
	case ADST_INT16:
		return "INT16"
	case ADST_INT32:
		return "INT32"
	case ADST_REAL32:
		return "REAL32"
	case ADST_REAL64:
		return "REAL64"
	case ADST_INT8:
		return "INT8"
	case ADST_UINT8:
		return "UINT8"
	case ADST_UINT16:
		return "UINT16"
	case ADST_UINT32:
		return "UINT32"
	case ADST_INT64:
		return "INT64"
	case ADST_UINT64:
		return "UINT64"
	case ADST_STRING:
		return "STRING"
	case ADST_WSTRING:
		return "WSTRING"
	case ADST_REAL80:
		return "REAL80"
	case ADST_BIT:
		return "BIT"
	case ADST_MAXTYPES:
		return "MAXTYPES"
	case ADST_BIGTYPE:
		return "BIGTYPE"
	default:
		return "UNKNOWN"
	}
}

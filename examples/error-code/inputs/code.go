package inputs

// ErrorCode is an enumerated message that corresponds to the status
type ErrorCode string

// ErrCodeCommonToMessage is map of error code to their messages
var ErrCodeCommonToMessage = map[int]string{
	ErrCodeCommonUnspecifiedError:   "unspecified error of %s",
	ErrCodeCommonJSONUnmarshalError: "failed to unmarshal json %s",
	ErrCodeCommonJSONMarshalError:   "failed to marshal json %s",
}

// ErrCodeCommonPrefix is the prefix of common error message
// +ErrCode=Common
const ErrCodeCommonPrefix = "01-01"

// list of common error codes
// +ErrCode=Common,
const (
	ErrCodeCommonBegin = iota
	ErrCodeCommonUnspecifiedError
	ErrCodeCommonJSONUnmarshalError
	ErrCodeCommonJSONMarshalError
	ErrCodeCommonEnd
)

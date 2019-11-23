package inputs

// Code defines the type of error code
type Code int

// ErrCodeRbd defines the id of rbd module
// +ErrCode
const ErrCodeRbd Code = 0x02

// the sub module of rbd
// +ErrCode=Rbd
const (
	ErrCodeRbdCommon Code = iota + 1
	ErrCodeRbdVolume
)

// list of rbd common error codes
// +ErrCode=Rbd,Common
const (
	ErrCodeRbdCommonBegin Code = iota
	ErrCodeRbdCommonUnspecifiedError
	ErrCodeRbdCommonEnd
)

// ErrCodeRbdCommonToMessage is map of common error code to their messages
var ErrCodeRbdCommonToMessage = map[Code]string{
	ErrCodeRbdCommonUnspecifiedError: "The %s operation failed due to unspecified error",
}

// list of rbd volume error codes
// +ErrCode=Rbd,Volume
const (
	ErrCodeRbdVolumeBegin = iota
	ErrCodeRbdVolumeUnknownParameter
	ErrCodeRbdVolumeNoEntry
	ErrCodeRbdVolumeEnd
)

// ErrCodeRbdVolumeToMessage is map of volume error code to their messages
var ErrCodeRbdVolumeToMessage = map[int]string{
	ErrCodeRbdVolumeUnknownParameter: "The %s operation failed due to unknown parameter",
	ErrCodeRbdVolumeNoEntry:          "The %s operation failed due to no such directory entry",
}

// Only for test
const (
	// only for test A
	TestA = iota
	// only for test B
	TestB

	TestC
	_
	_
	TestD
	TestE
	TestF
)

// hi, hello

var test string

const testString = "we are Chinese!"

const testBool = true

const testFloat = 3.14

const testInt = 234

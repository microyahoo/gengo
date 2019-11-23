package inputs

import (
	"fmt"
)

// ErrorCode is an enumerated message that corresponds to the status
// An error code is to identify fault, its pattern is "AA-BB-CCCC", contain three parts:
//   AA code identifies module or component. For examples, rbd, ceph, net etc.
//   BB code identifies sub module. For examples, volume, snapshot in `rbd` module.
//   CCCC code identifies the concrete fault id.
//
// Notice:
// 1. The format of the error code is uppercase hexadecimal.
// 2. The sub module stands for common part if BB is equal to 01.
// 3. The id of 0000 is reserved and 0001 stands for unspecified error.
//
// eg: 02-02-0002
// +-------------------+---------------------+-----------------------+
// |         02        |          02         |         0002          |
// +-------------------+---------------------+-----------------------+
// | `rbd` module      | `volume` sub module | volume already exists |
// +-------------------+---------------------+-----------------------+
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

// ErrCodeDescriptor defines the unique code of error info
// the format is AA-BB-CCCC
type ErrCodeDescriptor struct {
	Module    int64
	SubModule int64
	ID        int64
}

func (code *ErrCodeDescriptor) String() string {
	return fmt.Sprintf("%02X-%02X-%04X", code.Module, code.SubModule, code.ID)
}

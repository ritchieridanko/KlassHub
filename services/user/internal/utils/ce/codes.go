package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Internal Errors
var (
	ErrDBAffectNoRows error = errors.New("no rows affected")
	ErrDBQueryNoRows  error = pgx.ErrNoRows
)

// Internal Error Codes
const (
	CodeDBQueryExec          errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTransaction        errCode = "ERR_DB_TRANSACTION"
	CodeUnknown              errCode = "ERR_UNKNOWN"
	CodeUserNotFound         errCode = "ERR_USER_NOT_FOUND"
	CodeUUIDGenerationFailed errCode = "ERR_UUID_GENERATION_FAILED"
)

// External Error Messages
const (
	MsgInternalServer string = "Internal server error"
	MsgUserNotFound   string = "User not found"
)

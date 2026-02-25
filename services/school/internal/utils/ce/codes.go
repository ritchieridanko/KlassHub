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
	CodeDBQueryExec    errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTransaction  errCode = "ERR_DB_TRANSACTION"
	CodeSchoolNotFound errCode = "ERR_SCHOOL_NOT_FOUND"
	CodeUnknown        errCode = "ERR_UNKNOWN"
)

// External Error Messages
const (
	MsgInternalServer string = "Internal server error"
	MsgSchoolNotFound string = "School not found"
)

package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Internal Errors
var (
	ErrDBAffectNoRows error = errors.New("no rows affected")
	ErrDBQueryNoRows  error = pgx.ErrNoRows
	ErrEventOnProcess error = errors.New("event is being processed by another subscriber")
)

// Internal Error Codes
const (
	CodeDBQueryExec           errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTransaction         errCode = "ERR_DB_TRANSACTION"
	CodeEmailDeliveryFailed   errCode = "ERR_EMAIL_DELIVERY_FAILED"
	CodeEmailTemplatingFailed errCode = "ERR_EMAIL_TEMPLATING_FAILED"
	CodeEventCommittingFailed errCode = "ERR_EVENT_COMMITTING_FAILED"
	CodeEventFetchingFailed   errCode = "ERR_EVENT_FETCHING_FAILED"
	CodeEventNotFound         errCode = "ERR_EVENT_NOT_FOUND"
	CodeEventOnProcess        errCode = "ERR_EVENT_ON_PROCESS"
	CodeJSONRawEncodingFailed errCode = "ERR_JSON_RAW_ENCODING_FAILED"
	CodePanicOccurred         errCode = "ERR_PANIC_OCCURRED"
	CodeProtobufParsingFailed errCode = "ERR_PROTOBUF_PARSING_FAILED"
	CodeTypeConversionFailed  errCode = "ERR_TYPE_CONVERSION_FAILED"
	CodeURLGenerationFailed   errCode = "ERR_URL_GENERATION_FAILED"
)

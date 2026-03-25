package ce

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FromGRPCErr(err error) *Error {
	st, ok := status.FromError(err)
	if !ok {
		return NewError(CodeUnknown, MsgInternalServer, err)
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return NewError(CodeInvalidPayload, st.Message(), err)
	case codes.Unauthenticated:
		return NewError(CodeUnauthenticated, st.Message(), err)
	case codes.FailedPrecondition, codes.PermissionDenied:
		return NewError(CodeUnauthorized, st.Message(), err)
	case codes.NotFound:
		return NewError(CodeNotFound, st.Message(), err)
	case codes.AlreadyExists:
		return NewError(CodeAlreadyExists, st.Message(), err)
	case codes.Internal:
		return NewError(CodeInternal, st.Message(), err)
	default:
		return NewError(CodeInternal, st.Message(), err)
	}
}

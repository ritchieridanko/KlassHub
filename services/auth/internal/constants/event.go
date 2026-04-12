package constants

// Event Topics
var (
	EventTopicAC  string = "auth.created"
	EventTopicAVR string = "auth.verification.requested"
	EventTopicUCF string = "user.creation.failed"
)

// Event Header Keys
const (
	EventHeaderKeyRequestID string = "x-request-id"
)

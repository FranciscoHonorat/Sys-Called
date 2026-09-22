package valueobjects

type Status string

const (
	TicketStatusOpen       Status = "Open"
	TicketStatusInProgress Status = "In Progress"
	TicketStatusClosed     Status = "Closed"
)

func NewStatus(status string) Status {
	return Status(status)
}

func (s Status) IsValid() bool {
	return isOneOf(s, TicketStatusOpen, TicketStatusInProgress, TicketStatusClosed)
}

func (s Status) String() string {
	return string(s)
}

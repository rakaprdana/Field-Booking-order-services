package constants

type FieldStatusString string

const (
	AvailableFieldStatus FieldStatusString = "available"
	BookedFieldStatus    FieldStatusString = "booked"
)

func (f FieldStatusString) String() string {
	return string(f)
}

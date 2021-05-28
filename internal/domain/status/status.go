package status

type LinkStatus int

const (
	Unknown LinkStatus = iota
	OK
	Failed
)

func (s LinkStatus) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case OK:
		return "OK"
	case Failed:
		return "Failed"
	default:
		panic("Unknown status")
	}
}

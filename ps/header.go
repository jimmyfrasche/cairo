package ps

type dscSection int

const (
	Header dscSection = iota
	Setup
	PageSetup
)

func (d dscSection) String() (s string) {
	switch d {
	case Header:
		s = "Header"
	case Setup:
		s = "Setup"
	case PageSetup:
		s = "Page Setup"
	default:
	}
	return s + " header section"
}

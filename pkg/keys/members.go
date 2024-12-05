package keys

import (
	"github.com/go-ole/go-ole"
)

// Chrome CLSID_Elevator and IID_IElevator
var (
	ChromeCLSIDElevator = ole.NewGUID("{708860E0-F641-4611-8895-7D867DD3675B}")
	ChromeIIDIElevator  = ole.NewGUID("{463ABECF-410D-407F-8AF5-0DF35A005CC8}")
)

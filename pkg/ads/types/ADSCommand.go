package types

// ADSCommand defines the ADS commands.
type ADSCommand uint16

const (
	ADSCommandInvalid            ADSCommand = 0 // Invalid
	ADSCommandNone               ADSCommand = 0 // uninitialized
	ADSCommandReadDeviceInfo     ADSCommand = 1 // ReadDeviceInfo command
	ADSCommandRead               ADSCommand = 2 // Read Command
	ADSCommandWrite              ADSCommand = 3 // Write Command
	ADSCommandReadState          ADSCommand = 4 // Read State Command
	ADSCommandWriteControl       ADSCommand = 5 // WriteControl Command
	ADSCommandAddNotification    ADSCommand = 6 // Add Device notification
	ADSCommandDeleteNotification ADSCommand = 7 // Delete Device notification
	ADSCommandNotification       ADSCommand = 8 // Device notification
	ADSCommandReadWrite          ADSCommand = 9 // Read Write command
)

func (c ADSCommand) String() string {
	switch c {
	case ADSCommandReadDeviceInfo:
		return "ReadDeviceInfo"
	case ADSCommandRead:
		return "Read"
	case ADSCommandWrite:
		return "Write"
	case ADSCommandReadState:
		return "ReadState"
	case ADSCommandWriteControl:
		return "WriteControl"
	case ADSCommandAddNotification:
		return "AddNotification"
	case ADSCommandDeleteNotification:
		return "DeleteNotification"
	case ADSCommandNotification:
		return "Notification"
	case ADSCommandReadWrite:
		return "ReadWrite"
	default:
		return "UNKNOWN"
	}
}

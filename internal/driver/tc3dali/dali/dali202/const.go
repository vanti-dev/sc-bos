package dali202

const (
	StatusBitInhibit uint8 = 1 << iota
	StatusBitFunctionTestDone
	StatusBitDurationTestDone
	StatusBitBatteryFull
	StatusBitFunctionTestPending
	StatusBitDurationTestPending
	StatusBitIdentificationActive
	StatusBitPhysicallySelected
)

const (
	ModeBitRestActive uint8 = 1 << iota
	ModeBitNormalModeActive
	ModeBitEmergencyModeActive
	ModeBitExtendedEmergencyModeActive
	ModeBitFunctionTestInProgress
	ModeBitDurationTestInProgress
	ModeBitHardwiredInhibit
	ModeBitHardwiredSwitch
)

const (
	FailureBitCircuit uint8 = 1 << iota
	FailureBitBatteryDuration
	FailureBitBattery
	FailureBitEmergencyLamp
	FailureBitFunctionMaxDelayExceeded
	FailureBitDurationMaxDelayExceeded
	FailureBitFunctionTest
	FailureBitDurationTest
)

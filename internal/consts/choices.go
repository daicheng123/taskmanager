package consts

type (
	ExecutorStatus = uint
	SecretStatus   = uint
	ScriptStatus   = uint
	ScriptType     = uint
	TaskMode       = uint
	TaskStatus     = uint
	TransferStatus = uint
	AnsibleTask    = string
)

const (
	HostUnknown ExecutorStatus = iota + 1
	HostAvail
	HostUnreachable
)
const (
	KeyUndistributed = iota + 1
	KeyDistributing
	KeyDistributed
	KeyDistributeFailed
)

const (
	ManualTask TaskMode = iota + 1
	AutoTask
)

const (
	TaskExecuting TaskStatus = iota + 1
	TaskSuccess
	TaskFailed
	TaskTimeout
	TaskAbort
	TaskNOtExecute
)

const (
	NoAudit ScriptStatus = iota + 1
	Auditing
	PassAudit
	FailAudit
)
const (
	Shell ScriptType = iota + 1
	Python
)

const (
	TransferSuccess TransferStatus = iota + 1
	TransferFailed
)

const (
	TransferScript AnsibleTask = "transfer script"
	ExecuteScript  AnsibleTask = "execute script"
	DeleteScript   AnsibleTask = "delete script"
)

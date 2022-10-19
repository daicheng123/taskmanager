package utils

type (
	ExecutorStatus = uint
	SecretStatus   = uint
	ScriptStatus   = uint
	ScriptType     = uint
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
	NoAudit ScriptStatus = iota + 1
	Auditing
	PassAudit
	FailAudit
)
const (
	Shell ScriptType = iota + 1
	Python
)

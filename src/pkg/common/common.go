package common

const SqsMessageBodySeparator = "|"
const RdsTypeClusterKey = "cluster"
const RdsTypeInstanceKey = "instance"

type MySQLConnectionInfo struct {
	Endpoint       string
	Port           int32
	Username       string
	Password       string
	SecurityGroups []string
}

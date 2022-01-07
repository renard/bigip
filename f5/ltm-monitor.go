package f5

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

// Monitor is a real mess to parse for a limited gain. Report it as
// plain text for now.
// See:
// https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_monitor_http.html
type LtmMonitor struct {
	OriginalConfig ParsedConfig
	Pos            lexer.Position
	Type           string `("ltm" "monitor" @Ident | "monitor")`
	Name           string `@( F5Name | QF5Name | Ident) "{"`
	// Type           string `("ltm" "monitor" @Ident | "monitor" @Ident )`
	// Name           string `( @F5Name | @QF5Name )"{"?` // + version 10.x
	// Adaptive                  string  `(  "adaptive" @( "enabled" | "disabled" )`
	// AdaptiveDivergenceType    string  ` | "adaptive-divergence-type" @( "relative" | "absolute" )`
	// AdaptiveDivergenceValue   int     ` | "adaptive-divergence-value" @( Ident )`
	// AdaptiveDivergenceLimit   int     ` | "adaptive-divergence-limit" @( Ident )`
	// AdaptiveSampelingTimespan int     ` | "adaptive-sampling-timespan" @( Ident )`
	// AgentType                 string  ` | "agent-type" @( Ident )`
	// AppService                string  ` | "app-service" @( F5Name | QF5Name )`
	// Cipherlist                string  ` | "cipherlist" @( Ident )`
	// Community                 string  ` | "community" @( Ident )`
	// Compatibility             string  ` | "compatibility" @( "enabled" | "disabled" )`
	// Count                     int     ` | "count" @( Ident )`
	// CpuCoefficient            float64 ` | "cpu-coefficient" @Ident`
	// CpuThreshold              int     ` | "cpu-threshold" @Ident`
	// Database                  string  ` | "database" @( Ident | QString)`
	// Debug                     string  ` | "debug" @( "yes" | "no" )`
	// DefaultsFrom              string  ` | "defaults-from" @( F5Name | QF5Name )`
	// Description               string  ` | "description" @( Ident | QString )`
	// Destination               string  ` | "destination" @( F5IPPort | Ident )`
	// DiskCoefficient           float64 ` | "disk-coefficient" @Ident`
	// DiskThreshold             int     ` | "disk-threshold" @Ident`
	// Interval                  int     ` | "interval" @( Ident )`
	// IpTos                     int     ` | "ip-tos" @( Ident )`
	// ManualResume              string  ` | "manual-resume" @( "enabled" | "disabled" )`
	// MemoryCoefficient         float64 ` | "memory-coefficient" @Ident`
	// MemoryThreshold           int     ` | "memory-threshold" @Ident`
	// Password                  string  ` | "password" @( F5Password | Ident | QString )`
	// Recv                      string  ` | "recv" @( Ident | QString )`
	// RecvDisable               string  ` | "recv-disable" @( "none" | Ident | QString )`
	// Reverse                   string  ` | "reverse" @( "enabled" | "disabled" )`
	// IpDscp                    int     ` | "ip-dscp" @( Ident )`
	// Send                      string  ` | "send" @( "none" | Ident | F5Name | QF5Name | QString )`
	// TimeUntilUp               int     ` | "time-until-up" @( Ident )`
	// Timeout                   int     ` | "timeout" @( Ident )`
	// Transparent               string  ` | "transparent" @( "enabled" | "disabled" )`
	// Upinterval                int     ` | "up-interval" @( Ident )`
	// Version                   string  ` | "version" @( Ident )`
	// Username                  string  ` | "username" @( Ident | F5Name | QF5Name ) )* "}"`
}

// newLtmMonitor parses data and creates a new LtmMonitor struct.
func newLtmMonitor(data ParsedConfig) (ret *LtmMonitor, err error) {
	ret = &LtmMonitor{}
	// Only analyze first line
	monitor := strings.Split(data.Content, "\n")
	err = parseString("", monitor[0], ret)
	ret.OriginalConfig = data
	return
}

func (o *LtmMonitor) Original() string {
	return o.OriginalConfig.Content
}

func (o *LtmMonitor) GetName() string {
	return o.Name
}

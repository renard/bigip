package f5

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// See https://clouddocs.f5.com/cli/tmsh-reference/latest/modules/ltm/ltm_pool.html
type LtmPool struct {
	OriginalConfig         ParsedConfig
	Pos                    lexer.Position
	Name                   string           `("ltm" "pool" ( @F5Name | @QF5Name | @Ident ) | "pool" ( @F5Name | @QF5Name | @Ident)) "{"` // + version 10.x
	AllowNat               string           `(  "allow-nat" @( "yes" | "no" )`
	AllowSnat              string           ` | "allow-snat" @( "yes" | "no" )`
	AppService             string           ` | "app-service" @( "none" | QString | Ident )`
	AutoscaleGroupID       string           ` | "autoscale-group-id" @( "none" | QString | Ident )`
	Description            string           ` | "description" @( QString | Ident )`
	GatewayFailsafeDevice  string           ` | "gateway-failsafe-device" @( QString | Ident )`
	IgnorePersistedWeight  string           ` | "ignore-persisted-weight" @( "yes" | "no" | "enabled" )`
	LoadBalancingMode      string           ` | "load-balancing-mode" @( F5lbMode  )` // TODO: limit algos
	Members                *[]LtmPoolMember ` | "members" ("{" @@* "}" | @@ ) `
	MinActiveMembers       int              ` | "min-active-members" @Ident`
	MinUpMembers           int              ` | "min-up-members" @Ident`
	MinUpMembersAction     string           ` | "min-up-members-action" @Ident`
	MinUpMembersChecking   string           ` | "min-up-members-checking" @Ident`
	Monitor                []string         ` | ("monitor" ( ( "min" Ident "of" "{" @( F5Name | QF5Name ) ("and"? @(F5Name | QF5Name) )* "}" ) | @( "none" | F5Name | QF5Name | (?!"all") Ident ) ("and"? @(F5Name | QF5Name | Ident) )* )) | ( "monitor" "all" @( Ident ) ("and" @( Ident ) )*  )` // all and indent for v10.x |  (?!"all") Ident -> do not match all ...
	Profiles               string           ` | "profiles" @( "none" | @F5Name | @QF5Name )`
	QueueOnConnectionLimit string           ` | "queue-on-connection-limit" @( "enabled" | "disabled" )`
	QueueDepthLimit        int              ` | "queue-depth-limit" @Ident`
	QueueTimeLimit         int              ` | "queue-time-limit" @Ident`
	ReselectTries          int              ` | "reselect-tries" @Ident`
	ServiceDownAction      string           ` | "service-down-action" @Ident`
	SlowRampTime           int              ` | ("slow-ramp-time" | "slow" "ramp" "time") @Ident  )* "}"`
}

type LtmPoolMember struct {
	Name            string       `@(F5Name | QF5Name | Ident) "{"`
	Address         string       ` ( "address" @Ident`
	AppService      string       ` | "app-service" @( "none" | QString | Ident )`
	ConnectionLimit int          ` | "connection-limit" @Ident`
	Description     string       ` | "description" @( QString | Ident )`
	DynamicRatio    int          ` | "dynamic-ratio" @Ident`
	InheritProfile  *bool        ` | "inherit-profile" @("enabled" | "disabled")`
	Logging         *bool        ` | "logging" @("enabled" | "disabled")`
	Monitor         []string     ` | "monitor" ( ( "min" Ident "of" "{" @( F5Name | QF5Name ) ("and"? @(F5Name | QF5Name) )* "}" ) | @( "none" | F5Name | QF5Name | Ident) ("and"? @(F5Name | QF5Name | Ident) )* )`
	PriorityGroup   string       ` | "priority-group" @( "none" | Ident )`
	RateLimit       int          ` | "rate-limit" @Ident`
	Ratio           int          ` | "ratio" @Ident`
	Session         string       ` | "session" @( "user-enabled" | "user-disabled" | "user" "disabled" | "monitor-enabled")`
	State           string       ` | "state" @( "user-up" | "user-down" | "down" | "up")`
	FQDN            *LtmPoolFQDN ` | "fqdn" "{" @@ "}" )* "}"`
}

type LtmPoolFQDN struct {
	Autopopulate string `(  "autopopulate" @("enabled" | "disabled")`
	Name         string ` | "name" @Ident)*`
}

// newLtmPool parses data and creates a new LtmPool struct.
func newLtmPool(data ParsedConfig) (ret *LtmPool, err error) {
	ret = &LtmPool{}
	err = parseString("", data.Content, ret)
	ret.OriginalConfig = data
	return
}

func (o *LtmPool) Original() string {
	return o.OriginalConfig.Content
}

func (o *LtmPool) GetName() string {
	return o.Name
}

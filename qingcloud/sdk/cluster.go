package sdk

import (
	"fmt"
	"time"

	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/request"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	"github.com/yunify/qingcloud-sdk-go/request/errors"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

var _ fmt.State
var _ time.Time

type ClusterService struct {
	Config     *config.Config
	Properties *ClusterServiceProperties
}

type ClusterServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func NewCluster(s *qc.QingCloudService, zone string) (*ClusterService, error) {
	properties := &ClusterServiceProperties{
		Zone: &zone,
	}

	return &ClusterService{Config: s.Config, Properties: properties}, nil
}

func (s *ClusterService) DescribeClusters(i *DescribeClustersInput) (*DescribeClustersOutput, error) {
	if i == nil {
		i = &DescribeClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeClusters",
		RequestMethod: "GET",
	}

	x := &DescribeClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeClustersInput struct {
	AppVersions       []*string `json:"app_versions" name:"app_versions" location:"params"`
	Apps              []*string `json:"apps" name:"apps" location:"params"`
	CfgmgmtID         *string   `json:"cfgmgmt_id" name:"cfgmgmt_id" location:"params"`
	Clusters          []*string `json:"clusters" name:"clusters" location:"params"`
	Console           *string   `json:"console" name:"console" location:"params"`
	ExternalClusterID *string   `json:"external_cluster_id" name:"external_cluster_id" location:"params"`
	Limit             *int      `json:"limit" name:"limit" location:"params"`
	Link              *string   `json:"link" name:"link" location:"params"`
	Name              *string   `json:"name" name:"name" location:"params"`
	Offset            *int      `json:"offset" name:"offset" location:"params"`
	Owner             *string   `json:"owner" name:"owner" location:"params"`
	Reverse           *int      `json:"reverse" name:"reverse" location:"params"`
	Role              *string   `json:"role" name:"role" location:"params"`
	// Scope's available values: all, cfgmgmt
	Scope            *string   `json:"scope" name:"scope" location:"params"`
	SearchWord       *string   `json:"search_word" name:"search_word" location:"params"`
	SortKey          *string   `json:"sort_key" name:"sort_key" location:"params"`
	Status           *string   `json:"status" name:"status" location:"params"`
	TransitionStatus *string   `json:"transition_status" name:"transition_status" location:"params"`
	Users            []*string `json:"users" name:"users" location:"params"`
	Verbose          *int      `json:"verbose" name:"verbose" location:"params"`
	VxNet            *string   `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *DescribeClustersInput) Validate() error {

	if v.Scope != nil {
		scopeValidValues := []string{"all", "cfgmgmt"}
		scopeParameterValue := fmt.Sprint(*v.Scope)

		scopeIsValid := false
		for _, value := range scopeValidValues {
			if value == scopeParameterValue {
				scopeIsValid = true
			}
		}

		if !scopeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Scope",
				ParameterValue: scopeParameterValue,
				AllowedValues:  scopeValidValues,
			}
		}
	}

	return nil
}

type DescribeClustersOutput struct {
	Message    *string    `json:"message" name:"message"`
	Action     *string    `json:"action" name:"action" location:"elements"`
	ClusterSet []*Cluster `json:"cluster_set" name:"cluster_set" location:"elements"`
	RetCode    *int       `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int       `json:"total_count" name:"total_count" location:"elements"`
}

type Cluster struct {
	AdvancedActions            map[string]*string `json:"advanced_actions" name:"advanced_actions"`
	AppID                      *string            `json:"app_id" name:"app_id"`
	AppInfo                    interface{}        `json:"app_info" name:"app_info"`
	AppVersion                 *string            `json:"app_version" name:"app_version"`
	AppVersionInfo             interface{}        `json:"app_version_info" name:"app_version_info"`
	AutoBackupTime             *int               `json:"auto_backup_time" name:"auto_backup_time"`
	Backup                     map[string]*bool   `json:"backup" name:"backup"`
	BackupPolicy               *string            `json:"backup_policy" name:"backup_policy"`
	BackupService              interface{}        `json:"backup_service" name:"backup_service"`
	CfgmgmtID                  *string            `json:"cfgmgmt_id" name:"cfgmgmt_id"`
	ClusterID                  *string            `json:"cluster_id" name:"cluster_id"`
	ClusterType                *int               `json:"cluster_type" name:"cluster_type"`
	ConsoleID                  *string            `json:"console_id" name:"console_id"`
	Controller                 *string            `json:"controller" name:"controller"`
	CreateTime                 *time.Time         `json:"create_time" name:"create_time" format:"ISO 8601"`
	CustomService              interface{}        `json:"custom_service" name:"custom_service"`
	Debug                      *bool              `json:"debug" name:"debug"`
	Description                *string            `json:"description" name:"description"`
	DisplayTabs                interface{}        `json:"display_tabs" name:"display_tabs"`
	Endpoints                  interface{}        `json:"endpoints" name:"endpoints"`
	GlobalUUID                 *string            `json:"global_uuid" name:"global_uuid"`
	HealthCheckEnablement      map[string]*bool   `json:"health_check_enablement" name:"health_check_enablement"`
	IncrementalBackupSupported *bool              `json:"incremental_backup_supported" name:"incremental_backup_supported"`
	LatestSnapshotTime         *string            `json:"latest_snapshot_time" name:"latest_snapshot_time"`
	Links                      map[string]*string `json:"links" name:"links"`
	MetadataRootAccess         *bool              `json:"metadata_root_access" name:"metadata_root_access"`
	Name                       *string            `json:"name" name:"name"`
	NodeCount                  *int               `json:"node_count" name:"node_count"`
	Nodes                      []*ClusterNode     `json:"nodes" name:"nodes"`
	Owner                      *string            `json:"owner" name:"owner"`
	PartnerAccess              *bool              `json:"partner_access" name:"partner_access"`
	RestoreService             interface{}        `json:"restore_service" name:"restore_service"`
	ReuseHyper                 *bool              `json:"reuse_hyper" name:"reuse_hyper"`
	RoleCount                  map[string]*int    `json:"role_count" name:"role_count"`
	Roles                      []*string          `json:"roles" name:"roles"`
	RootUserID                 *string            `json:"root_user_id" name:"root_user_id"`
	SecurityGroupID            *string            `json:"security_group_id" name:"security_group_id"`
	Status                     *string            `json:"status" name:"status"`
	StatusTime                 *time.Time         `json:"status_time" name:"status_time" format:"ISO 8601"`
	SubCode                    *int               `json:"sub_code" name:"sub_code"`
	TransitionStatus           *string            `json:"transition_status" name:"transition_status"`
	UpgradePolicy              []interface{}      `json:"upgrade_policy" name:"upgrade_policy"`
	UpgradeStatus              *string            `json:"upgrade_status" name:"upgrade_status"`
	UpgradeTime                *time.Time         `json:"upgrade_time" name:"upgrade_time" format:"ISO 8601"`
	VxNet                      *VxNet             `json:"vxnet" name:"vxnet"`
}

type ClusterNode struct {
	AdvancedActions            *string     `json:"advanced_actions" name:"advanced_actions"`
	AgentInstalled             *bool       `json:"agent_installed" name:"agent_installed"`
	AlarmStatus                *string     `json:"alarm_status" name:"alarm_status"`
	AppID                      *string     `json:"app_id" name:"app_id"`
	AppVersion                 *string     `json:"app_version" name:"app_version"`
	AutoBackup                 *int        `json:"auto_backup" name:"auto_backup"`
	BackupPolicy               *string     `json:"backup_policy" name:"backup_policy"`
	BackupService              interface{} `json:"backup_service" name:"backup_service"`
	ClusterID                  *string     `json:"cluster_id" name:"cluster_id"`
	ConsoleID                  *string     `json:"console_id" name:"console_id"`
	Controller                 *string     `json:"controller" name:"controller"`
	CPU                        *int        `json:"cpu" name:"cpu"`
	CreateTime                 *time.Time  `json:"create_time" name:"create_time" format:"ISO 8601"`
	CustomMetadataScript       interface{} `json:"custom_metadata_script" name:"custom_metadata_script"`
	CustomService              interface{} `json:"custom_service" name:"custom_service"`
	Debug                      *bool       `json:"debug" name:"debug"`
	DestroyService             interface{} `json:"destroy_service" name:"destroy_service"`
	DisplayTabs                interface{} `json:"display_tabs" name:"display_tabs"`
	EIP                        *string     `json:"eip" name:"eip"`
	Env                        *string     `json:"env" name:"env"`
	GlobalServerID             *int        `json:"global_server_id" name:"global_server_id"`
	Gpu                        *int        `json:"gpu" name:"gpu"`
	GpuClass                   *int        `json:"gpu_class" name:"gpu_class"`
	GroupID                    *int        `json:"group_id" name:"group_id"`
	HealthCheck                interface{} `json:"health_check" name:"health_check"`
	HealthStatus               *string     `json:"health_status" name:"health_status"`
	Hypervisor                 *string     `json:"hypervisor" name:"hypervisor"`
	ImageID                    *string     `json:"image_id" name:"image_id"`
	IncrementalBackupSupported *bool       `json:"incremental_backup_supported" name:"incremental_backup_supported"`
	InitService                interface{} `json:"init_service" name:"init_service"`
	InstanceID                 *string     `json:"instance_id" name:"instance_id"`
	IsBackup                   *int        `json:"is_backup" name:"is_backup"`
	Memory                     *int        `json:"memory" name:"memory"`
	Monitor                    interface{} `json:"monitor" name:"monitor"`
	Name                       *string     `json:"name" name:"name"`
	NodeID                     *string     `json:"node_id" name:"node_id"`
	Owner                      *string     `json:"owner" name:"owner"`
	Passphraseless             *string     `json:"passphraseless" name:"passphraseless"`
	PrivateIP                  *string     `json:"private_ip" name:"private_ip"`
	Repl                       *string     `json:"repl" name:"repl"`
	ResourceClass              *int        `json:"resource_class" name:"resource_class"`
	RestartService             interface{} `json:"restart_service" name:"restart_service"`
	RestoreService             interface{} `json:"restore_service" name:"restore_service"`
	Role                       *string     `json:"role" name:"role"`
	RootUserID                 *string     `json:"root_user_id" name:"root_user_id"`
	ScaleInService             interface{} `json:"scale_in_service" name:"scale_in_service"`
	ScaleOutService            interface{} `json:"scale_out_service" name:"scale_out_service"`
	// 这里在SDK中的json处理有报错，这里并不使用，所以先去掉。
	// SecurityGroup              *string     `json:"security_group" name:"security_group"`
	ServerID              *int        `json:"server_id" name:"server_id"`
	ServerIDUpperBound    *int        `json:"server_id_upper_bound" name:"server_id_upper_bound"`
	SingleNodeRepl        *string     `json:"single_node_repl" name:"single_node_repl"`
	StartService          interface{} `json:"start_service" name:"start_service"`
	Status                *string     `json:"status" name:"status"`
	StatusTime            *time.Time  `json:"status_time" name:"status_time" format:"ISO 8601"`
	StopService           interface{} `json:"stop_service" name:"stop_service"`
	StorageSize           *int        `json:"storage_size" name:"storage_size"`
	TransitionStatus      *string     `json:"transition_status" name:"transition_status"`
	UserAccess            *int        `json:"user_access" name:"user_access"`
	VerticalScalingPolicy *string     `json:"vertical_scaling_policy" name:"vertical_scaling_policy"`
	VolumeIDs             *string     `json:"volume_ids" name:"volume_ids"`
	VolumeType            *int        `json:"volume_type" name:"volume_type"`
	VxNetID               *string     `json:"vxnet_id" name:"vxnet_id"`
}

type VxNet struct {
	AvailableIPCount *int       `json:"available_ip_count" name:"available_ip_count"`
	CreateTime       *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description      *string    `json:"description" name:"description"`
	InstanceIDs      []*string  `json:"instance_ids" name:"instance_ids"`
	Owner            *string    `json:"owner" name:"owner"`
	Router           *Router    `json:"router" name:"router"`
	// 目前不需要，所以不从SDK中拉出来了。
	// Tags             []*Tag     `json:"tags" name:"tags"`
	VpcRouterID *string `json:"vpc_router_id" name:"vpc_router_id"`
	VxNetID     *string `json:"vxnet_id" name:"vxnet_id"`
	VxNetName   *string `json:"vxnet_name" name:"vxnet_name"`
	// VxNetType's available values: 0, 1
	VxNetType *int `json:"vxnet_type" name:"vxnet_type"`
}

type Router struct {
	CreateTime  *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description *string    `json:"description" name:"description"`
	DYNIPEnd    *string    `json:"dyn_ip_end" name:"dyn_ip_end"`
	DYNIPStart  *string    `json:"dyn_ip_start" name:"dyn_ip_start"`
	EIP         *EIP       `json:"eip" name:"eip"`
	IPNetwork   *string    `json:"ip_network" name:"ip_network"`
	// IsApplied's available values: 0, 1
	IsApplied  *int    `json:"is_applied" name:"is_applied"`
	ManagerIP  *string `json:"manager_ip" name:"manager_ip"`
	Mode       *int    `json:"mode" name:"mode"`
	PrivateIP  *string `json:"private_ip" name:"private_ip"`
	RouterID   *string `json:"router_id" name:"router_id"`
	RouterName *string `json:"router_name" name:"router_name"`
	// RouterType's available values: 1
	RouterType      *int    `json:"router_type" name:"router_type"`
	SecurityGroupID *string `json:"security_group_id" name:"security_group_id"`
	// Status's available values: pending, active, poweroffed, suspended, deleted, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	// 目前不需要，所以不从SDK中拉出来了。
	// Tags       []*Tag     `json:"tags" name:"tags"`
	// TransitionStatus's available values: creating, updating, suspending, resuming, poweroffing, poweroning, deleting
	TransitionStatus *string  `json:"transition_status" name:"transition_status"`
	VpcNetwork       *string  `json:"vpc_network" name:"vpc_network"`
	VxNets           []*VxNet `json:"vxnets" name:"vxnets"`
}

type EIP struct {
	AlarmStatus   *string `json:"alarm_status" name:"alarm_status"`
	AssociateMode *int    `json:"associate_mode" name:"associate_mode"`
	Bandwidth     *int    `json:"bandwidth" name:"bandwidth"`
	// BillingMode's available values: bandwidth, traffic
	BillingMode *string      `json:"billing_mode" name:"billing_mode"`
	CreateTime  *time.Time   `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description *string      `json:"description" name:"description"`
	EIPAddr     *string      `json:"eip_addr" name:"eip_addr"`
	EIPGroup    *EIPGroup    `json:"eip_group" name:"eip_group"`
	EIPID       *string      `json:"eip_id" name:"eip_id"`
	EIPName     *string      `json:"eip_name" name:"eip_name"`
	ICPCodes    *string      `json:"icp_codes" name:"icp_codes"`
	NeedICP     *int         `json:"need_icp" name:"need_icp"`
	Resource    *EIPResource `json:"resource" name:"resource"`
	// Status's available values: pending, available, associated, suspended, released, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	SubCode    *int       `json:"sub_code" name:"sub_code"`
	// 目前不需要，所以不从SDK中拉出来了。
	// Tags       []*Tag     `json:"tags" name:"tags"`
	// TransitionStatus's available values: associating, dissociating, suspending, resuming, releasing
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
}

type EIPGroup struct {
	EIPGroupID   *string `json:"eip_group_id" name:"eip_group_id"`
	EIPGroupName *string `json:"eip_group_name" name:"eip_group_name"`
}

type EIPResource struct {
	ResourceID   *string `json:"resource_id" name:"resource_id"`
	ResourceName *string `json:"resource_name" name:"resource_name"`
	ResourceType *string `json:"resource_type" name:"resource_type"`
}

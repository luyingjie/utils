package sdk

import (
	"fmt"
	"time"

	verror "utils/os/error"

	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/request"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

var _ fmt.State
var _ time.Time

type UserService struct {
	Config     *config.Config
	Properties *UserServiceProperties
}

type UserServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func NewUser(s *qc.QingCloudService, zone string) (*UserService, error) {
	properties := &UserServiceProperties{
		Zone: &zone,
	}

	return &UserService{Config: s.Config, Properties: properties}, nil
}

// CreateSession : 获取用户Sission
func (s *UserService) CreateSession(i *UserInput) (*UserOutput, error) {
	if i == nil {
		i = &UserInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateSession",
		RequestMethod: "GET",
	}

	x := &UserOutput{}
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

type UserInput struct {
	User   *string `json:"user" name:"user" location:"params"` // Required
	Passwd *string `json:"passwd" name:"passwd" location:"params"`
}

func (v *UserInput) Validate() error {

	if v.User == nil || v.Passwd == nil {
		// return errors.ParameterRequiredError{
		// 	User:   "User",
		// 	Passwd: "Passwd",
		// }
		return verror.New("用户和密码不能为空")
	}

	return nil
}

type UserOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	SK      *string `json:"sk" name:"sk" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// DescribeUsers 查询账户信息
func (s *UserService) DescribeUsers(i *DescribeUserInput) (*DescribeUserOutput, error) {
	if i == nil {
		i = &DescribeUserInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeUsers",
		RequestMethod: "GET",
	}

	x := &DescribeUserOutput{}
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

type DescribeUserInput struct {
	Email *string `json:"email" name:"email" location:"params"`
}

func (v *DescribeUserInput) Validate() error {

	if v.Email == nil {
		return verror.New("用户和密码不能为空")
	}

	return nil
}

type DescribeUserOutput struct {
	TotalCount *int         `json:"total_count" name:"total_count"`
	Action     *string      `json:"action" name:"action" location:"elements"`
	UserSet    []*UserModel `json:"user_set" name:"user_set" location:"elements"`
	RetCode    *int         `json:"ret_code" name:"ret_code" location:"elements"`
}

type UserModel struct {
	UserID      *string `json:"user_id" name:"user_id" location:"elements"`
	Status      *string `json:"status" name:"status" location:"elements"`
	Email       *string `json:"email" name:"email" location:"elements"`
	NotifyEmail *string `json:"notify_email" name:"notify_email" location:"elements"`
	Role        *string `json:"role" name:"role" location:"elements"`
	Privilege   *int    `json:"privilege" name:"privilege" location:"elements"`
	Phone       *string `json:"phone" name:"phone" location:"elements"`
	Remark      *string `json:"remark" name:"remark" location:"elements"`
}

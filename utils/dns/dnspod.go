package dns

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

var client *dnspod.Client

func InitClient(secretId string, secretKey string) {
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ = dnspod.NewClient(credential, "", cpf)
}

func SelectRecord(domain string, subdomain string) *dnspod.RecordListItem {

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewDescribeRecordListRequest()

	request.Domain = common.StringPtr(domain)
	request.Subdomain = common.StringPtr(subdomain)

	// 返回的resp是一个DescribeRecordListResponse的实例，与请求对象对应
	response, err := client.DescribeRecordList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("API client.DescribeRecordList error has returned: %s", err)
		return nil
	}
	if err != nil {
		fmt.Printf("client.DescribeRecordList failed, err: %s", err)
		return nil
	}
	if len(response.Response.RecordList) > 0 {
		return response.Response.RecordList[0]
	}
	return nil
}

func UpdateRecord(domain string, subdomain string, ipv6 string, recordId uint64) {

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewModifyRecordRequest()

	request.Domain = common.StringPtr(domain)
	request.SubDomain = common.StringPtr(subdomain)
	request.RecordType = common.StringPtr("AAAA")
	request.RecordLine = common.StringPtr("默认")
	request.Value = common.StringPtr(ipv6)
	request.RecordId = common.Uint64Ptr(recordId)

	// 返回的resp是一个ModifyRecordResponse的实例，与请求对象对应
	response, err := client.ModifyRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("API client.ModifyRecord error has returned: %s", err)
		return
	}
	if err != nil {
		return
	}
	// 输出json格式的字符串回包
	fmt.Println(response.ToJsonString())
}

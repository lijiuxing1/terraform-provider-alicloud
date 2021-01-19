package emr

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeFlowAgentUser invokes the emr.DescribeFlowAgentUser API synchronously
func (client *Client) DescribeFlowAgentUser(request *DescribeFlowAgentUserRequest) (response *DescribeFlowAgentUserResponse, err error) {
	response = CreateDescribeFlowAgentUserResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeFlowAgentUserWithChan invokes the emr.DescribeFlowAgentUser API asynchronously
func (client *Client) DescribeFlowAgentUserWithChan(request *DescribeFlowAgentUserRequest) (<-chan *DescribeFlowAgentUserResponse, <-chan error) {
	responseChan := make(chan *DescribeFlowAgentUserResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeFlowAgentUser(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeFlowAgentUserWithCallback invokes the emr.DescribeFlowAgentUser API asynchronously
func (client *Client) DescribeFlowAgentUserWithCallback(request *DescribeFlowAgentUserRequest, callback func(response *DescribeFlowAgentUserResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeFlowAgentUserResponse
		var err error
		defer close(result)
		response, err = client.DescribeFlowAgentUser(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeFlowAgentUserRequest is the request struct for api DescribeFlowAgentUser
type DescribeFlowAgentUserRequest struct {
	*requests.RpcRequest
	ResourceOwnerId requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ClusterBizId    string           `position:"Query" name:"ClusterBizId"`
	UserId          string           `position:"Query" name:"UserId"`
}

// DescribeFlowAgentUserResponse is the response struct for api DescribeFlowAgentUser
type DescribeFlowAgentUserResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Data      string `json:"Data" xml:"Data"`
}

// CreateDescribeFlowAgentUserRequest creates a request to invoke DescribeFlowAgentUser API
func CreateDescribeFlowAgentUserRequest() (request *DescribeFlowAgentUserRequest) {
	request = &DescribeFlowAgentUserRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "DescribeFlowAgentUser", "emr", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeFlowAgentUserResponse creates a response to parse from DescribeFlowAgentUser response
func CreateDescribeFlowAgentUserResponse() (response *DescribeFlowAgentUserResponse) {
	response = &DescribeFlowAgentUserResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
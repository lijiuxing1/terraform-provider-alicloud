package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudSimpleApplicationServerInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSimpleApplicationServerInstanceCreate,
		Read:   resourceAlicloudSimpleApplicationServerInstanceRead,
		Update: resourceAlicloudSimpleApplicationServerInstanceUpdate,
		Delete: resourceAlicloudSimpleApplicationServerInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 12, 24, 3, 36, 6}),
			},
			"data_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 16380),
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 12, 24, 3, 36, 6}),
			},
			"plan_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Resetting", "Running", "Stopped", "Upgrading"}, false),
			},
		},
	}
}

func resourceAlicloudSimpleApplicationServerInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateInstances"
	request := make(map[string]interface{})
	var err error
	request["Amount"] = 1
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	if v, ok := d.GetOk("data_disk_size"); ok {
		request["DataDiskSize"] = v
	}
	request["ImageId"] = d.Get("image_id")
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertSimpleApplicationServerInstancePaymentTypeRequest(v.(string))
	}
	request["Period"] = d.Get("period")
	request["PlanId"] = d.Get("plan_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateInstances")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_simple_application_server_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceIds"].([]interface{})[0]))

	swasOpenService := SwasOpenService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, swasOpenService.SimpleApplicationServerInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudSimpleApplicationServerInstanceUpdate(d, meta)
}
func resourceAlicloudSimpleApplicationServerInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	swasOpenService := SwasOpenService{client}
	object, err := swasOpenService.DescribeSimpleApplicationServerInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_simple_application_server_instance swasOpenService.DescribeSimpleApplicationServerInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("image_id", object["ImageId"])
	d.Set("instance_name", object["InstanceName"])
	d.Set("payment_type", convertSimpleApplicationServerInstancePaymentTypeResponse(object["ChargeType"]))
	d.Set("plan_id", object["PlanId"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudSimpleApplicationServerInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	swasOpenService := SwasOpenService{client}
	var err error
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("image_id") {
		update = true
	}
	request["ImageId"] = d.Get("image_id")
	request["RegionId"] = client.RegionId
	if update {
		action := "ResetSystem"
		request["ClientToken"] = buildClientToken("ResetSystem")
		wait := incrementalWait(3*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, swasOpenService.SimpleApplicationServerInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("image_id")
	}
	update = false
	upgradeInstanceReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("plan_id") {
		update = true
	}
	upgradeInstanceReq["PlanId"] = d.Get("plan_id")
	upgradeInstanceReq["RegionId"] = client.RegionId
	if update {
		action := "UpgradeInstance"
		request["ClientToken"] = buildClientToken("UpgradeInstance")
		wait := incrementalWait(3*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, upgradeInstanceReq, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, upgradeInstanceReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, swasOpenService.SimpleApplicationServerInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("plan_id")
	}
	update = false
	updateInstanceAttributeReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	updateInstanceAttributeReq["RegionId"] = client.RegionId
	if d.HasChange("instance_name") {
		update = true
		if v, ok := d.GetOk("instance_name"); ok {
			updateInstanceAttributeReq["InstanceName"] = v
		}
	}
	if d.HasChange("password") {
		update = true
		if v, ok := d.GetOk("password"); ok {
			updateInstanceAttributeReq["Password"] = v
		}
	}
	if update {
		action := "UpdateInstanceAttribute"
		request["ClientToken"] = buildClientToken("UpdateInstanceAttribute")
		wait := incrementalWait(3*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, updateInstanceAttributeReq, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateInstanceAttributeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, swasOpenService.SimpleApplicationServerInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_name")
		d.SetPartial("password")
	}
	if d.HasChange("status") {
		object, err := swasOpenService.DescribeSimpleApplicationServerInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Resetting" {
				request := map[string]interface{}{
					"InstanceId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "RebootInstance"
				request["ClientToken"] = buildClientToken("RebootInstance")
				wait := incrementalWait(3*time.Second, 10*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, swasOpenService.SimpleApplicationServerInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "Running" {
				request := map[string]interface{}{
					"InstanceId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "StartInstance"
				request["ClientToken"] = buildClientToken("StartInstance")
				wait := incrementalWait(3*time.Second, 10*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, swasOpenService.SimpleApplicationServerInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "Stopped" {
				request := map[string]interface{}{
					"InstanceId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "StopInstance"
				request["ClientToken"] = buildClientToken("StopInstance")
				wait := incrementalWait(3*time.Second, 10*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "Throttling.User"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, swasOpenService.SimpleApplicationServerInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudSimpleApplicationServerInstanceRead(d, meta)
}
func resourceAlicloudSimpleApplicationServerInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudSimpleApplicationServerInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
func convertSimpleApplicationServerInstancePaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "PrePaid"
	}
	return source
}
func convertSimpleApplicationServerInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PrePaid":
		return "Subscription"
	}
	return source
}

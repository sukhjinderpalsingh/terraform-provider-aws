// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kafkaconnect

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kafkaconnect"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_mskconnect_custom_plugin")
func ResourceCustomPlugin() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCustomPluginCreate,
		ReadWithoutTimeout:   resourceCustomPluginRead,
		DeleteWithoutTimeout: resourceCustomPluginDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			names.AttrARN: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrContentType: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(kafkaconnect.CustomPluginContentType_Values(), false),
			},
			names.AttrDescription: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"latest_revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"s3": {
							Type:     schema.TypeList,
							MaxItems: 1,
							ForceNew: true,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_arn": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: verify.ValidARN,
									},
									"file_key": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"object_version": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			names.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			names.AttrState: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCustomPluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).KafkaConnectConn(ctx)

	name := d.Get(names.AttrName).(string)
	input := &kafkaconnect.CreateCustomPluginInput{
		ContentType: aws.String(d.Get(names.AttrContentType).(string)),
		Location:    expandCustomPluginLocation(d.Get("location").([]interface{})[0].(map[string]interface{})),
		Name:        aws.String(name),
	}

	if v, ok := d.GetOk(names.AttrDescription); ok {
		input.Description = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating MSK Connect Custom Plugin: %s", input)
	output, err := conn.CreateCustomPluginWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating MSK Connect Custom Plugin (%s): %s", name, err)
	}

	d.SetId(aws.StringValue(output.CustomPluginArn))

	_, err = waitCustomPluginCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for MSK Connect Custom Plugin (%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceCustomPluginRead(ctx, d, meta)...)
}

func resourceCustomPluginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).KafkaConnectConn(ctx)

	plugin, err := FindCustomPluginByARN(ctx, conn, d.Id())

	if tfresource.NotFound(err) && !d.IsNewResource() {
		log.Printf("[WARN] MSK Connect Custom Plugin (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading MSK Connect Custom Plugin (%s): %s", d.Id(), err)
	}

	d.Set(names.AttrARN, plugin.CustomPluginArn)
	d.Set(names.AttrDescription, plugin.Description)
	d.Set(names.AttrName, plugin.Name)
	d.Set(names.AttrState, plugin.CustomPluginState)

	if plugin.LatestRevision != nil {
		d.Set(names.AttrContentType, plugin.LatestRevision.ContentType)
		d.Set("latest_revision", plugin.LatestRevision.Revision)
		if plugin.LatestRevision.Location != nil {
			if err := d.Set("location", []interface{}{flattenCustomPluginLocationDescription(plugin.LatestRevision.Location)}); err != nil {
				return sdkdiag.AppendErrorf(diags, "setting location: %s", err)
			}
		} else {
			d.Set("location", nil)
		}
	} else {
		d.Set(names.AttrContentType, nil)
		d.Set("latest_revision", nil)
		d.Set("location", nil)
	}

	return diags
}

func resourceCustomPluginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).KafkaConnectConn(ctx)

	log.Printf("[DEBUG] Deleting MSK Connect Custom Plugin: %s", d.Id())
	_, err := conn.DeleteCustomPluginWithContext(ctx, &kafkaconnect.DeleteCustomPluginInput{
		CustomPluginArn: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, kafkaconnect.ErrCodeNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting MSK Connect Custom Plugin (%s): %s", d.Id(), err)
	}

	_, err = waitCustomPluginDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for MSK Connect Custom Plugin (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func expandCustomPluginLocation(tfMap map[string]interface{}) *kafkaconnect.CustomPluginLocation {
	if tfMap == nil {
		return nil
	}

	apiObject := &kafkaconnect.CustomPluginLocation{}

	if v, ok := tfMap["s3"].([]interface{}); ok && len(v) > 0 {
		apiObject.S3Location = expandS3Location(v[0].(map[string]interface{}))
	}

	return apiObject
}

func expandS3Location(tfMap map[string]interface{}) *kafkaconnect.S3Location {
	if tfMap == nil {
		return nil
	}

	apiObject := &kafkaconnect.S3Location{}

	if v, ok := tfMap["bucket_arn"].(string); ok && v != "" {
		apiObject.BucketArn = aws.String(v)
	}

	if v, ok := tfMap["file_key"].(string); ok && v != "" {
		apiObject.FileKey = aws.String(v)
	}

	if v, ok := tfMap["object_version"].(string); ok && v != "" {
		apiObject.ObjectVersion = aws.String(v)
	}

	return apiObject
}

func flattenCustomPluginLocationDescription(apiObject *kafkaconnect.CustomPluginLocationDescription) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.S3Location; v != nil {
		tfMap["s3"] = []interface{}{flattenS3LocationDescription(v)}
	}

	return tfMap
}

func flattenS3LocationDescription(apiObject *kafkaconnect.S3LocationDescription) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.BucketArn; v != nil {
		tfMap["bucket_arn"] = aws.StringValue(v)
	}

	if v := apiObject.FileKey; v != nil {
		tfMap["file_key"] = aws.StringValue(v)
	}

	if v := apiObject.ObjectVersion; v != nil {
		tfMap["object_version"] = aws.StringValue(v)
	}

	return tfMap
}

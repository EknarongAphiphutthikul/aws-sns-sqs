package ass_sns

import (
	"context"
	"github.com/EknarongAphiphutthikul/aws-sns-sqs/ass_aws"
	"github.com/EknarongAphiphutthikul/aws-sns-sqs/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"time"
)

var (
	DefaultTimeout               = 5 * time.Second
	DefaultTimeoutPublishMessage = 5 * time.Second
)

func NewClient(ctx context.Context, awsConfig aws.Config) Client {
	if ctx == nil {
		ctx = context.Background()
	}
	return &client{
		ctx:                    ctx,
		client:                 sns.NewFromConfig(awsConfig),
		defaultPublishInputOps: GetDefaultPublishInputOps(),
		defaultTimeout:         DefaultTimeout,
		appendAttributes:       false,
	}
}

func NewClientWithRegion(ctx context.Context, region string) (Client, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	cfg, err := ass_aws.NewConfig(ctx, region)
	return NewClient(ctx, cfg), err
}

type Client interface {
	SetDefaultPublishInputOps(ops PublishInputOps)
	SetDefaultTimeout(timeout time.Duration)
	AppendAttributesEnable()
	Publish(topicArn, message string, ops *PublishInputOps) (*sns.PublishOutput, error)
	GetTopicAttributes(topicArn string) (*sns.GetTopicAttributesOutput, error)
}

type PublishInputOps struct {
	MessageAttributes      map[string]types.MessageAttributeValue
	MessageDeduplicationId *string
	MessageGroupId         *string
	Timeout                time.Duration
}

type client struct {
	ctx                    context.Context
	client                 *sns.Client
	defaultPublishInputOps *PublishInputOps
	defaultTimeout         time.Duration
	appendAttributes       bool
}

func (p *client) SetDefaultPublishInputOps(ops PublishInputOps) {
	p.defaultPublishInputOps = &ops
}

func (p *client) SetDefaultTimeout(timeout time.Duration) {
	p.defaultTimeout = timeout
}

func (p *client) AppendAttributesEnable() {
	p.appendAttributes = true
}

func (p *client) Publish(topicArn, message string, ops *PublishInputOps) (*sns.PublishOutput, error) {
	var finalOps = p.defaultPublishInputOps
	if ops != nil {
		ops.overrideZeroValueBy(p.defaultPublishInputOps, p.appendAttributes)
		finalOps = ops
	}

	var publishInput = generatePublishInput(topicArn, message, finalOps)

	var timeout = p.defaultTimeout
	if finalOps != nil && !utils.IsZeroValue(finalOps.Timeout) {
		timeout = finalOps.Timeout
	}

	var ctx = p.ctx
	if !utils.IsZeroValue(timeout) {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(p.ctx, timeout)
		defer cancelFunc()
	}

	return p.client.Publish(ctx, &publishInput)
}

func (p *client) GetTopicAttributes(topicArn string) (*sns.GetTopicAttributesOutput, error) {
	var ctx = p.ctx
	if !utils.IsZeroValue(p.defaultTimeout) {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(p.ctx, p.defaultTimeout)
		defer cancelFunc()
	}

	return p.client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
		TopicArn: &topicArn,
	})
}

func GetDefaultPublishInputOps() *PublishInputOps {
	return &PublishInputOps{
		Timeout: DefaultTimeoutPublishMessage,
	}
}

func generatePublishInput(topicArn, message string, ops *PublishInputOps) sns.PublishInput {
	msgInput := sns.PublishInput{
		TopicArn: &topicArn,
		Message:  &message,
	}
	if ops == nil {
		return msgInput
	}

	if !utils.IsZeroValue(ops.MessageAttributes) {
		msgInput.MessageAttributes = ops.MessageAttributes
	}
	if !utils.IsZeroValue(ops.MessageGroupId) {
		msgInput.MessageGroupId = ops.MessageGroupId
	}
	if !utils.IsZeroValue(ops.MessageDeduplicationId) {
		msgInput.MessageDeduplicationId = ops.MessageDeduplicationId
	}
	return msgInput
}

func (ops *PublishInputOps) overrideZeroValueBy(refOps *PublishInputOps, appendAttributes bool) {
	if refOps == nil {
		return
	}
	if utils.IsZeroValue(ops.Timeout) {
		ops.Timeout = refOps.Timeout
	}
	if utils.IsZeroValue(ops.MessageAttributes) {
		ops.MessageAttributes = refOps.MessageAttributes
	} else {
		if appendAttributes {
			for k, v := range refOps.MessageAttributes {
				ops.MessageAttributes[k] = v
			}
		}
	}
	if utils.IsZeroValue(ops.MessageGroupId) {
		ops.MessageGroupId = refOps.MessageGroupId
	}
	if utils.IsZeroValue(ops.MessageDeduplicationId) {
		ops.MessageDeduplicationId = refOps.MessageDeduplicationId
	}
}

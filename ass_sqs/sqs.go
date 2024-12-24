package ass_sqs

import (
	"context"
	"github.com/EknarongAphiphutthikul/aws-sns-sqs/ass_aws"
	"github.com/EknarongAphiphutthikul/aws-sns-sqs/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"time"
)

var (
	DefaultTimeout                    = 5 * time.Second
	DefaultTimeoutSendMessage         = 5 * time.Second
	DefaultWaitTimeReceiveMessage     = int32(20)
	DefaultMaxNumberReceiveMessage    = int32(1)
	DefaultMessageSystemAttributeName = []types.MessageSystemAttributeName{types.MessageSystemAttributeNameAll}
	DefaultMessageAttributeNames      = []string{"All"}
)

func NewClient(ctx context.Context, awsConfig aws.Config) Client {
	if ctx == nil {
		ctx = context.Background()
	}
	return &client{
		ctx:                           ctx,
		client:                        sqs.NewFromConfig(awsConfig),
		defaultReceiveMessageInputOps: GetDefaultReceiveMessageInputOps(),
		defaultSendMessageInputOps:    GetDefaultSendMessageInputOps(),
		defaultTimeout:                DefaultTimeout,
		appendAttributes:              false,
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
	SetDefaultSendMessageInputOps(ops SendMessageInputOps)
	SetDefaultReceiveMessageInputOps(ops ReceiveMessageInputOps)
	SetDefaultTimeout(timeout time.Duration)
	AppendAttributesEnable()
	GetQueueURL(queueName, awsAccountId string) (*sqs.GetQueueUrlOutput, error)
	SendMessage(queueUrl, message string, ops *SendMessageInputOps) (*sqs.SendMessageOutput, error)
	GetMessages(queueUrl string, ops *ReceiveMessageInputOps) (*sqs.ReceiveMessageOutput, error)
	GetQueueAttribute(queueUrl string, attributeNames []types.QueueAttributeName) (*sqs.GetQueueAttributesOutput, error)
	DeleteMessage(queueUrl string, message types.Message) (*sqs.DeleteMessageOutput, error)
}

type SendMessageInputOps struct {
	DelaySeconds            int32
	MessageAttributes       map[string]types.MessageAttributeValue
	MessageSystemAttributes map[string]types.MessageSystemAttributeValue
	MessageDeduplicationId  *string
	MessageGroupId          *string
	Timeout                 time.Duration
}

type ReceiveMessageInputOps struct {
	MessageSystemAttributeName []types.MessageSystemAttributeName
	MaxNumberOfMessages        int32
	MessageAttributeNames      []string
	VisibilityTimeout          int32
	WaitTimeSeconds            int32
}

type client struct {
	ctx                           context.Context
	client                        *sqs.Client
	defaultSendMessageInputOps    *SendMessageInputOps
	defaultReceiveMessageInputOps *ReceiveMessageInputOps
	defaultTimeout                time.Duration
	appendAttributes              bool
}

func (p *client) SetDefaultSendMessageInputOps(ops SendMessageInputOps) {
	p.defaultSendMessageInputOps = &ops
}

func (p *client) SetDefaultReceiveMessageInputOps(ops ReceiveMessageInputOps) {
	p.defaultReceiveMessageInputOps = &ops
}

func (p *client) SetDefaultTimeout(timeout time.Duration) {
	p.defaultTimeout = timeout
}

func (p *client) AppendAttributesEnable() {
	p.appendAttributes = true
}

func (p *client) GetQueueURL(queueName, awsAccountId string) (*sqs.GetQueueUrlOutput, error) {
	var ctx = p.ctx
	if !utils.IsZeroValue(p.defaultTimeout) {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(p.ctx, p.defaultTimeout)
		defer cancelFunc()
	}

	return p.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName:              &queueName,
		QueueOwnerAWSAccountId: aws.String(awsAccountId),
	})
}

func (p *client) SendMessage(queueUrl, message string, ops *SendMessageInputOps) (*sqs.SendMessageOutput, error) {
	var finalOps = p.defaultSendMessageInputOps
	if ops != nil {
		ops.overrideZeroValueBy(p.defaultSendMessageInputOps, p.appendAttributes)
		finalOps = ops
	}

	var sendMessageInput = generateSendMessageInput(queueUrl, message, finalOps)

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
	return p.client.SendMessage(ctx, &sendMessageInput)
}

func (p *client) GetMessages(queueUrl string, ops *ReceiveMessageInputOps) (*sqs.ReceiveMessageOutput, error) {
	var finalOps = p.defaultReceiveMessageInputOps
	if ops != nil {
		ops.overrideZeroValueBy(p.defaultReceiveMessageInputOps, p.appendAttributes)
		finalOps = ops
	}

	var receiveMessageInput = generateReceiveMessageInput(queueUrl, finalOps)

	return p.client.ReceiveMessage(p.ctx, &receiveMessageInput)
}

func (p *client) GetQueueAttribute(queueUrl string, attributeNames []types.QueueAttributeName) (*sqs.GetQueueAttributesOutput, error) {
	if len(attributeNames) == 0 {
		attributeNames = []types.QueueAttributeName{types.QueueAttributeNameAll}
	}

	var ctx = p.ctx
	if !utils.IsZeroValue(p.defaultTimeout) {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(p.ctx, p.defaultTimeout)
		defer cancelFunc()
	}

	return p.client.GetQueueAttributes(ctx,
		&sqs.GetQueueAttributesInput{
			QueueUrl:       &queueUrl,
			AttributeNames: attributeNames,
		},
	)
}

func (p *client) DeleteMessage(queueUrl string, message types.Message) (*sqs.DeleteMessageOutput, error) {
	var ctx = p.ctx
	if !utils.IsZeroValue(p.defaultTimeout) {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(p.ctx, p.defaultTimeout)
		defer cancelFunc()
	}

	return p.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: message.ReceiptHandle,
	})
}

func GetDefaultSendMessageInputOps() *SendMessageInputOps {
	return &SendMessageInputOps{
		Timeout: DefaultTimeoutSendMessage,
	}
}

func GetDefaultReceiveMessageInputOps() *ReceiveMessageInputOps {
	return &ReceiveMessageInputOps{
		MessageSystemAttributeName: DefaultMessageSystemAttributeName,
		MessageAttributeNames:      DefaultMessageAttributeNames,
		MaxNumberOfMessages:        DefaultMaxNumberReceiveMessage,
		WaitTimeSeconds:            DefaultWaitTimeReceiveMessage,
	}
}

func generateSendMessageInput(queueUrl, message string, ops *SendMessageInputOps) sqs.SendMessageInput {
	msgInput := sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: &message,
	}
	if ops == nil {
		return msgInput
	}

	if !utils.IsZeroValue(ops.DelaySeconds) {
		msgInput.DelaySeconds = ops.DelaySeconds
	}
	if !utils.IsZeroValue(ops.MessageAttributes) {
		msgInput.MessageAttributes = ops.MessageAttributes
	}
	if !utils.IsZeroValue(ops.MessageSystemAttributes) {
		msgInput.MessageSystemAttributes = ops.MessageSystemAttributes
	}
	if !utils.IsZeroValue(ops.MessageGroupId) {
		msgInput.MessageGroupId = ops.MessageGroupId
	}
	if !utils.IsZeroValue(ops.MessageDeduplicationId) {
		msgInput.MessageDeduplicationId = ops.MessageDeduplicationId
	}
	return msgInput
}

func generateReceiveMessageInput(queueUrl string, ops *ReceiveMessageInputOps) sqs.ReceiveMessageInput {
	msgInput := sqs.ReceiveMessageInput{
		QueueUrl: &queueUrl,
	}
	if ops == nil {
		return msgInput
	}
	if !utils.IsZeroValue(ops.MaxNumberOfMessages) {
		msgInput.MaxNumberOfMessages = ops.MaxNumberOfMessages
	}
	if !utils.IsZeroValue(ops.MessageAttributeNames) {
		msgInput.MessageAttributeNames = ops.MessageAttributeNames
	}
	if !utils.IsZeroValue(ops.MessageSystemAttributeName) {
		msgInput.MessageSystemAttributeNames = ops.MessageSystemAttributeName
	}
	if !utils.IsZeroValue(ops.VisibilityTimeout) {
		msgInput.VisibilityTimeout = ops.VisibilityTimeout
	}
	if !utils.IsZeroValue(ops.WaitTimeSeconds) {
		msgInput.WaitTimeSeconds = ops.WaitTimeSeconds
	}
	return msgInput
}

func (ops *SendMessageInputOps) overrideZeroValueBy(refOps *SendMessageInputOps, appendAttributes bool) {
	if refOps == nil {
		return
	}
	if utils.IsZeroValue(ops.Timeout) {
		ops.Timeout = refOps.Timeout
	}
	if utils.IsZeroValue(ops.DelaySeconds) {
		ops.DelaySeconds = refOps.DelaySeconds
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
	if utils.IsZeroValue(ops.MessageSystemAttributes) {
		ops.MessageSystemAttributes = refOps.MessageSystemAttributes
	} else {
		if appendAttributes {
			for k, v := range refOps.MessageSystemAttributes {
				ops.MessageSystemAttributes[k] = v
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

func (ops *ReceiveMessageInputOps) overrideZeroValueBy(refOps *ReceiveMessageInputOps, appendAttributes bool) {
	if refOps == nil {
		return
	}
	if utils.IsZeroValue(ops.MaxNumberOfMessages) {
		ops.MaxNumberOfMessages = refOps.MaxNumberOfMessages
	}
	if utils.IsZeroValue(ops.MessageAttributeNames) {
		ops.MessageAttributeNames = refOps.MessageAttributeNames
	} else {
		if appendAttributes {
			for _, v := range refOps.MessageAttributeNames {
				ops.MessageAttributeNames = append(ops.MessageAttributeNames, v)
			}
		}
	}
	if utils.IsZeroValue(ops.MessageSystemAttributeName) {
		ops.MessageSystemAttributeName = refOps.MessageSystemAttributeName
	} else {
		if appendAttributes {
			for _, v := range refOps.MessageSystemAttributeName {
				ops.MessageSystemAttributeName = append(ops.MessageSystemAttributeName, v)
			}
		}
	}
	if utils.IsZeroValue(ops.VisibilityTimeout) {
		ops.VisibilityTimeout = refOps.VisibilityTimeout
	}
	if utils.IsZeroValue(ops.WaitTimeSeconds) {
		ops.WaitTimeSeconds = refOps.WaitTimeSeconds
	}
}

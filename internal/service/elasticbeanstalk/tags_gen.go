// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package elasticbeanstalk

import (
	"context"

	"github.com/YakDriver/smarterr"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	awstypes "github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/logging"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types/option"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// listTags lists elasticbeanstalk service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func listTags(ctx context.Context, conn *elasticbeanstalk.Client, identifier string, optFns ...func(*elasticbeanstalk.Options)) (tftags.KeyValueTags, error) {
	input := elasticbeanstalk.ListTagsForResourceInput{
		ResourceArn: aws.String(identifier),
	}

	output, err := conn.ListTagsForResource(ctx, &input, optFns...)

	if err != nil {
		return tftags.New(ctx, nil), smarterr.NewError(err)
	}

	return keyValueTags(ctx, output.ResourceTags), nil
}

// ListTags lists elasticbeanstalk service tags and set them in Context.
// It is called from outside this package.
func (p *servicePackage) ListTags(ctx context.Context, meta any, identifier string) error {
	tags, err := listTags(ctx, meta.(*conns.AWSClient).ElasticBeanstalkClient(ctx), identifier)

	if err != nil {
		return smarterr.NewError(err)
	}

	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(tags)
	}

	return nil
}

// []*SERVICE.Tag handling

// svcTags returns elasticbeanstalk service tags.
func svcTags(tags tftags.KeyValueTags) []awstypes.Tag {
	result := make([]awstypes.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := awstypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// keyValueTags creates tftags.KeyValueTags from elasticbeanstalk service tags.
func keyValueTags(ctx context.Context, tags []awstypes.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.ToString(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns elasticbeanstalk service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []awstypes.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := svcTags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets elasticbeanstalk service tags in Context.
func setTagsOut(ctx context.Context, tags []awstypes.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(keyValueTags(ctx, tags))
	}
}

// updateTags updates elasticbeanstalk service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn *elasticbeanstalk.Client, identifier string, oldTagsMap, newTagsMap any, optFns ...func(*elasticbeanstalk.Options)) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.ElasticBeanstalk)
	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.ElasticBeanstalk)

	// Ensure we do not send empty requests.
	if len(removedTags) == 0 && len(updatedTags) == 0 {
		return nil
	}

	input := elasticbeanstalk.UpdateTagsForResourceInput{
		ResourceArn: aws.String(identifier),
	}

	if len(updatedTags) > 0 {
		input.TagsToAdd = svcTags(updatedTags)
	}

	if len(removedTags) > 0 {
		input.TagsToRemove = removedTags.Keys()
	}

	_, err := conn.UpdateTagsForResource(ctx, &input, optFns...)

	if err != nil {
		return smarterr.NewError(err)
	}

	return nil
}

// UpdateTags updates elasticbeanstalk service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).ElasticBeanstalkClient(ctx), identifier, oldTags, newTags)
}

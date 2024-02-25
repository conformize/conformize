// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
	"github.com/conformize/conformize/serialization/unmarshal/hcl"
)

func ExpectedMockTree() *ds.Node[string, any] {
	root := ds.NewNode[string, any]()

	provider := root.AddChild("provider").AddChild("aws")
	provider.AddAttribute("region", "us-west-2")

	resourceNode := root.AddChild("resource")
	awsInstance := resourceNode.AddChild("aws_instance")
	webInstance := awsInstance.AddChild("web")
	webInstance.AddAttribute("ami", "ami-abc123")
	webInstance.AddAttribute("instance_type", "t2.micro")

	tags := webInstance.AddChild("tags")
	tags.AddAttribute("Name", "webserver")
	tags.AddAttribute("Env", "dev")

	rbd1 := webInstance.AddChild("root_block_device")
	rbd1.AddAttribute("volume_type", "gp2")
	rbd1.AddAttribute("volume_size", 20.0)

	rbd2 := webInstance.AddChild("root_block_device")
	rbd2.AddAttribute("volume_type", "gp2")
	rbd2.AddAttribute("volume_size", 30.0)

	metadata := webInstance.AddChild("metadata_options")
	metadata.AddAttribute("http_endpoint", "enabled")
	metadata.AddAttribute("http_tokens", "required")

	beInstance := awsInstance.AddChild("be")
	beInstance.AddAttribute("ami", "ami-abc123")
	beInstance.AddAttribute("instance_type", "t2.micro")

	beTags := beInstance.AddChild("tags")
	beTags.AddAttribute("Name", "webserver")
	beTags.AddAttribute("Env", "dev")

	beRbd1 := beInstance.AddChild("root_block_device")
	beRbd1.AddAttribute("volume_type", "gp2")
	beRbd1.AddAttribute("volume_size", 20.0)

	beRbd2 := beInstance.AddChild("root_block_device")
	beRbd2.AddAttribute("volume_type", "gp2")
	beRbd2.AddAttribute("volume_size", 30.0)

	beMetadata := beInstance.AddChild("metadata_options")
	beMetadata.AddAttribute("http_endpoint", "enabled")
	beMetadata.AddAttribute("http_tokens", "required")

	bucketInstance := resourceNode.AddChild("aws_s3_bucket")
	bucket := bucketInstance.AddChild("my_bucket")
	bucket.AddAttribute("bucket", "my-unique-bucket-name")
	bucket.AddAttribute("acl", "private")

	versioning := bucket.AddChild("versioning")
	versioning.AddAttribute("enabled", true)

	bucketTags := bucket.AddChild("tags")
	bucketTags.AddAttribute("Name", "MyBucket")
	bucketTags.AddAttribute("Environment", "Dev")

	securityGroup := resourceNode.AddChild("aws_security_group")
	webSG := securityGroup.AddChild("web_sg")
	webSG.AddAttribute("name", "web_sg")
	webSG.AddAttribute("description", "Security group for web server")

	ingress := webSG.AddChild("ingress")
	ingress.AddAttribute("from_port", 80.0)
	ingress.AddAttribute("to_port", 80.0)
	ingress.AddAttribute("protocol", "tcp")

	ingress.AddAttribute("cidr_blocks", []string{"0.0.0.0/0", "10.0.0.0/8"})

	return root
}

func TestHCLFileUnmarshalling(t *testing.T) {
	startTime := time.Now()
	fileSource, _ := serialization.NewFileSource("../../mocks/stack.hcl")
	hclFileUnmarshal := hcl.HclFileUnmarshal{}
	content, err := hclFileUnmarshal.Unmarshal(fileSource)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	duration := time.Since(startTime)
	ms := float64(duration) / float64(time.Millisecond)
	fmt.Printf("execution time: %.2f ms\n", ms)

	expectedTree := ExpectedMockTree()
	if !reflect.DeepEqual(expectedTree, content) {
		fmt.Println("Expected tree:")
		expectedTree.PrintTree()
		fmt.Println("Actual tree:")
		content.PrintTree()
		t.Errorf("Unmarshalled content does not match expected structure")
	}
}

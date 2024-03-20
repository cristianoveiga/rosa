package machinepool

import (
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/spf13/cobra"

	"github.com/openshift/rosa/pkg/helper/features"
	"github.com/openshift/rosa/pkg/rosa"
)

var _ = Describe("Machine pool helper", func() {
	Context("Validates cluster's subnet list isn't empty", func() {
		var r *rosa.Runtime
		var cmd *cobra.Command

		aws := cmv1.NewAWS()
		cluster, err := cmv1.NewCluster().AWS(aws).Build()
		Expect(err).ToNot(HaveOccurred())

		It("should return an error if subnets list is empty", func() {
			_, err := getSecurityGroupsOption(r, cmd, cluster)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(
				"Expected cluster's subnets to contain subnets IDs, but got an empty list"))
		})

		It("Should return an error is subnet is missing the VPC ID", func() {
			subnet := ec2types.Subnet{}

			_, err := getVpcIdFromSubnet(subnet)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(
				"Unexpected situation a VPC ID should have been selected based on chosen subnets"))
		})

		It("Should return VPC ID from the subnet object", func() {
			vpcId := "123"
			subnet := ec2types.Subnet{
				VpcId: &vpcId,
			}

			vpcId, err := getVpcIdFromSubnet(subnet)
			Expect(err).ToNot(HaveOccurred())
			Expect(vpcId).To(Equal("123"))
		})
	})

	Context("It create an AWS node pool builder successfully", func() {
		It("Create AWS node pool with security group IDs when provided", func() {
			instanceType := "123"
			securityGroupIds := []string{"123"}

			awsNpBuilder := createAwsNodePoolBuilder(instanceType, securityGroupIds)
			awsNodePool, err := awsNpBuilder.Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(awsNodePool.AdditionalSecurityGroupIds()).To(Equal(securityGroupIds))
			Expect(awsNodePool.InstanceType()).To(Equal(instanceType))
		})

		It("Create AWS node pool without security group IDs if not provided", func() {
			instanceType := "123"

			awsNpBuilder := createAwsNodePoolBuilder(instanceType, []string{})
			awsNodePool, err := awsNpBuilder.Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(awsNodePool.AdditionalSecurityGroupIds())).To(Equal(0))
			Expect(awsNodePool.InstanceType()).To(Equal(instanceType))
		})
	})

	Context("It validate version is compatible for security groups", func() {
		It("Skips validation if the version isn't provided", func() {
			version := ""
			isCompatible, err := features.IsFeatureSupported(features.AdditionalDay2SecurityGroupsHcpFeature, version)
			Expect(err).ToNot(HaveOccurred())
			Expect(isCompatible).To(BeTrue())
		})
		It("Returns false for 4.14.0", func() {
			version := "4.14.0"
			isCompatible, err := features.IsFeatureSupported(features.AdditionalDay2SecurityGroupsHcpFeature, version)
			Expect(err).ToNot(HaveOccurred())
			Expect(isCompatible).To(BeFalse())
		})
		It("Returns true for 4.15.0", func() {
			version := "4.15.0"
			isCompatible, err := features.IsFeatureSupported(features.AdditionalDay2SecurityGroupsHcpFeature, version)
			Expect(err).ToNot(HaveOccurred())
			Expect(isCompatible).To(BeTrue())
		})
	})
})

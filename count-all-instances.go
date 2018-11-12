package main

import (
	"bufio"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/howeyc/gopass"
	"os"
	"strings"
)

var (
	ecsClient *ecs.Client
	key, secret string
)

func main() {
	key, secret = GetCreds()
	ecsClient = createEcsClient(key, secret, "ap-southeast-1")
	allRegions, err := GetRegions()
	if err != nil {
		panic(err)
	}

	for _, r := range allRegions {
		ecsClient = createEcsClient(key, secret, r)
		allVpcs, err := GetVpcs(r)
		if err != nil {
			panic(err)
		}
		for _, v := range allVpcs {
			instanceCount, err := countInstancesInVpc(v)
			if err != nil {
				panic(err)
			}
			fmt.Println("\nRegion: ", r)
			fmt.Println("VPC: ", v)
			fmt.Println("Instance count: ", instanceCount)
		}
	}

}

func GetCreds() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter Alicloud access key: ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)
	fmt.Printf("Alicloud access secret: ")
	secretByte, err := gopass.GetPasswd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	secret := string(secretByte)
	return key, secret
}

func createEcsClient(key, secret, region string) *ecs.Client {
	ecsClient, err := ecs.NewClientWithAccessKey(
		region,
		key,
		secret)
	if err != nil {
		panic(err)
	}
	return ecsClient
}

func countInstancesInVpc(vpc string) (int, error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.VpcId = vpc
	response, err := ecsClient.DescribeInstances(request)
	if err != nil {
		panic(err)
	}
	return response.TotalCount, nil
}

func GetRegions() ([]string, error) {

	regions := make([]string, 0)

	request := ecs.CreateDescribeRegionsRequest()
	response, err := ecsClient.DescribeRegions(request)
	if err != nil {
		panic(err)
	}

	for _, region := range response.Regions.Region {
		regions = append(regions, region.RegionId)
	}
	return regions, nil
}

func GetVpcs(region string) ([]string, error) {

	vpcs := make([]string, 0)

	request := ecs.CreateDescribeVpcsRequest()
	request.RegionId = region
	response, err := ecsClient.DescribeVpcs(request)
	if err != nil {
		panic(err)
	}

	for _, i := range response.Vpcs.Vpc {
		vpcs = append(vpcs, i.VpcId)
	}
	return vpcs, nil
}

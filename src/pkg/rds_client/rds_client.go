package rds_client

import (
	"context"
	"github.com/amit7itz/goset"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
	log "github.com/sirupsen/logrus"
)

type RdsClient struct {
	region    string
	queueUrl  string
	rdsClient *rds.Client
}

func NewRdsClient(region string, queueUrl string) *RdsClient {
	return &RdsClient{
		region:    region,
		queueUrl:  queueUrl,
		rdsClient: rds.NewFromConfig(common.GetAWSConfig(region)),
	}
}

func (r *RdsClient) getAllDBInstances() []types.DBInstance {
	var instances []types.DBInstance

	log.Infof("getting list of RDS instances in region: %s", r.region)
	paginator := rds.NewDescribeDBInstancesPaginator(r.rdsClient, &rds.DescribeDBInstancesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get a page, %s", err)
		}
		log.Debugf("retrieved %d items", len(page.DBInstances))

		for _, dbInstance := range page.DBInstances {
			log.Debug(*dbInstance.DBInstanceArn)
		}

		instances = append(instances, page.DBInstances...)
	}
	return instances
}

func (r *RdsClient) GetAllDatabases() []string {
	identifiers := goset.NewSet[string]()

	for _, instance := range r.getAllDBInstances() {

		if instance.DBClusterIdentifier != nil {
			log.Debugf("found cluster: %s", *instance.DBClusterIdentifier)
			identifiers.Add(*instance.DBClusterIdentifier)
		} else {
			log.Debugf("found instance: %s", *instance.DBInstanceIdentifier)
			identifiers.Add(*instance.DBInstanceIdentifier)
		}
	}
	log.Infof("found %d RDS clusters and instances", identifiers.Len())

	return identifiers.Items()
}

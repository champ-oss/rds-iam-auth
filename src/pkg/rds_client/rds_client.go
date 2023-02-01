package rds_client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/rds"
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

func (r *RdsClient) GetAllDBClusters() []string {
	var identifiers []string

	log.Infof("getting list of RDS clusters in region: %s", r.region)
	paginator := rds.NewDescribeDBClustersPaginator(r.rdsClient, &rds.DescribeDBClustersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get a page, %s", err)
		}
		log.Debugf("retrieved %d items", len(page.DBClusters))

		for _, dbCluster := range page.DBClusters {
			log.Debug(*dbCluster.DBClusterIdentifier)
			identifiers = append(identifiers, *dbCluster.DBClusterIdentifier)
		}
	}

	log.Infof("found %d RDS clusters", len(identifiers))
	return identifiers
}

func (r *RdsClient) GetAllDBInstances() []string {
	var identifiers []string

	log.Infof("getting list of RDS instances in region: %s", r.region)
	paginator := rds.NewDescribeDBInstancesPaginator(r.rdsClient, &rds.DescribeDBInstancesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get a page, %s", err)
		}
		log.Debugf("retrieved %d items", len(page.DBInstances))

		for _, dbInstance := range page.DBInstances {
			log.Debug(*dbInstance.DBInstanceIdentifier)
			identifiers = append(identifiers, *dbInstance.DBInstanceIdentifier)
		}
	}

	log.Infof("found %d RDS instances", len(identifiers))
	return identifiers
}

package rds_client

import (
	"context"
	"fmt"
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

func (r *RdsClient) GetAllDatabases() []types.DBCluster {
	var dbClusters []types.DBCluster

	log.Infof("getting list of RDS clusters in region: %s", r.region)
	paginator := rds.NewDescribeDBClustersPaginator(r.rdsClient, &rds.DescribeDBClustersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			fmt.Printf("failed to get a page, %w", err)
		}
		log.Debugf("retrieved %d items", len(page.DBClusters))

		for _, dbCluster := range page.DBClusters {
			log.Debug(*dbCluster.DBClusterIdentifier)
		}

		dbClusters = append(dbClusters, page.DBClusters...)
	}

	log.Infof("found %d RDS clusters", len(dbClusters))
	return dbClusters
}

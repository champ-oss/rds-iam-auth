package rds_client

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
)

type RdsClient struct {
	queueUrl  string
	rdsClient *rds.Client
}

func NewRdsClient(region string, queueUrl string) *RdsClient {
	return &RdsClient{
		queueUrl:  queueUrl,
		rdsClient: rds.NewFromConfig(common.GetAWSConfig(region)),
	}
}

func (r *RdsClient) GetAllDatabases() {
	var dbClusters []types.DBCluster

	paginator := rds.NewDescribeDBClustersPaginator(r.rdsClient, &rds.DescribeDBClustersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			fmt.Printf("failed to get a page, %w", err)
		}

		for _, dbCluster := range page.DBClusters {
			fmt.Println(*dbCluster.DBClusterArn)
		}

		dbClusters = append(dbClusters, page.DBClusters...)
	}
}

package rds_client

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	log "github.com/sirupsen/logrus"
)

type RdsClientInterface interface {
	GetAllDBClusters() []string
	GetAllDBInstances() []string
	GetDBCluster(identifier string) (*types.DBCluster, error)
	GetDBInstance(identifier string) (*types.DBInstance, error)
}

type RdsClient struct {
	region    string
	queueUrl  string
	rdsClient *rds.Client
}

func NewRdsClient(config *cfg.Config) *RdsClient {
	return &RdsClient{
		region:    config.AwsRegion,
		queueUrl:  config.QueueUrl,
		rdsClient: rds.NewFromConfig(config.AwsConfig),
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
			log.Infof("found RDS cluster: %s", *dbCluster.DBClusterIdentifier)
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
			if dbInstance.DBClusterIdentifier == nil {
				log.Infof("found RDS instance: %s", *dbInstance.DBInstanceIdentifier)
				identifiers = append(identifiers, *dbInstance.DBInstanceIdentifier)
			}
		}
	}

	log.Infof("found %d RDS instances", len(identifiers))
	return identifiers
}

func (r *RdsClient) GetDBCluster(identifier string) (*types.DBCluster, error) {
	log.Infof("getting information for RDS cluster: %s", identifier)
	output, err := r.rdsClient.DescribeDBClusters(context.TODO(), &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(identifier),
	})
	if err != nil {
		return nil, err
	}
	if len(output.DBClusters) != 1 {
		return nil, fmt.Errorf("unable to find RDS cluster with identifier %s", identifier)
	}
	return &output.DBClusters[0], nil
}

func (r *RdsClient) GetDBInstance(identifier string) (*types.DBInstance, error) {
	log.Infof("getting information for RDS instance: %s", identifier)
	output, err := r.rdsClient.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(identifier),
	})
	if err != nil {
		return nil, err
	}
	if len(output.DBInstances) != 1 {
		return nil, fmt.Errorf("unable to find RDS instance with identifier %s", identifier)
	}
	return &output.DBInstances[0], nil
}

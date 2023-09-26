package interfaces

import (
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
)

type EKSProvider struct{}

func (EKSProvider) CreateHostedCluster(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (EKSProvider) DeleteHostedCluster(cluster *management.Cluster, client *rancher.Client) error {
	//TODO implement me
	panic("implement me")
}

func (EKSProvider) CreateCloudCredential(cluster *management.Cluster, client *rancher.Client) (*management.CloudCredential, error) {
	//TODO implement me
	panic("implement me")
}

func (EKSProvider) UpdateCloudCredential(cluster *management.Cluster, client *rancher.Client) (*management.CloudCredential, error) {
	//TODO implement me
	panic("implement me")
}

func (EKSProvider) DeleteCloudCredential(cluster *management.Cluster, client *rancher.Client) error {
	//TODO implement me
	panic("implement me")
}

func (EKSProvider) UpgradeKubernetesVersion(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (EKSProvider) UpgradeNodePoolKubernetesVersion(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

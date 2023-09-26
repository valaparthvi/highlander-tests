package interfaces

import (
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
)

type AKSProvider struct{}

func (AKSProvider) CreateHostedCluster(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (AKSProvider) DeleteHostedCluster(cluster *management.Cluster, client *rancher.Client) error {
	//TODO implement me
	panic("implement me")
}

func (AKSProvider) CreateCloudCredential(cluster *management.Cluster, client *rancher.Client) (*management.CloudCredential, error) {
	//TODO implement me
	panic("implement me")
}

func (AKSProvider) UpdateCloudCredential(cluster *management.Cluster, client *rancher.Client) (*management.CloudCredential, error) {
	//TODO implement me
	panic("implement me")
}

func (AKSProvider) DeleteCloudCredential(cluster *management.Cluster, client *rancher.Client) error {
	//TODO implement me
	panic("implement me")
}

func (AKSProvider) UpgradeKubernetesVersion(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (AKSProvider) UpgradeNodePoolKubernetesVersion(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

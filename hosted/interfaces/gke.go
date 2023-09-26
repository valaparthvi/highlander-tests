package interfaces

import (
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
)

type GKEProvider struct{}

func (GKEProvider) CreateHostedCluster(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (GKEProvider) DeleteHostedCluster(cluster *management.Cluster, client *rancher.Client) error {
	//TODO implement me
	panic("implement me")
}

func (GKEProvider) CreateCloudCredential(cluster *management.Cluster, client *rancher.Client) (*management.CloudCredential, error) {
	//TODO implement me
	panic("implement me")
}

func (GKEProvider) UpdateCloudCredential(cluster *management.Cluster, client *rancher.Client) (*management.CloudCredential, error) {
	//TODO implement me
	panic("implement me")
}

func (GKEProvider) DeleteCloudCredential(cluster *management.Cluster, client *rancher.Client) error {
	//TODO implement me
	panic("implement me")
}

func (GKEProvider) UpgradeKubernetesVersion(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (GKEProvider) UpgradeNodePoolKubernetesVersion(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

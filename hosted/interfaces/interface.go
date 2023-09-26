package interfaces

import (
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
)

type HostedProvider interface {
	CreateHostedCluster(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error)
	//UpdateHostedCluster(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error)
	DeleteHostedCluster(cluster *management.Cluster, client *rancher.Client) error
	CreateCloudCredential(cluster *management.Cluster, client *rancher.Client) (*management.CloudCredential, error)
	UpdateCloudCredential(cluster *management.Cluster, client *rancher.Client) (*management.CloudCredential, error)
	DeleteCloudCredential(cluster *management.Cluster, client *rancher.Client) error
	UpgradeKubernetesVersion(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error)
	UpgradeNodePoolKubernetesVersion(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error)
}

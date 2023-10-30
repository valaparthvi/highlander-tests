package helpers

import (
	"fmt"
	"os"

	"github.com/epinio/epinio/acceptance/helpers/proc"
	"github.com/pkg/errors"
)

// Create Azure AKS cluster
func CreateClusterAKS(location string, clusterName string, k8sVersion string, nodes string) error {

	fmt.Println("Creating AKS resource group ...")
	out, err := proc.RunW("az", "group", "create", "--location", location, "--resource-group", clusterName)
	if err != nil {
		return errors.Wrap(err, "Failed to create cluster: "+out)
	}

	fmt.Println("Creating AKS cluster ...")
	out, err = proc.RunW("az", "aks", "create", "--resource-group", clusterName, "--kubernetes-version", k8sVersion, "--enable-managed-identity", "--name", clusterName, "--node-count", nodes)
	if err != nil {
		return errors.Wrap(err, "Failed to create cluster: "+out)
	}

	fmt.Println("Created AKS cluster: ", clusterName)

	return nil
}

// Create AWS EKS cluster
func CreateClusterEKS(eks_region string, clusterName string, k8sVersion string, nodes string) error {

	fmt.Println("Creating EKS cluster ...")
	out, err := proc.RunW("eksctl", "create", "cluster", "--region="+eks_region, "--name="+clusterName, "--version="+k8sVersion, "--nodegroup-name", "ranchernodes", "--nodes", nodes, "--managed")
	if err != nil {
		return errors.Wrap(err, "Failed to create cluster: "+out)
	}
	fmt.Println("Created EKS cluster: ", clusterName)

	return nil
}

// Create Google GKE cluster
func CreateClusterGKE(clusterName string) error {
	gke_zone := os.Getenv("GKE_ZONE")

	fmt.Println("Creating GKE cluster ...")
	os.Setenv("USE_GKE_GCLOUD_AUTH_PLUGIN", "true")
	out, err := proc.RunW("gcloud", "container", "clusters", "delete", clusterName, "--zone", gke_zone, "--quiet")
	if err != nil {
		return errors.Wrap(err, "Failed to create cluster: "+out)
	}

	fmt.Println("Created GKE cluster: ", clusterName)

	return nil
}

// Complete cleanup steps for Azure AKS
func DeleteClusterAKS(clusterName string) error {

	fmt.Println("Deleting AKS resource group which will delete cluster too ...")
	out, err := proc.RunW("az", "group", "delete", "--name", clusterName, "--yes")
	if err != nil {
		return errors.Wrap(err, "Failed to delete resource group: "+out)
	}

	fmt.Println("Deleted AKS resource group: ", clusterName)

	return nil
}

// Complete cleanup steps for Amazon EKS
func DeleteClusterEKS(eks_region string, clusterName string) error {

	fmt.Println("Deleting EKS cluster ...")
	out, err := proc.RunW("eksctl", "delete", "cluster", "--region="+eks_region, "--name="+clusterName)
	if err != nil {
		return errors.Wrap(err, "Failed to delete cluster: "+out)
	}

	fmt.Println("Deleted EKS cluster: ", clusterName)

	return nil
}

// Complete cleanup steps for Google GKE
func DeleteClusterGKE(clusterName string) error {
	gke_zone := os.Getenv("GKE_ZONE")

	fmt.Println("Deleting GKE cluster ...")
	os.Setenv("USE_GKE_GCLOUD_AUTH_PLUGIN", "true")
	out, err := proc.RunW("gcloud", "container", "clusters", "delete", clusterName, "--zone", gke_zone, "--quiet")
	if err != nil {
		return errors.Wrap(err, "Failed to delete cluster: "+out)
	}

	fmt.Println("Deleted GKE cluster: ", clusterName)

	return nil
}

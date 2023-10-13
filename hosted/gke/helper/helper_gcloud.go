package helper

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateGKEClusterUsingGCloud(clusterName, location, project, k8sVersion, nodeCount string, labels []string) error {

	args := []string{"container", "--project", project, "clusters", "create", clusterName, "--location", location, "--disk-size=50", "--no-enable-basic-auth", "--no-issue-client-certificate", "--no-enable-autoupgrade", "--metadata=disable-legacy-endpoints=true"}
	if k8sVersion != "" {
		args = append(args, "---cluster-version", k8sVersion)
	}
	if nodeCount != "" {
		args = append(args, "--num-nodes", nodeCount)
	}
	labels = append(labels, "owner=highlander-qa")
	args = append(args, "--labels", strings.Join(labels, ","))

	cmd := exec.Command("gcloud", args...)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func GetGKEClusterKubeConfigUsingGCloud(clusterName, location, project string) (*rest.Config, error) {

	kubeconfigPath, err := os.CreateTemp("", "imported-gke-")
	if err != nil {
		return nil, err
	}
	// TODO: defer to delete the temp file(?)
	os.Setenv("KUBECONFIG", kubeconfigPath.Name())

	args := []string{"container", "clusters", "get-credentials", clusterName, "--location", location, "--project", project}

	cmd := exec.Command("gcloud", args...)
	_, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	kubeconfigData, err := os.ReadFile(kubeconfigPath.Name())
	if err != nil {
		return nil, err
	}
	return clientcmd.RESTConfigFromKubeConfig(kubeconfigData)
}

func DeleteGKEClusterUsingGCloud(clusterName, location, project string, wait bool) error {

	args := []string{"container", "clusters", "delete", clusterName, "--location", location, "--project", project, "--quiet"}
	if !wait {
		args = append(args, "--async")
	}
	cmd := exec.Command("gcloud", args...)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func CheckGKEClusterExistsUsingGCloud(clusterName, location, project string) bool {
	args := []string{"container", "clusters", "list", "--filter", fmt.Sprintf("%v AND status:RUNNING", clusterName), "--location", location, "--project", project, "--quiet"}
	cmd := exec.Command("gcloud", args...)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	if strings.Contains(string(out), clusterName) {
		return true
	}
	return false
}

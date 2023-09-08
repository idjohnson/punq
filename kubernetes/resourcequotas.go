package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	core "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllResourceQuotas(namespaceName string) utils.K8sWorkloadResult {
	result := []core.ResourceQuota{}

	provider := NewKubeProvider()
	rqList, err := provider.ClientSet.CoreV1().ResourceQuotas(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllResourceQuotas ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, rq := range rqList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, rq.ObjectMeta.Namespace) {
			result = append(result, rq)
		}
	}
	return WorkloadResult(result, nil)
}

func GetResourceQuota(namespaceName string, name string) (*core.ResourceQuota, error) {
	provider := NewKubeProvider()
	return provider.ClientSet.CoreV1().ResourceQuotas(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sResourceQuota(data core.ResourceQuota) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().ResourceQuotas(data.Namespace)
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sResourceQuota(data core.ResourceQuota) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().ResourceQuotas(data.Namespace)
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sResourceQuotaBy(namespace string, name string) error {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().ResourceQuotas(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sResourceQuota(namespace string, name string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "resourcequotas", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sResourceQuota(data core.ResourceQuota) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().ResourceQuotas(data.Namespace)
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sResourceQuota() K8sNewWorkload {
	return NewWorkload(
		RES_RESOURCEQUOTAS,
		utils.InitResourceQuotaYaml(),
		"A ResourceQuota is a Kubernetes object that provides constraints that limit aggregate resource consumption per namespace. It can limit the quantity of objects that can be created in a namespace by type, as well as the total amount of compute resources that may be consumed by resources in that namespace. In this example, the quota named 'compute-resources' restricts the namespace to a maximum of 10 pods, request up to 1 CPU, request up to 1Gi of memory, limit up to 2 CPUs, and limit up to 2Gi of memory.")
}

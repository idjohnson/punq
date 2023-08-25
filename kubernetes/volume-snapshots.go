package kubernetes

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	snap "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllVolumeSnapshots(namespace string) K8sWorkloadResult {
	result := []snap.VolumeSnapshot{}

	provider := NewKubeProviderSnapshot()
	volSnapshotsList, err := provider.ClientSet.SnapshotV1().VolumeSnapshots(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllVolumeSnapshots ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	result = append(result, volSnapshotsList.Items...)
	return WorkloadResult(result, nil)
}

func UpdateK8sVolumeSnapshot(data snap.VolumeSnapshot) K8sWorkloadResult {
	return WorkloadResult(nil, fmt.Errorf("UPDATE not available in VolumeSnapshot."))
}

func DeleteK8sVolumeSnapshot(data snap.VolumeSnapshot) K8sWorkloadResult {
	kubeProvider := NewKubeProviderSnapshot()
	client := kubeProvider.ClientSet.SnapshotV1().VolumeSnapshots(data.Namespace)
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sVolumeSnapshot(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "volumesnapshots", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sVolumeSnapshot(data snap.VolumeSnapshot) K8sWorkloadResult {
	kubeProvider := NewKubeProviderSnapshot()
	client := kubeProvider.ClientSet.SnapshotV1().VolumeSnapshots(data.Namespace)
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sVolumeSnapshots() K8sNewWorkload {
	return NewWorkload(
		RES_VOLUMESNAPSHOTS,
		utils.InitVolumeSnapshotYaml(),
		"A VolumeSnapshot in Kubernetes is a representation of a storage volume at a particular point in time. It's part of the Kubernetes storage system and is used for creating backups of data.	This YAML file will create a VolumeSnapshot named 'snapshot-test' from the PersistentVolumeClaim named 'pvc-test'. The snapshot will be taken using the VolumeSnapshotClass named 'snapshot-class'. The VolumeSnapshotClass would typically be defined by your storage provider and would specify the underlying snapshotting technology to use.")
}

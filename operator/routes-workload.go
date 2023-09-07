package operator

import (
	"github.com/mogenius/punq/utils"
	"net/http"

	v1Cert "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	v6Snap "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	v1Apps "k8s.io/api/apps/v1"
	v2Scale "k8s.io/api/autoscaling/v2"
	v1Job "k8s.io/api/batch/v1"
	v1Coordination "k8s.io/api/coordination/v1"
	v1 "k8s.io/api/core/v1"
	v1Networking "k8s.io/api/networking/v1"
	v1Rbac "k8s.io/api/rbac/v1"
	v1Scheduling "k8s.io/api/scheduling/v1"
	v1Storage "k8s.io/api/storage/v1"
	apiExt "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

func InitWorkloadRoutes(router *gin.Engine) {

	workloadRoutes := router.Group("/workload")
	{
		workloadRoutes.GET("/templates", Auth(dtos.USER), allWorkloadTemplates)
		workloadRoutes.GET("/available-resources", Auth(dtos.READER), allKubernetesResources)

		workloadRoutes.GET("/namespace/all", Auth(dtos.USER), allNamespaces)   // QUERY: -
		workloadRoutes.DELETE("/namespace", Auth(dtos.ADMIN), deleteNamespace) // BODY: json-object
		workloadRoutes.POST("/namespace", Auth(dtos.USER), createNamespace)    // BODY: yaml-object

		workloadRoutes.GET("/pod", Auth(dtos.USER), allPods)              // QUERY: namespace
		workloadRoutes.GET("/pod/describe", Auth(dtos.USER), describePod) // QUERY: namespace, name
		workloadRoutes.DELETE("/pod", Auth(dtos.USER), deletePod)         // BODY: json-object
		workloadRoutes.PATCH("/pod", Auth(dtos.USER), patchPod)           // BODY: json-object
		workloadRoutes.POST("/pod", Auth(dtos.USER), createPod)           // BODY: yaml-object

		workloadRoutes.GET("/deployment", Auth(dtos.USER), allDeployments)              // QUERY: namespace
		workloadRoutes.GET("/deployment/describe", Auth(dtos.USER), describeDeployment) // QUERY: namespace, name
		workloadRoutes.DELETE("/deployment", Auth(dtos.USER), deleteDeployment)         // BODY: json-object
		workloadRoutes.PATCH("/deployment", Auth(dtos.USER), patchDeployment)           // BODY: json-object
		workloadRoutes.POST("/deployment", Auth(dtos.USER), createDeployment)           // BODY: yaml-object

		workloadRoutes.GET("/service", Auth(dtos.USER), allServices)              // QUERY: namespace
		workloadRoutes.GET("/service/describe", Auth(dtos.USER), describeService) // QUERY: namespace, name
		workloadRoutes.DELETE("/service", Auth(dtos.USER), deleteService)         // BODY: json-object
		workloadRoutes.PATCH("/service", Auth(dtos.USER), patchService)           // BODY: json-object
		workloadRoutes.POST("/service", Auth(dtos.USER), createService)           // BODY: yaml-object

		workloadRoutes.GET("/ingress", Auth(dtos.USER), allIngresses)             // QUERY: namespace
		workloadRoutes.GET("/ingress/describe", Auth(dtos.USER), describeIngress) // QUERY: namespace, name
		workloadRoutes.DELETE("/ingress", Auth(dtos.USER), deleteIngress)         // BODY: json-object
		workloadRoutes.PATCH("/ingress", Auth(dtos.USER), patchIngress)           // BODY: json-object
		workloadRoutes.POST("/ingress", Auth(dtos.USER), createIngress)           // BODY: yaml-object

		workloadRoutes.GET("/configmap", Auth(dtos.USER), allConfigmaps)              // QUERY: namespace
		workloadRoutes.GET("/configmap/describe", Auth(dtos.USER), describeConfigmap) // QUERY: namespace, name
		workloadRoutes.DELETE("/configmap", Auth(dtos.USER), deleteConfigmap)         // BODY: json-object
		workloadRoutes.PATCH("/configmap", Auth(dtos.USER), patchConfigmap)           // BODY: json-object
		workloadRoutes.POST("/configmap", Auth(dtos.USER), createConfigmap)           // BODY: yaml-object

		workloadRoutes.GET("/secret", Auth(dtos.ADMIN), allSecrets)              // QUERY: namespace
		workloadRoutes.GET("/secret/describe", Auth(dtos.ADMIN), describeSecret) // QUERY: namespace, name
		workloadRoutes.DELETE("/secret", Auth(dtos.ADMIN), deleteSecret)         // BODY: json-object
		workloadRoutes.PATCH("/secret", Auth(dtos.ADMIN), patchSecret)           // BODY: json-object
		workloadRoutes.POST("/secret", Auth(dtos.ADMIN), createSecret)           // BODY: yaml-object

		workloadRoutes.GET("/node", Auth(dtos.USER), allNodes)              // QUERY: -
		workloadRoutes.GET("/node/describe", Auth(dtos.USER), describeNode) // QUERY:  name

		workloadRoutes.GET("/daemon_set", Auth(dtos.USER), allDaemonSets)              // QUERY: namespace
		workloadRoutes.GET("/daemon_set/describe", Auth(dtos.USER), describeDaemonSet) // QUERY: namespace, name
		workloadRoutes.DELETE("/daemon_set", Auth(dtos.USER), deleteDaemonSet)         // BODY: json-object
		workloadRoutes.PATCH("/daemon_set", Auth(dtos.USER), patchDaemonSet)           // BODY: json-object
		workloadRoutes.POST("/daemon_set", Auth(dtos.USER), createDaemonSet)           // BODY: yaml-object

		workloadRoutes.GET("/stateful_set", Auth(dtos.USER), allStatefulSets)              // QUERY: namespace
		workloadRoutes.GET("/stateful_set/describe", Auth(dtos.USER), describeStatefulSet) // QUERY: namespace, name
		workloadRoutes.DELETE("/stateful_set", Auth(dtos.USER), deleteStatefulSet)         // BODY: json-object
		workloadRoutes.PATCH("/stateful_set", Auth(dtos.USER), patchStatefulSet)           // BODY: json-object
		workloadRoutes.POST("/stateful_set", Auth(dtos.USER), createStatefulSet)           // BODY: yaml-object

		workloadRoutes.GET("/job", Auth(dtos.USER), allJobs)              // QUERY: namespace
		workloadRoutes.GET("/job/describe", Auth(dtos.USER), describeJob) // QUERY: namespace, name
		workloadRoutes.DELETE("/job", Auth(dtos.USER), deleteJob)         // BODY: json-object
		workloadRoutes.PATCH("/job", Auth(dtos.USER), patchJob)           // BODY: json-object
		workloadRoutes.POST("/job", Auth(dtos.USER), createJob)           // BODY: yaml-object

		workloadRoutes.GET("/cron_job", Auth(dtos.USER), allCronJobs)              // QUERY: namespace
		workloadRoutes.GET("/cron_job/describe", Auth(dtos.USER), describeCronJob) // QUERY: namespace, name
		workloadRoutes.DELETE("/cron_job", Auth(dtos.USER), deleteCronJob)         // BODY: json-object
		workloadRoutes.PATCH("/cron_job", Auth(dtos.USER), patchCronJob)           // BODY: json-object
		workloadRoutes.POST("/cron_job", Auth(dtos.USER), createCronJob)           // BODY: yaml-object

		workloadRoutes.GET("/replica_set", Auth(dtos.USER), allReplicasets)              // QUERY: namespace
		workloadRoutes.GET("/replica_set/describe", Auth(dtos.USER), describeReplicaset) // QUERY: namespace, name
		workloadRoutes.DELETE("/replica_set", Auth(dtos.USER), deleteReplicaset)         // BODY: json-object
		workloadRoutes.PATCH("/replica_set", Auth(dtos.USER), patchReplicaset)           // BODY: json-object
		workloadRoutes.POST("/replica_set", Auth(dtos.USER), createReplicaset)           // BODY: yaml-object

		workloadRoutes.GET("/persistent_volume", Auth(dtos.ADMIN), allPersistentVolumes)              // QUERY: -
		workloadRoutes.GET("/persistent_volume/describe", Auth(dtos.ADMIN), describePersistentVolume) // QUERY: name
		workloadRoutes.DELETE("/persistent_volume", Auth(dtos.ADMIN), deletePersistentVolume)         // BODY: json-object
		workloadRoutes.PATCH("/persistent_volume", Auth(dtos.ADMIN), patchPersistentVolume)           // BODY: json-object
		workloadRoutes.POST("/persistent_volume", Auth(dtos.ADMIN), createPersistentVolume)           // BODY: yaml-object

		workloadRoutes.GET("/persistent_volume_claim", Auth(dtos.USER), allPersistentVolumeClaims)              // QUERY: namespace
		workloadRoutes.GET("/persistent_volume_claim/describe", Auth(dtos.USER), describePersistentVolumeClaim) // QUERY: namespace, name
		workloadRoutes.DELETE("/persistent_volume_claim", Auth(dtos.ADMIN), deletePersistentVolumeClaim)        // BODY: json-object
		workloadRoutes.PATCH("/persistent_volume_claim", Auth(dtos.ADMIN), patchPersistentVolumeClaim)          // BODY: json-object
		workloadRoutes.POST("/persistent_volume_claim", Auth(dtos.ADMIN), createPersistentVolumeClaim)          // BODY: yaml-object

		workloadRoutes.GET("/horizontal_pod_autoscaler", Auth(dtos.USER), allHpas)              // QUERY: namespace
		workloadRoutes.GET("/horizontal_pod_autoscaler/describe", Auth(dtos.USER), describeHpa) // QUERY: namespace, name
		workloadRoutes.DELETE("/horizontal_pod_autoscaler", Auth(dtos.ADMIN), deleteHpa)        // BODY: json-object
		workloadRoutes.PATCH("/horizontal_pod_autoscaler", Auth(dtos.ADMIN), patchHpa)          // BODY: json-object
		workloadRoutes.POST("/horizontal_pod_autoscaler", Auth(dtos.ADMIN), createHpa)          // BODY: yaml-object

		workloadRoutes.GET("/event", Auth(dtos.USER), allEvents)              // QUERY: namespace
		workloadRoutes.GET("/event/describe", Auth(dtos.USER), describeEvent) // QUERY: namespace, name

		workloadRoutes.GET("/certificate", Auth(dtos.USER), allCertificates)              // QUERY: namespace
		workloadRoutes.GET("/certificate/describe", Auth(dtos.USER), describeCertificate) // QUERY: namespace, name
		workloadRoutes.DELETE("/certificate", Auth(dtos.USER), deleteCertificate)         // BODY: json-object
		workloadRoutes.PATCH("/certificate", Auth(dtos.USER), patchCertificate)           // BODY: json-object
		workloadRoutes.POST("/certificate", Auth(dtos.USER), createCertificate)           // BODY: yaml-object

		workloadRoutes.GET("/certificaterequest", Auth(dtos.USER), allCertificateRequests)              // QUERY: namespace
		workloadRoutes.GET("/certificaterequest/describe", Auth(dtos.USER), describeCertificateRequest) // QUERY: namespace, name
		workloadRoutes.DELETE("/certificaterequest", Auth(dtos.USER), deleteCertificateRequest)         // BODY: json-object
		workloadRoutes.PATCH("/certificaterequest", Auth(dtos.USER), patchCertificateRequest)           // BODY: json-object
		workloadRoutes.POST("/certificaterequest", Auth(dtos.USER), createCertificateRequest)           // BODY: yaml-object

		workloadRoutes.GET("/orders", Auth(dtos.USER), allOrders)              // QUERY: namespace
		workloadRoutes.GET("/orders/describe", Auth(dtos.USER), describeOrder) // QUERY: namespace, name
		workloadRoutes.DELETE("/orders", Auth(dtos.USER), deleteOrder)         // BODY: json-object
		workloadRoutes.PATCH("/orders", Auth(dtos.USER), patchOrder)           // BODY: json-object
		workloadRoutes.POST("/orders", Auth(dtos.USER), createOrder)           // BODY: yaml-object

		workloadRoutes.GET("/issuer", Auth(dtos.USER), allIssuers)              // QUERY: namespace
		workloadRoutes.GET("/issuer/describe", Auth(dtos.USER), describeIssuer) // QUERY: namespace, name
		workloadRoutes.DELETE("/issuer", Auth(dtos.USER), deleteIssuer)         // BODY: json-object
		workloadRoutes.PATCH("/issuer", Auth(dtos.USER), patchIssuer)           // BODY: json-object
		workloadRoutes.POST("/issuer", Auth(dtos.USER), createIssuer)           // BODY: yaml-object

		workloadRoutes.GET("/clusterissuer", Auth(dtos.ADMIN), allClusterIssuers)              // QUERY: -
		workloadRoutes.GET("/clusterissuer/describe", Auth(dtos.ADMIN), describeClusterIssuer) // QUERY: name
		workloadRoutes.DELETE("/clusterissuer", Auth(dtos.ADMIN), deleteClusterIssuer)         // BODY: json-object
		workloadRoutes.PATCH("/clusterissuer", Auth(dtos.ADMIN), patchClusterIssuer)           // BODY: json-object
		workloadRoutes.POST("/clusterissuer", Auth(dtos.ADMIN), createClusterIssuer)           // BODY: yaml-object

		workloadRoutes.GET("/service_account", Auth(dtos.ADMIN), allServiceAccounts)              // QUERY: namespace
		workloadRoutes.GET("/service_account/describe", Auth(dtos.ADMIN), describeServiceAccount) // QUERY: namespace, name
		workloadRoutes.DELETE("/service_account", Auth(dtos.ADMIN), deleteServiceAccount)         // BODY: json-object
		workloadRoutes.PATCH("/service_account", Auth(dtos.ADMIN), patchServiceAccount)           // BODY: json-object
		workloadRoutes.POST("/service_account", Auth(dtos.ADMIN), createServiceAccount)           // BODY: yaml-object

		workloadRoutes.GET("/role", Auth(dtos.USER), allRoles)              // QUERY: namespace
		workloadRoutes.GET("/role/describe", Auth(dtos.USER), describeRole) // QUERY: namespace, name
		workloadRoutes.DELETE("/role", Auth(dtos.ADMIN), deleteRole)        // BODY: json-object
		workloadRoutes.PATCH("/role", Auth(dtos.ADMIN), patchRole)          // BODY: json-object
		workloadRoutes.POST("/role", Auth(dtos.ADMIN), createRole)          // BODY: yaml-object

		workloadRoutes.GET("/role_binding", Auth(dtos.USER), allRoleBindings)              // QUERY: namespace
		workloadRoutes.GET("/role_binding/describe", Auth(dtos.USER), describeRoleBinding) // QUERY: namespace, name
		workloadRoutes.DELETE("/role_binding", Auth(dtos.ADMIN), deleteRoleBinding)        // BODY: json-object
		workloadRoutes.PATCH("/role_binding", Auth(dtos.ADMIN), patchRoleBinding)          // BODY: json-object
		workloadRoutes.POST("/role_binding", Auth(dtos.ADMIN), createRoleBinding)          // BODY: yaml-object

		workloadRoutes.GET("/cluster_role", Auth(dtos.ADMIN), allClusterRoles)              // QUERY: -
		workloadRoutes.GET("/cluster_role/describe", Auth(dtos.ADMIN), describeClusterRole) // QUERY: name
		workloadRoutes.DELETE("/cluster_role", Auth(dtos.ADMIN), deleteClusterRole)         // BODY: json-object
		workloadRoutes.PATCH("/cluster_role", Auth(dtos.ADMIN), patchClusterRole)           // BODY: json-object
		workloadRoutes.POST("/cluster_role", Auth(dtos.ADMIN), createClusterRole)           // BODY: yaml-object

		workloadRoutes.GET("/cluster_role_binding", Auth(dtos.ADMIN), allClusterRoleBindings)              // QUERY: -
		workloadRoutes.GET("/cluster_role_binding/describe", Auth(dtos.ADMIN), describeClusterRoleBinding) // QUERY: name
		workloadRoutes.DELETE("/cluster_role_binding", Auth(dtos.ADMIN), deleteClusterRoleBinding)         // BODY: json-object
		workloadRoutes.PATCH("/cluster_role_binding", Auth(dtos.ADMIN), patchClusterRoleBinding)           // BODY: json-object
		workloadRoutes.POST("/cluster_role_binding", Auth(dtos.ADMIN), createClusterRoleBinding)           // BODY: yaml-object

		workloadRoutes.GET("/volume_attachment", Auth(dtos.ADMIN), allVolumeAttachments)              // QUERY: -
		workloadRoutes.GET("/volume_attachment/describe", Auth(dtos.ADMIN), describeVolumeAttachment) // QUERY: name
		workloadRoutes.DELETE("/volume_attachment", Auth(dtos.ADMIN), deleteVolumeAttachment)         // BODY: json-object
		workloadRoutes.PATCH("/volume_attachment", Auth(dtos.ADMIN), patchVolumeAttachment)           // BODY: json-object
		workloadRoutes.POST("/volume_attachment", Auth(dtos.ADMIN), createVolumeAttachment)           // BODY: yaml-object

		workloadRoutes.GET("/network_policy", Auth(dtos.USER), allNetworkPolicies)             // QUERY: namespace
		workloadRoutes.GET("/network_policy/describe", Auth(dtos.USER), describeNetworkPolicy) // QUERY: namespace, name
		workloadRoutes.DELETE("/network_policy", Auth(dtos.ADMIN), deleteNetworkPolicy)        // BODY: json-object
		workloadRoutes.PATCH("/network_policy", Auth(dtos.ADMIN), patchNetworkPolicy)          // BODY: json-object
		workloadRoutes.POST("/network_policy", Auth(dtos.ADMIN), createNetworkPolicy)          // BODY: yaml-object

		workloadRoutes.GET("/storageclass", Auth(dtos.USER), allStorageClasses)             // QUERY: namespace
		workloadRoutes.GET("/storageclass/describe", Auth(dtos.USER), describeStorageClass) // QUERY: namespace, name
		workloadRoutes.DELETE("/storageclass", Auth(dtos.ADMIN), deleteStorageClass)        // BODY: json-object
		workloadRoutes.PATCH("/storageclass", Auth(dtos.ADMIN), patchStorageClass)          // BODY: json-object
		workloadRoutes.POST("/storageclass", Auth(dtos.ADMIN), createStorageClass)          // BODY: yaml-object

		workloadRoutes.GET("/crds", Auth(dtos.ADMIN), allCrds)              // QUERY: -
		workloadRoutes.GET("/crds/describe", Auth(dtos.ADMIN), describeCrd) // QUERY: name
		workloadRoutes.DELETE("/crds", Auth(dtos.ADMIN), deleteCrd)         // BODY: json-object
		workloadRoutes.PATCH("/crds", Auth(dtos.ADMIN), patchCrd)           // BODY: json-object
		workloadRoutes.POST("/crds", Auth(dtos.ADMIN), createCrd)           // BODY: yaml-object

		workloadRoutes.GET("/endpoints", Auth(dtos.USER), allEndpoints)              // QUERY: namespace
		workloadRoutes.GET("/endpoints/describe", Auth(dtos.USER), describeEndpoint) // QUERY: namespace, name
		workloadRoutes.DELETE("/endpoints", Auth(dtos.USER), deleteEndpoint)         // BODY: json-object
		workloadRoutes.PATCH("/endpoints", Auth(dtos.USER), patchEndpoint)           // BODY: json-object
		workloadRoutes.POST("/endpoints", Auth(dtos.USER), createEndpoint)           // BODY: yaml-object

		workloadRoutes.GET("/leases", Auth(dtos.USER), allLeases)              // QUERY: namespace
		workloadRoutes.GET("/leases/describe", Auth(dtos.USER), describeLease) // QUERY: namespace, name
		workloadRoutes.DELETE("/leases", Auth(dtos.USER), deleteLease)         // BODY: json-object
		workloadRoutes.PATCH("/leases", Auth(dtos.USER), patchLease)           // BODY: json-object
		workloadRoutes.POST("/leases", Auth(dtos.USER), createLease)           // BODY: yaml-object

		workloadRoutes.GET("/priorityclasses", Auth(dtos.ADMIN), allPriorityClasses)             // QUERY: -
		workloadRoutes.GET("/priorityclasses/describe", Auth(dtos.ADMIN), describePriorityClass) // QUERY: name
		workloadRoutes.DELETE("/priorityclasses", Auth(dtos.ADMIN), deletePriorityClass)         // BODY: json-object
		workloadRoutes.PATCH("/priorityclasses", Auth(dtos.ADMIN), patchPriorityClass)           // BODY: json-object
		workloadRoutes.POST("/priorityclasses", Auth(dtos.ADMIN), createPriorityClass)           // BODY: yaml-object

		workloadRoutes.GET("/volumesnapshots", Auth(dtos.USER), allVolumeSnapshots)              // QUERY: namespace
		workloadRoutes.GET("/volumesnapshots/describe", Auth(dtos.USER), describeVolumeSnapshot) // QUERY: namespace, name
		workloadRoutes.DELETE("/volumesnapshots", Auth(dtos.USER), deleteVolumeSnapshot)         // BODY: json-object
		workloadRoutes.PATCH("/volumesnapshots", Auth(dtos.USER), patchVolumeSnapshot)           // BODY: json-object
		workloadRoutes.POST("/volumesnapshots", Auth(dtos.USER), createVolumeSnapshot)           // BODY: yaml-object

		workloadRoutes.GET("/resourcequota", Auth(dtos.ADMIN), allResourceQuotas)              // QUERY: namespace
		workloadRoutes.GET("/resourcequota/describe", Auth(dtos.ADMIN), describeResourceQuota) // QUERY: namespace, name
		workloadRoutes.DELETE("/resourcequota", Auth(dtos.ADMIN), deleteResourceQuota)         // BODY: json-object
		workloadRoutes.PATCH("/resourcequota", Auth(dtos.ADMIN), patchResourceQuota)           // BODY: json-object
		workloadRoutes.POST("/resourcequota", Auth(dtos.ADMIN), createResourceQuota)           // BODY: yaml-object
	}
}

// GENERAL
// @Tags General
// @Produce json
// @Success 200 {array} kubernetes.K8sNewWorkload
// @Router /workload/templates [get]
// @Security Bearer
func allWorkloadTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ListCreateTemplates())
}

// @Tags General
// @Produce json
// @Success 200 {array} string
// @Router /workload/available-resources [get]
// @Security Bearer
func allKubernetesResources(c *gin.Context) {
	user, err := CheckUserAuthorization(c)
	if err != nil || user == nil {
		utils.MalformedMessage(c, "User not found.")
		return
	}
	c.JSON(http.StatusOK, kubernetes.WorkloadsForAccesslevel(user.AccessLevel))
}

// NAMESPACES
// @Tags Workloads
// @Produce json
// @Success 200 {array} v1.Namespace
// @Router /workload/namespace/all [get]
// @Security Bearer
func allNamespaces(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ListK8sNamespaces(""))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/namespace [post]
// @Security Bearer
func createNamespace(c *gin.Context) {
	var data v1.Namespace
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, kubernetes.CreateK8sNamespace(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/namespace [delete]
// @Security Bearer
func deleteNamespace(c *gin.Context) {
	var data v1.Namespace
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, kubernetes.DeleteK8sNamespace(data))
}

// PODS
// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Pod
// @Router /workload/pod [get]
// @Security Bearer
func allPods(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sPods(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/pod/describe [get]
// @Security Bearer
// @Param namespace query string false  "namespace"
// @Param name query string false  "resource name"
func describePod(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPod(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/pod [delete]
// @Security Bearer
func deletePod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sPod(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/pod [patch]
// @Security Bearer
func patchPod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPod(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/pod [post]
// @Security Bearer
func createPod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPod(data))
}

// DEPLOYMENTS
// @Tags Workloads
// @Produce json
// @Success 200 {array} v1Apps.Deployment
// @Router /workload/deployment [get]
// @Security Bearer
// @Param namespace query string false  "namespace"
func allDeployments(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sDeployments(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/deployment/describe [get]
// @Security Bearer
// @Param namespace query string false  "namespace"
// @Param name query string false  "resource name"
func describeDeployment(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sDeployment(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Deployment
// @Router /workload/deployment [delete]
// @Security Bearer
func deleteDeployment(c *gin.Context) {
	var data v1Apps.Deployment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sDeployment(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Deployment
// @Router /workload/deployment [patch]
// @Security Bearer
func patchDeployment(c *gin.Context) {
	var data v1Apps.Deployment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sDeployment(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Deployment
// @Router /workload/deployment [post]
// @Security Bearer
func createDeployment(c *gin.Context) {
	var data v1Apps.Deployment
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sDeployment(data))
}

// SERVICES
// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Service
// @Router /workload/service/describe [get]
// @Security Bearer
// @Param namespace query string false  "namespace"
func allServices(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sServices(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/service/describe [get]
// @Security Bearer
// @Param namespace query string false  "namespace"
// @Param name query string false  "resource name"
func describeService(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sService(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Service
// @Router /workload/service [delete]
// @Security Bearer
func deleteService(c *gin.Context) {
	var data v1.Service
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sService(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Service
// @Router /workload/service [patch]
// @Security Bearer
func patchService(c *gin.Context) {
	var data v1.Service
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sService(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Service
// @Router /workload/service [post]
// @Security Bearer
func createService(c *gin.Context) {
	var data v1.Service
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sService(data))
}

// INGRESSES
// @Tags Workloads
// @Produce json
// @Success 200 {object} v1Networking.Ingress
// @Router /workload/ingress/describe [get]
// @Security Bearer
// @Param namespace query string false  "namespace"
func allIngresses(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sIngresses(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/ingress/describe [get]
// @Security Bearer
// @Param namespace query string false  "namespace"
// @Param name query string false  "resource name"
func describeIngress(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sIngress(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1Networking.Ingress
// @Router /workload/ingress [delete]
// @Security Bearer
func deleteIngress(c *gin.Context) {
	var data v1Networking.Ingress
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sIngress(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1Networking.Ingress
// @Router /workload/ingress [patch]
// @Security Bearer
func patchIngress(c *gin.Context) {
	var data v1Networking.Ingress
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sIngress(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1Networking.Ingress
// @Router /workload/ingress [post]
// @Security Bearer
func createIngress(c *gin.Context) {
	var data v1Networking.Ingress
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sIngress(data))
}

// CONFIGMAPS
// @Tags Workloads
// @Produce json
// @Success 200 {array} v1.ConfigMap
// @Router /workload/configmap [get]
// @Security Bearer
// @Param namespace query string false "namespace"
func allConfigmaps(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sConfigmaps(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/configmap/describe [get]
// @Security Bearer
// @Param namespace query string false "namespace"
// @Param name query string false "resource name"
func describeConfigmap(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sConfigmap(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.ConfigMap
// @Router /workload/configmap [delete]
// @Security Bearer
func deleteConfigmap(c *gin.Context) {
	var data v1.ConfigMap
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sConfigmap(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.ConfigMap
// @Router /workload/configmap [patch]
// @Security Bearer
func patchConfigmap(c *gin.Context) {
	var data v1.ConfigMap
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sConfigMap(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.ConfigMap
// @Router /workload/configmap [post]
// @Security Bearer
func createConfigmap(c *gin.Context) {
	var data v1.ConfigMap
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sConfigMap(data))
}

// SECRETS
// @Tags Workloads
// @Produce json
// @Success 200 {array} v1.Secret
// @Router /workload/secret [get]
// @Security Bearer
// @Param namespace query string false  "namespace"
func allSecrets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sSecrets(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/secret/describe [get]
// @Security Bearer
// @Param namespace query string false  "namespace"
// @Param name query string false  "resource name"
func describeSecret(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sSecret(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Secret
// @Router /workload/secret [delete]
// @Security Bearer
func deleteSecret(c *gin.Context) {
	var data v1.Secret
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sSecret(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Secret
// @Router /workload/secret [patch]
// @Security Bearer
func patchSecret(c *gin.Context) {
	var data v1.Secret
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sSecret(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1.Secret
// @Router /workload/secret [post]
// @Security Bearer
func createSecret(c *gin.Context) {
	var data v1.Secret
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sSecret(data))
}

// NODES
// @Tags Workloads
// @Produce json
// @Success 200 {array} utils.HttpResult
// @Router /workload/node/all [get]
// @Security Bearer
func allNodes(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.ListK8sNodes())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.HttpResult
// @Router /workload/node/describe [get]
// @Security Bearer
// @Param name query string false  "resource name"
func describeNode(c *gin.Context) {
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sNode(name))
}

// DAEMONSETS
// @Tags Workloads
// @Produce json
// @Success 200 {array} v1Apps.DaemonSet
// @Router /workload/daemonset [get]
// @Security Bearer
// @Param namespace query string false "namespace"
func allDaemonSets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sDaemonsets(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1Apps.DaemonSet
// @Router /workload/daemonset/describe [get]
// @Security Bearer
// @Param namespace query string false "namespace"
func describeDaemonSet(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sDaemonSet(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1Apps.DaemonSet
// @Router /workload/daemonset [delete]
// @Security Bearer
func deleteDaemonSet(c *gin.Context) {
	var data v1Apps.DaemonSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sDaemonSet(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} v1Apps.DaemonSet
// @Router /workload/daemonset [patch]
// @Security Bearer
func patchDaemonSet(c *gin.Context) {
	var data v1Apps.DaemonSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sDaemonSet(data))
}
func createDaemonSet(c *gin.Context) {
	var data v1Apps.DaemonSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sDaemonSet(data))
}

// STATEFULSETS
func allStatefulSets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllStatefulSets(namespace))
}
func describeStatefulSet(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sStatefulset(namespace, name))
}
func deleteStatefulSet(c *gin.Context) {
	var data v1Apps.StatefulSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sStatefulset(data))
}
func patchStatefulSet(c *gin.Context) {
	var data v1Apps.StatefulSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sStatefulset(data))
}
func createStatefulSet(c *gin.Context) {
	var data v1Apps.StatefulSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sStatefulset(data))
}

// JOBS
func allJobs(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllJobs(namespace))
}
func describeJob(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sJob(namespace, name))
}
func deleteJob(c *gin.Context) {
	var data v1Job.Job
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sJob(data))
}
func patchJob(c *gin.Context) {
	var data v1Job.Job
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sJob(data))
}
func createJob(c *gin.Context) {
	var data v1Job.Job
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sJob(data))
}

// CRONJOBS
func allCronJobs(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllCronjobs(namespace))
}
func describeCronJob(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCronJob(namespace, name))
}
func deleteCronJob(c *gin.Context) {
	var data v1Job.CronJob
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sCronJob(data))
}
func patchCronJob(c *gin.Context) {
	var data v1Job.CronJob
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCronJob(data))
}
func createCronJob(c *gin.Context) {
	var data v1Job.CronJob
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCronJob(data))
}

// REPLICASETS
func allReplicasets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sReplicasets(namespace))
}
func describeReplicaset(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sReplicaset(namespace, name))
}
func deleteReplicaset(c *gin.Context) {
	var data v1Apps.ReplicaSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sReplicaset(data))
}
func patchReplicaset(c *gin.Context) {
	var data v1Apps.ReplicaSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sReplicaset(data))
}
func createReplicaset(c *gin.Context) {
	var data v1Apps.ReplicaSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sReplicaSet(data))
}

// PERSISTENTVOLUMES
func allPersistentVolumes(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllPersistentVolumes())
}
func describePersistentVolume(c *gin.Context) {
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPersistentVolume(name))
}
func deletePersistentVolume(c *gin.Context) {
	var data v1.PersistentVolume
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sPersistentVolume(data))
}
func patchPersistentVolume(c *gin.Context) {
	var data v1.PersistentVolume
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPersistentVolume(data))
}
func createPersistentVolume(c *gin.Context) {
	var data v1.PersistentVolume
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPersistentVolume(data))
}

// PERSISTENTVOLUMECLAIMS
func allPersistentVolumeClaims(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sPersistentVolumeClaims(namespace))
}
func describePersistentVolumeClaim(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPersistentVolumeClaim(namespace, name))
}
func deletePersistentVolumeClaim(c *gin.Context) {
	var data v1.PersistentVolumeClaim
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sPersistentVolumeClaim(data))
}
func patchPersistentVolumeClaim(c *gin.Context) {
	var data v1.PersistentVolumeClaim
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPersistentVolumeClaim(data))
}
func createPersistentVolumeClaim(c *gin.Context) {
	var data v1.PersistentVolumeClaim
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPersistentVolumeClaim(data))
}

// HPA
func allHpas(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllHpas(namespace))
}
func describeHpa(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sHpa(namespace, name))
}
func deleteHpa(c *gin.Context) {
	var data v2Scale.HorizontalPodAutoscaler
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sHpa(data))
}
func patchHpa(c *gin.Context) {
	var data v2Scale.HorizontalPodAutoscaler
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sHpa(data))
}
func createHpa(c *gin.Context) {
	var data v2Scale.HorizontalPodAutoscaler
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sHpa(data))
}

// EVENTS
func allEvents(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllEvents(namespace))
}
func describeEvent(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sEvent(namespace, name))
}

// CERTIFICATES
func allCertificates(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sCertificates(namespace))
}
func describeCertificate(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCertificate(namespace, name))
}
func deleteCertificate(c *gin.Context) {
	var data cmapi.Certificate
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sCertificate(data))
}
func patchCertificate(c *gin.Context) {
	var data cmapi.Certificate
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCertificate(data))
}
func createCertificate(c *gin.Context) {
	var data cmapi.Certificate
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCertificate(data))
}

// CERTIFICATEREQUESTS
func allCertificateRequests(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllCertificateSigningRequests(namespace))
}
func describeCertificateRequest(c *gin.Context) {
	name := c.Query("name")
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCertificateSigningRequest(namespace, name))
}
func deleteCertificateRequest(c *gin.Context) {
	var data cmapi.CertificateRequest
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sCertificateSigningRequest(data))
}
func patchCertificateRequest(c *gin.Context) {
	var data cmapi.CertificateRequest
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCertificateSigningRequest(data))
}
func createCertificateRequest(c *gin.Context) {
	var data cmapi.CertificateRequest
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCertificateSigningRequest(data))
}

// ORDERS
func allOrders(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllOrders(namespace))
}
func describeOrder(c *gin.Context) {
	name := c.Query("name")
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sOrder(namespace, name))
}
func deleteOrder(c *gin.Context) {
	var data v1Cert.Order
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sOrder(data))
}
func patchOrder(c *gin.Context) {
	var data v1Cert.Order
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sOrder(data))
}
func createOrder(c *gin.Context) {
	var data v1Cert.Order
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sOrder(data))
}

// ISSUERS
func allIssuers(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllIssuer(namespace))
}
func describeIssuer(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sIssuer(namespace, name))
}
func deleteIssuer(c *gin.Context) {
	var data cmapi.Issuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sIssuer(data))
}
func patchIssuer(c *gin.Context) {
	var data cmapi.Issuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sIssuer(data))
}
func createIssuer(c *gin.Context) {
	var data cmapi.Issuer
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sIssuer(data))
}

// CLUSTERISSUERS
func allClusterIssuers(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllClusterIssuers())
}
func describeClusterIssuer(c *gin.Context) {
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sClusterIssuer(name))
}
func deleteClusterIssuer(c *gin.Context) {
	var data cmapi.ClusterIssuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sClusterIssuer(data))
}
func patchClusterIssuer(c *gin.Context) {
	var data cmapi.ClusterIssuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sClusterIssuer(data))
}
func createClusterIssuer(c *gin.Context) {
	var data cmapi.ClusterIssuer
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sClusterIssuer(data))
}

// SERVICEACCOUNTS
func allServiceAccounts(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllServiceAccounts(namespace))
}
func describeServiceAccount(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sServiceAccount(namespace, name))
}
func deleteServiceAccount(c *gin.Context) {
	var data v1.ServiceAccount
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sServiceAccount(data))
}
func patchServiceAccount(c *gin.Context) {
	var data v1.ServiceAccount
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sServiceAccount(data))
}
func createServiceAccount(c *gin.Context) {
	var data v1.ServiceAccount
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sServiceAccount(data))
}

// ROLES
func allRoles(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllRoles(namespace))
}
func describeRole(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sRole(namespace, name))
}
func deleteRole(c *gin.Context) {
	var data v1Rbac.Role
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sRole(data))
}
func patchRole(c *gin.Context) {
	var data v1Rbac.Role
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sRole(data))
}
func createRole(c *gin.Context) {
	var data v1Rbac.Role
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sRole(data))
}

// ROLEBINDINGS
func allRoleBindings(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllRoleBindings(namespace))
}
func describeRoleBinding(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sRoleBinding(namespace, name))
}
func deleteRoleBinding(c *gin.Context) {
	var data v1Rbac.RoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sRoleBinding(data))
}
func patchRoleBinding(c *gin.Context) {
	var data v1Rbac.RoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sRoleBinding(data))
}
func createRoleBinding(c *gin.Context) {
	var data v1Rbac.RoleBinding
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sRoleBinding(data))
}

// CLUSTERROLES
func allClusterRoles(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllClusterRoles())
}
func describeClusterRole(c *gin.Context) {
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sClusterRole(name))
}
func deleteClusterRole(c *gin.Context) {
	var data v1Rbac.ClusterRole
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sClusterRole(data))
}
func patchClusterRole(c *gin.Context) {
	var data v1Rbac.ClusterRole
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sClusterRole(data))
}
func createClusterRole(c *gin.Context) {
	var data v1Rbac.ClusterRole
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sClusterRole(data))
}

// CLUSTERROLEBINDINGS
func allClusterRoleBindings(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllClusterRoleBindings())
}
func describeClusterRoleBinding(c *gin.Context) {
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sClusterRoleBinding(name))
}
func deleteClusterRoleBinding(c *gin.Context) {
	var data v1Rbac.ClusterRoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sClusterRoleBinding(data))
}
func patchClusterRoleBinding(c *gin.Context) {
	var data v1Rbac.ClusterRoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sClusterRoleBinding(data))
}
func createClusterRoleBinding(c *gin.Context) {
	var data v1Rbac.ClusterRoleBinding
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sClusterRoleBinding(data))
}

// VOLUMEATTACHMENTS
func allVolumeAttachments(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllVolumeAttachments())
}
func describeVolumeAttachment(c *gin.Context) {
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sVolumeAttachment(name))
}
func deleteVolumeAttachment(c *gin.Context) {
	var data v1Storage.VolumeAttachment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sVolumeAttachment(data))
}
func patchVolumeAttachment(c *gin.Context) {
	var data v1Storage.VolumeAttachment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sVolumeAttachment(data))
}
func createVolumeAttachment(c *gin.Context) {
	var data v1Storage.VolumeAttachment
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sVolumeAttachment(data))
}

// NETWORKPOLICIES
func allNetworkPolicies(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllNetworkPolicies(namespace))
}
func describeNetworkPolicy(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sNetworkPolicy(namespace, name))
}
func deleteNetworkPolicy(c *gin.Context) {
	var data v1Networking.NetworkPolicy
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sNetworkPolicy(data))
}
func patchNetworkPolicy(c *gin.Context) {
	var data v1Networking.NetworkPolicy
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sNetworkPolicy(data))
}
func createNetworkPolicy(c *gin.Context) {
	var data v1Networking.NetworkPolicy
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sNetworkpolicy(data))
}

// STORAGECLASSES
func allStorageClasses(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllStorageClasses())
}
func describeStorageClass(c *gin.Context) {
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sStorageClass(name))
}
func deleteStorageClass(c *gin.Context) {
	var data v1Storage.StorageClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sStorageClass(data))
}
func patchStorageClass(c *gin.Context) {
	var data v1Storage.StorageClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sStorageClass(data))
}
func createStorageClass(c *gin.Context) {
	var data v1Storage.StorageClass
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sStorageClass(data))
}

// CRDS
func allCrds(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllCustomResourceDefinitions())
}
func describeCrd(c *gin.Context) {
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCustomResourceDefinition(name))
}
func deleteCrd(c *gin.Context) {
	var data apiExt.CustomResourceDefinition
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sCustomResourceDefinition(data))
}
func patchCrd(c *gin.Context) {
	var data apiExt.CustomResourceDefinition
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCustomResourceDefinition(data))
}
func createCrd(c *gin.Context) {
	var data apiExt.CustomResourceDefinition
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCustomResourceDefinition(data))
}

// ENDPOINTS
func allEndpoints(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllEndpoints(namespace))
}
func describeEndpoint(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sEndpoint(namespace, name))
}
func deleteEndpoint(c *gin.Context) {
	var data v1.Endpoints
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sEndpoint(data))
}
func patchEndpoint(c *gin.Context) {
	var data v1.Endpoints
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sEndpoint(data))
}
func createEndpoint(c *gin.Context) {
	var data v1.Endpoints
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sEndpoint(data))
}

// LEASES
func allLeases(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllLeases(namespace))
}
func describeLease(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sLease(namespace, name))
}
func deleteLease(c *gin.Context) {
	var data v1Coordination.Lease
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sLease(data))
}
func patchLease(c *gin.Context) {
	var data v1Coordination.Lease
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sLease(data))
}
func createLease(c *gin.Context) {
	var data v1Coordination.Lease
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sLease(data))
}

// PRIORITYCLASSES
func allPriorityClasses(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllPriorityClasses())
}
func describePriorityClass(c *gin.Context) {
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPriorityClass(name))
}
func deletePriorityClass(c *gin.Context) {
	var data v1Scheduling.PriorityClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sPriorityClass(data))
}
func patchPriorityClass(c *gin.Context) {
	var data v1Scheduling.PriorityClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPriorityClass(data))
}
func createPriorityClass(c *gin.Context) {
	var data v1Scheduling.PriorityClass
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPriorityClass(data))
}

// VOLUMESNAPSHOTS
func allVolumeSnapshots(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllVolumeSnapshots(namespace))
}
func describeVolumeSnapshot(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sVolumeSnapshot(namespace, name))
}
func deleteVolumeSnapshot(c *gin.Context) {
	var data v6Snap.VolumeSnapshot
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sVolumeSnapshot(data))
}
func patchVolumeSnapshot(c *gin.Context) {
	var data v6Snap.VolumeSnapshot
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sVolumeSnapshot(data))
}
func createVolumeSnapshot(c *gin.Context) {
	var data v6Snap.VolumeSnapshot
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sVolumeSnapshot(data))
}

// RESOURCEQUOTAS
func allResourceQuotas(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllResourceQuotas(namespace))
}
func describeResourceQuota(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sResourceQuota(namespace, name))
}
func deleteResourceQuota(c *gin.Context) {
	var data v1.ResourceQuota
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.DeleteK8sResourceQuota(data))
}
func patchResourceQuota(c *gin.Context) {
	var data v1.ResourceQuota
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sResourceQuota(data))
}
func createResourceQuota(c *gin.Context) {
	var data v1.ResourceQuota
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sResourceQuota(data))
}

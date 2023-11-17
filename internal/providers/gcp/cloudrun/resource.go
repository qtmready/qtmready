package cloudrun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"cloud.google.com/go/iam/apiv1/iampb"
	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/gocql/gocql"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
	api "google.golang.org/genproto/googleapis/api"

	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/shared"
)

type (
	Resource struct {
		ID                         gocql.UUID
		Cpu                        string
		Memory                     string
		Generation                 uint8
		Port                       int32
		Envs                       map[string]string
		OutputEnvs                 map[string]string
		Region                     string // from blueprint
		Image                      string // from workload
		Config                     map[string]interface{}
		Name                       string
		Revision                   string
		LastRevision               string
		MinInstances               int32
		MaxInstances               int32
		AllowUnauthenticatedAccess bool
		CpuIdle                    bool
		Project                    string
		ServiceName                string
	}

	Workload struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	GCPConfig struct {
		Project string
	}
)

var (
	activities *Activities
)

// Marshal marshals the Resource object.
func (r *Resource) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Provision provisions the cloud resource.
func (r *Resource) Provision(ctx workflow.Context) (workflow.Future, error) {
	// do nothing, the infra will be provisioned with deployment
	return nil, nil
}

// DeProvision deprovisions the cloudrun resource.
func (r *Resource) DeProvision() error {
	return nil
}

// UpdateTraffic updates the traffic distribution on latest and previous revision as per the input
// parameter trafficpcnt is the percentage traffic to be deployed on latest revision.
// UpdateTraffic will execute a workflow to update the resource. This workflow is not directly called
// from provisioninfra workflow to avoid passing resource interface as argument.
func (r *Resource) UpdateTraffic(ctx workflow.Context, percent int32) error {
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		ChildWorkflowOptions(
			shared.WithWorkflowParent(ctx),
			shared.WithWorkflowBlock("Resource"),
			shared.WithWorkflowBlockID(r.Name),
			shared.WithWorkflowElement("UpdateTraffic"),
		)

	cctx := workflow.WithChildOptions(ctx, opts)

	shared.Logger().Info("Executing Update traffic workflow")

	w := &workflows{}
	err := workflow.
		ExecuteChildWorkflow(cctx, w.UpdateTraffic, r, percent).Get(cctx, nil)

	if err != nil {
		shared.Logger().Error("Could not execute UpdateTraffic workflow", "error", err)
		return err
	}

	return nil
}

// Deploy deploys the resource with the given changeset ID.
func (r *Resource) Deploy(ctx workflow.Context, wl []core.Workload, changesetID gocql.UUID) error {
	shared.Logger().Info("deploying", "cloudrun", r, "workload", wl)

	if len(wl) != 1 {
		shared.Logger().Error("Cannot deploy more than one workloads on cloud run", "number of workloads", len(wl))
		return errors.New("multiple workloads defined for cloud run")
	}

	// provision with execute a workflow to provision the resources. This workflow is not directly called
	// from provisioninfra workflow to avoid passing resource interface as argument

	workload := &Workload{}
	if err := json.Unmarshal([]byte(wl[0].Container), workload); err != nil {
		return err
	}

	workload.Image = workload.Image + ":" + changesetID.String()
	workload.Name = wl[0].Name

	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		ChildWorkflowOptions(
			shared.WithWorkflowParent(ctx),
			shared.WithWorkflowBlock("Resource"),
			shared.WithWorkflowBlockID(r.Name),
			shared.WithWorkflowElement("Deploy"),
		)

	cctx := workflow.WithChildOptions(ctx, opts)

	shared.Logger().Info("starting DeployCloudRun workflow")

	w := &workflows{}
	err := workflow.
		ExecuteChildWorkflow(cctx, w.DeployWorkflow, r, workload).Get(cctx, r)

	if err != nil {
		shared.Logger().Error("Could not start DeployCloudRun workflow", "error", err)
		return err
	}

	return nil
}

func (r *Resource) GetServiceClient() (*run.ServicesClient, error) {
	client, err := run.NewServicesRESTClient(context.Background())
	if err != nil {
		shared.Logger().Error("New service rest client", "error", err)
		return nil, err
	}

	return client, err
}

// GetService gets a cloud run service from GCP.
func (r *Resource) GetService(ctx context.Context) *runpb.Service {
	logger := activity.GetLogger(ctx)

	serviceClient, err := run.NewServicesRESTClient(ctx)
	if err != nil {
		shared.Logger().Error("New service rest client", "Error", err)
		return nil
	}

	defer serviceClient.Close()

	svcpath := r.GetParent() + "/services/" + r.ServiceName
	req := &runpb.GetServiceRequest{Name: svcpath}

	svc, err := serviceClient.GetService(ctx, req)

	if err != nil {
		logger.Error("Get Service", "Error, returning nil", err)
		return nil
	}

	logger.Debug("Get service", "service", svc, "error", err)

	return svc
}

// AllowAccessToAll Sets IAM policy to allow access to all users.
func (r *Resource) AllowAccessToAll(ctx context.Context) error {
	logger := activity.GetLogger(ctx)

	client, err := run.NewServicesRESTClient(context.Background())
	if err != nil {
		logger.Error("New service rest client", "Error", err)
		return nil
	}

	defer func() { _ = client.Close() }()

	rsc := r.GetParent() + "/services/" + r.ServiceName

	binding := new(iampb.Binding)
	binding.Members = []string{"allUsers"}
	binding.Role = "roles/run.invoker"

	_, err = client.SetIamPolicy(
		context.Background(),
		&iampb.SetIamPolicyRequest{Resource: rsc, Policy: &iampb.Policy{Bindings: []*iampb.Binding{binding}}},
	)
	if err != nil {
		logger.Error("Set policy", "Error", err)
		return err
	}

	return nil
}

// GetServiceTemplate creates and returns the revision template for cloud run from the workload to be deployed
// revision template specifies the resource requirements, image to be deployed and traffic distribution etc.
// this template will be used for first deployment only, from next deployments the already deployed template will be
// fetched from cloudrun and the same will be used for next revision.
// TODO: the above design will not work if resource definition is changed.
func (r *Resource) GetServiceTemplate(ctx context.Context, wl *Workload) *runpb.Service {
	activity.GetLogger(ctx).Info("setting service template for", "revision", r.Revision)

	// unmarshaling the container here assuming that container definition will be specific to a resource
	// this can be done at a common location if the container definition turns out to be same for all resources

	revisionTemplate := r.GetRevisionTemplate(ctx, wl)
	launchStage := r.GetLaunchStageConfig(ctx)
	ingress := r.GetIngressConfig(ctx)

	service := &runpb.Service{
		Template:    revisionTemplate,
		LaunchStage: launchStage,
		Ingress:     ingress,
	}

	tt := &runpb.TrafficTarget{
		Type:    runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST,
		Percent: 100,
	}
	service.Traffic = []*runpb.TrafficTarget{tt}

	return service
}

func (r *Resource) GetLaunchStageConfig(ctx context.Context) api.LaunchStage {
	launchStageConfig := r.Config["launch_stage"].(string)
	launchStage := api.LaunchStage(api.LaunchStage_value[launchStageConfig])

	return launchStage
}

func (r *Resource) GetIngressConfig(ctx context.Context) runpb.IngressTraffic {
	ingressConfig := r.Config["ingress"].(string)
	ingress := runpb.IngressTraffic(runpb.IngressTraffic_value[ingressConfig])

	return ingress
}

func (r *Resource) GetRevisionTemplate(ctx context.Context, wl *Workload) *runpb.RevisionTemplate {

	container := r.GetContainerConfig(ctx, wl)
	scaling := r.GetRevisionScalingConfig(ctx)
	executionEnv := r.GetExecEnvConfig(ctx)
	vpcAccess := r.GetVpcAccessConfig(ctx)
	volumes := r.GetVolumesConfig(ctx)

	revisionTemplate := &runpb.RevisionTemplate{
		Containers:           []*runpb.Container{container},
		Scaling:              scaling,
		ExecutionEnvironment: executionEnv,
		Revision:             r.Revision,
		VpcAccess:            vpcAccess,
		Volumes:              volumes,
	}

	return revisionTemplate
}

func (r *Resource) GetContainerConfig(ctx context.Context, wl *Workload) *runpb.Container {

	templateConfig := r.Config["template"].(map[string]interface{})
	templateContainersConfig := templateConfig["containers"].(map[string]interface{})

	resourceReq := r.GetResourceReqConfig(ctx)
	containerPorts := r.GetContainerPortsConfig(ctx)

	Envs := []*runpb.EnvVar{}

	env := templateContainersConfig["env"].([]interface{})
	for _, val := range env {
		envVal := val.(map[string]interface{})
		Envs = append(Envs, &runpb.EnvVar{
			Name: fmt.Sprint(envVal["name"]),
			Values: &runpb.EnvVar_Value{
				Value: fmt.Sprint(envVal["value"])},
		})
	}

	volumeMountsName := templateContainersConfig["volume_mounts"].(map[string]interface{})["name"].(string)
	volumeMountsPath := templateContainersConfig["volume_mounts"].(map[string]interface{})["mount_path"].(string)
	volumeMounts := &runpb.VolumeMount{
		Name:      volumeMountsName,
		MountPath: volumeMountsPath,
	}
	volumeMountsArray := []*runpb.VolumeMount{volumeMounts}

	container := &runpb.Container{
		Name:         wl.Name,
		Image:        wl.Image,
		Resources:    resourceReq,
		Ports:        containerPorts,
		Env:          Envs,
		VolumeMounts: volumeMountsArray,
	}

	return container
}

func (r *Resource) GetContainerPortsConfig(ctx context.Context) []*runpb.ContainerPort {

	templateConfig := r.Config["template"].(map[string]interface{})
	templateContainersConfig := templateConfig["containers"].(map[string]interface{})

	containerPortConfig := templateContainersConfig["ports"].(map[string]interface{})["container_port"].(string)
	containerPort64bit, _ := strconv.ParseInt(containerPortConfig, 10, 32)
	containerPort := &runpb.ContainerPort{ContainerPort: int32(containerPort64bit)}

	containerPortArray := []*runpb.ContainerPort{containerPort}

	return containerPortArray
}

func (r *Resource) GetResourceReqConfig(ctx context.Context) *runpb.ResourceRequirements {

	templateConfig := r.Config["template"].(map[string]interface{})
	templateContainersConfig := templateConfig["containers"].(map[string]interface{})

	cpuIdleStr := templateContainersConfig["resources"].(map[string]interface{})["cpu_idle"].(string)
	cpuIdle, _ := strconv.ParseBool(cpuIdleStr)
	resources := &runpb.ResourceRequirements{
		CpuIdle: cpuIdle,
	}

	return resources
}

func (r *Resource) GetExecEnvConfig(ctx context.Context) runpb.ExecutionEnvironment {

	templateConfig := r.Config["template"].(map[string]interface{})
	executionEnvConfig := templateConfig["execution_environment"].(string)
	executionEnv := runpb.ExecutionEnvironment(
		runpb.ExecutionEnvironment_value[executionEnvConfig],
	)

	return executionEnv
}

func (r *Resource) GetVpcAccessConfig(ctx context.Context) *runpb.VpcAccess {

	templateConfig := r.Config["template"].(map[string]interface{})
	templateVpcAccessConfig := templateConfig["vpc_access"].(map[string]interface{})

	networkInterfaces := templateVpcAccessConfig["network_interfaces"].(map[string]interface{})
	networkInterfaceArray := []*runpb.VpcAccess_NetworkInterface{
		{
			Network:    fmt.Sprint(networkInterfaces["network"]),
			Subnetwork: fmt.Sprint(networkInterfaces["subnetwork"]),
		},
	}

	egress := templateVpcAccessConfig["egress"].(string)
	vpcAccess := &runpb.VpcAccess{
		Egress:            runpb.VpcAccess_VpcEgress(runpb.VpcAccess_VpcEgress_value[egress]),
		NetworkInterfaces: networkInterfaceArray,
	}
	return vpcAccess
}

func (r *Resource) GetVolumesConfig(ctx context.Context) []*runpb.Volume {

	templateConfig := r.Config["template"].(map[string]interface{})
	volumeName := templateConfig["volumes"].(map[string]interface{})["name"].(string)
	volumeInstance := templateConfig["volumes"].(map[string]interface{})["cloud_sql_instance"].(map[string]interface{})["instances"].([]interface{})

	VolumeInstanceArray := []string{}
	for _, instanceVal := range volumeInstance {
		val := instanceVal.(map[string]interface{})
		VolumeInstanceArray = append(VolumeInstanceArray, val["name"].(string))
	}
	volume := &runpb.Volume{
		Name: volumeName,
		VolumeType: &runpb.Volume_CloudSqlInstance{
			CloudSqlInstance: &runpb.CloudSqlInstance{
				Instances: VolumeInstanceArray,
			},
		},
	}
	volumes := []*runpb.Volume{volume}
	return volumes
}

func (r *Resource) GetRevisionScalingConfig(ctx context.Context) *runpb.RevisionScaling {

	templateConfig := r.Config["template"].(map[string]interface{})

	scalingConfig := templateConfig["scaling"].(map[string]interface{})
	minInstanceConfig := scalingConfig["min_instance_count"].(string)
	minInstanceConfig64bit, _ := strconv.ParseInt(minInstanceConfig, 10, 32)
	maxInstanceConfig := scalingConfig["max_instance_count"].(string)
	maxInstanceConfig64bit, _ := strconv.ParseInt(maxInstanceConfig, 10, 32)

	scaling := &runpb.RevisionScaling{
		MinInstanceCount: int32(minInstanceConfig64bit),
		MaxInstanceCount: int32(maxInstanceConfig64bit),
	}

	return scaling
}

func (r *Resource) GetParent() string {
	return "projects/" + r.Project + "/locations/" + r.Region
}

func (r *Resource) GetFirstRevision() string {
	return r.ServiceName + "-0"
}

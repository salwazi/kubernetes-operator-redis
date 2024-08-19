package controller

import (
	"context"

	cachev1alpha1 "github.com/salwazi/kubernetes-operator-redis/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// deploymentForRedis returns a Redis Deployment object
func (r *RedisReconciler) deploymentForRedis(redis *cachev1alpha1.Redis, secretName string) *appsv1.Deployment {
	labels := map[string]string{
		"app": redis.Name,
	}
	replicas := redis.Spec.Replicas

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: redis.Spec.Image + ":" + redis.Spec.Version,
						Name:  redis.Name,
						Env: []corev1.EnvVar{
							{
								Name: "REDIS_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: "password",
									},
								},
							},
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse(redis.Spec.Resources.Requests.CPU),
								corev1.ResourceMemory: resource.MustParse(redis.Spec.Resources.Requests.Memory),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse(redis.Spec.Resources.Limits.CPU),
								corev1.ResourceMemory: resource.MustParse(redis.Spec.Resources.Limits.Memory),
							},
						},
					}},
				},
			},
		},
	}
}

// updateDeploymentAndStatus updates the Deployment and status of a Redis resource.
// It compares the Redis resource with the existing Deployment and makes necessary updates
// to the Deployment size, image, resources, and update strategy.
func (r *RedisReconciler) updateDeploymentAndStatus(ctx context.Context, redis *cachev1alpha1.Redis, foundDeployment *appsv1.Deployment) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Update the Deployment size if necessary
	size := redis.Spec.Replicas
	if *foundDeployment.Spec.Replicas != size {
		foundDeployment.Spec.Replicas = &size
		err := r.Update(ctx, foundDeployment)
		if err != nil {
			logger.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			return ctrl.Result{}, err
		}
		logger.Info("Updated Deployment size", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name, "Size", size)
	}

	// Check for updates in the image or version
	if foundDeployment.Spec.Template.Spec.Containers[0].Image != (redis.Spec.Image + ":" + redis.Spec.Version) {
		foundDeployment.Spec.Template.Spec.Containers[0].Image = (redis.Spec.Image + ":" + redis.Spec.Version)
		err := r.Update(ctx, foundDeployment)
		if err != nil {
			logger.Error(err, "Failed to update Deployment image and version", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			return ctrl.Result{}, err
		}
		logger.Info("Updated Deployment image and version", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name, "Image", redis.Spec.Image, "Version", redis.Spec.Version)
	}

	// Check for updates in other fields (e.g., resources, storage, etc.)
	if foundDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String() != redis.Spec.Resources.Requests.CPU ||
		foundDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String() != redis.Spec.Resources.Requests.Memory ||
		foundDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String() != redis.Spec.Resources.Limits.CPU ||
		foundDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String() != redis.Spec.Resources.Limits.Memory {
		resources := foundDeployment.Spec.Template.Spec.Containers[0].Resources
		resources.Requests[corev1.ResourceCPU] = resource.MustParse(redis.Spec.Resources.Requests.CPU)
		foundDeployment.Spec.Template.Spec.Containers[0].Resources = resources
		resources.Requests[corev1.ResourceMemory] = resource.MustParse(redis.Spec.Resources.Requests.Memory)
		resources.Requests[corev1.ResourceMemory] = resource.MustParse(redis.Spec.Resources.Requests.Memory)
		foundDeployment.Spec.Template.Spec.Containers[0].Resources = resources
		resources.Limits[corev1.ResourceCPU] = resource.MustParse(redis.Spec.Resources.Limits.CPU)
		resources.Limits[corev1.ResourceMemory] = resource.MustParse(redis.Spec.Resources.Limits.Memory)
		err := r.Update(ctx, foundDeployment)
		if err != nil {
			logger.Error(err, "Failed to update Deployment resources", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			return ctrl.Result{}, err
		}
		logger.Info("Updated Deployment resources", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name, "Resources", redis.Spec.Resources)
	}

	err := r.updateRedisStatus(ctx, redis, foundDeployment)
	if err != nil {
		logger.Error(err, "Failed to update Redis status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

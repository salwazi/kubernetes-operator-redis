package controller

import (
	"testing"

	cachev1alpha1 "github.com/salwazi/kubernetes-operator-redis/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestDeploymentForRedis tests the deploymentForRedis function
func TestDeploymentForRedis(t *testing.T) {
	// Arrange
	redis := &cachev1alpha1.Redis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-redis",
			Namespace: "default",
		},
		Spec: cachev1alpha1.RedisSpec{
			Replicas: 3,
			Image:    "redis",
			Version:  "6.2",
			Resources: cachev1alpha1.RedisResources{
				Requests: cachev1alpha1.Requests{
					CPU:    "500m",
					Memory: "512Mi",
				},
				Limits: cachev1alpha1.Limits{
					CPU:    "1",
					Memory: "1Gi",
				},
			},
		},
	}

	secretName := "redis-secret"
	r := &RedisReconciler{}

	// Act
	deployment, _ := r.deploymentForRedis(redis, secretName)

	// Assert
	assert.Equal(t, redis.Name, deployment.Name, "Deployment name should match Redis name")
	assert.Equal(t, redis.Namespace, deployment.Namespace, "Deployment namespace should match Redis namespace")
	assert.Equal(t, redis.Spec.Replicas, *deployment.Spec.Replicas, "Deployment replicas should match Redis replicas")
	assert.Equal(t, redis.Spec.Image+":6.2", deployment.Spec.Template.Spec.Containers[0].Image, "Deployment image should match Redis image and version")
	assert.Equal(t, secretName, deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.SecretKeyRef.Name, "Secret name should match")
	assert.Equal(t, "password", deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.SecretKeyRef.Key, "Secret key should match")

	requests := deployment.Spec.Template.Spec.Containers[0].Resources.Requests
	limits := deployment.Spec.Template.Spec.Containers[0].Resources.Limits
	assert.Equal(t, "500m", requests.Cpu().String(), "CPU request should match")
	assert.Equal(t, "512Mi", requests.Memory().String(), "Memory request should match")
	assert.Equal(t, "1", limits.Cpu().String(), "CPU limit should match")
	assert.Equal(t, "1Gi", limits.Memory().String(), "Memory limit should match")
}

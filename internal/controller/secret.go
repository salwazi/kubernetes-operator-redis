package controller

import (
	"encoding/base64"
	"math/rand"

	cachev1alpha1 "github.com/salwazi/kubernetes-operator-redis/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// createSecret creates a new Kubernetes Secret for storing the Redis password
func (r *RedisReconciler) createSecret(redis *cachev1alpha1.Redis, password string) (*corev1.Secret, error) {
	labels := map[string]string{
		"app": redis.Name,
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name + "-secret",
			Namespace: redis.Namespace,
			Labels:    labels,
		},
		Data: map[string][]byte{
			"password": []byte(base64.StdEncoding.EncodeToString([]byte(password))),
		},
	}
	// Set Redis instance as the owner of the Secret
	if err := controllerutil.SetControllerReference(redis, secret, r.Scheme); err != nil {
		return nil, err
	}

	return secret, nil
}

// generateRandomPassword creates a secure random password
func generateRandomPassword() (string, error) {
	bytePassword := make([]byte, 16) // Adjust length as needed
	_, err := rand.Read(bytePassword)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytePassword), nil
}

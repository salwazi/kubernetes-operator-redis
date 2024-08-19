package controller

import (
	"encoding/base64"
	"testing"

	cachev1alpha1 "github.com/salwazi/kubernetes-operator-redis/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestCreateSecret tests the createSecret function
func TestCreateSecret(t *testing.T) {
	// Arrange
	redis := &cachev1alpha1.Redis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-redis",
			Namespace: "default",
		},
	}

	password := "mysecurepassword"

	r := &RedisReconciler{}

	// Act
	secret := r.createSecret(redis, password)

	// Assert
	assert.Equal(t, redis.Namespace, secret.Namespace, "Secret namespace should match")
	assert.Equal(t, redis.Name, secret.Labels["app"], "Secret label 'app' should match Redis name")

	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))
	assert.Equal(t, encodedPassword, string(secret.Data["password"]), "Secret password should be base64 encoded")
}

// TestGenerateRandomPassword tests the generateRandomPassword function
func TestGenerateRandomPassword(t *testing.T) {
	// Act
	password, err := generateRandomPassword()

	// Assert
	assert.NoError(t, err, "generateRandomPassword should not return an error")
	assert.NotEmpty(t, password, "Generated password should not be empty")

	// Check that the password can be decoded from base64
	decodedPassword, err := base64.URLEncoding.DecodeString(password)
	assert.NoError(t, err, "Generated password should be valid base64")
	assert.Len(t, decodedPassword, 16, "Decoded password should be 16 bytes long")
}

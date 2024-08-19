package controller

import (
	"context"

	cachev1alpha1 "github.com/salwazi/kubernetes-operator-redis/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
func removeString(slice []string, s string) []string {
	result := []string{}
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}
func (r *RedisReconciler) reconcileFinalizer(ctx context.Context, redis *cachev1alpha1.Redis) (ctrl.Result, error) {
	if redis.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then we add the finalizer and update the object. This is necessary to
		// ensure that our controller can clean up resources before the object is deleted.
		if !containsString(redis.ObjectMeta.Finalizers, redisFinalizer) {
			redis.ObjectMeta.Finalizers = append(redis.ObjectMeta.Finalizers, redisFinalizer)
			if err := r.Update(ctx, redis); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(redis.ObjectMeta.Finalizers, redisFinalizer) {
			// our finalizer is present, so lets handle any dependency
			if err := r.deleteDependantResources(ctx, redis); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			redis.ObjectMeta.Finalizers = removeString(redis.ObjectMeta.Finalizers, redisFinalizer)
			if err := r.Update(ctx, redis); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// updateRedisStatus updates the status of the Redis CR
func (r *RedisReconciler) updateRedisStatus(ctx context.Context, redis *cachev1alpha1.Redis, deployment *appsv1.Deployment) error {
	redis.Status.ReadyReplicas = deployment.Status.ReadyReplicas
	redis.Status.TotalReplicas = deployment.Status.Replicas
	redis.Status.Conditions = deployment.Status.Conditions
	return r.Status().Update(ctx, redis)
}

// Implement the deleteExternalResources function to clean up any external resources
func (r *RedisReconciler) deleteDependantResources(ctx context.Context, redis *cachev1alpha1.Redis) error {

	// Delete the Secret
	secretName := redis.Spec.SecretName
	if secretName != "" {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: redis.Namespace,
			},
		}
		err := r.Delete(ctx, secret)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}

	// Delete the Deployment
	deploymentName := redis.Name
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: redis.Namespace,
		},
	}
	err := r.Delete(ctx, deployment)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

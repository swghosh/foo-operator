package controller

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/swghosh/foo-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"reconciler.io/runtime/reconcilers"
)

// TheFooController is a resource reconciler constituting multiple sub reconcilers
func TheFooController(c reconcilers.Config) *reconcilers.ResourceReconciler[*v1.FooManager] {
	return &reconcilers.ResourceReconciler[*v1.FooManager]{
		Name: "FooManager",
		Reconciler: &reconcilers.Sequence[*v1.FooManager]{
			FooUseController(c),
			FooDeploymentController(c),
		},
		Config: c,
	}
}

// FooUseController puts spec.useIn -> spec.wasUsedFor
func FooUseController(c reconcilers.Config) reconcilers.SubReconciler[*v1.FooManager] {
	return &reconcilers.SyncReconciler[*v1.FooManager]{
		Name: "SyncUseInWithUsedFor",
		Sync: func(ctx context.Context, resource *v1.FooManager) error {
			if resource.Status.WasUsedFor == resource.Spec.UseIn {
				return nil
			}

			resource.Status = v1.FooManagerStatus{
				WasUsedFor: resource.Spec.UseIn,
				LastTransitionTime: &metav1.Time{
					Time: time.Now(),
				},
			}
			return nil
		},
	}
}

func FooDeploymentController(c reconcilers.Config) reconcilers.SubReconciler[*v1.FooManager] {
	return &reconcilers.ChildReconciler[*v1.FooManager, *appsv1.Deployment, *appsv1.DeploymentList]{
		Name: "SyncDeployment",
		DesiredChild: func(ctx context.Context, resource *v1.FooManager) (*appsv1.Deployment, error) {
			if resource.Spec.TemplateSpec == nil {
				return nil, nil
			}

			labels := map[string]string{
				"app.kubernetes.io/name":       "foo-deployment",
				"app.kubernetes.io/instance":   resource.Name,
				"app.kubernetes.io/managed-by": "foo-operator",
			}

			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resource.Name + "-deployment",
					Namespace: resource.Namespace,
					Labels:    labels,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: labels,
						},
						Spec: *resource.Spec.TemplateSpec,
					},
				},
			}

			return deployment, nil
		},
		ReflectChildStatusOnParent: func(ctx context.Context, parent *v1.FooManager, child *appsv1.Deployment, err error) {
			if err != nil || child == nil {
				parent.Status.WeirdReport = fmt.Sprintf("we-saw-an-error-and-couldnt-avoid-reporting: %q", err)
				parent.Status.LastTransitionTime = &metav1.Time{
					Time: time.Now(),
				}

				parent.Status.DeploymentStatus = nil
				return
			}

			parent.Status.DeploymentStatus = child.Status.DeepCopy()
		},
		ChildObjectManager: &reconcilers.UpdatingObjectManager[*appsv1.Deployment]{
			MergeBeforeUpdate: func(current, desired *appsv1.Deployment) {},
		},
	}
}

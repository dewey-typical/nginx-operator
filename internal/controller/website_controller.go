/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	webv1alpha1 "github.com/dewey-typical/nginx-operator/api/v1alpha1"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// WebsiteReconciler reconciles a Website object
type WebsiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=web.dewey-typical.github.io,resources=websites,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=web.dewey-typical.github.io,resources=websites/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=web.dewey-typical.github.io,resources=websites/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Website object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *WebsiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	website := &webv1alpha1.Website{}

	err := r.Get(ctx, req.NamespacedName, website)

	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func buildDeployment(website webv1alpha1.Website) appv1.Deployment {
	var replicas int32 = 2
	var port int32 = 80
	deployName := "website-" + website.Spec.Title
	deploy := appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployName,
			Namespace: "default",
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{Image: "nginx:alpine",
							Name: website.Spec.Title,
							Ports: []v1.ContainerPort{
								{ContainerPort: port}},
							VolumeMounts: []v1.VolumeMount{
								{Name: "html-volume",
									MountPath: "/usr/share/nginx/html"},
							},
						},
					},
					Volumes: []v1.Volume{
						{Name: "html-volume",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: deployName,
									},
									Items: []v1.KeyToPath{
										{Key: "index.html"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return deploy
}

func buildContent(website webv1alpha1.Website) v1.ConfigMap {
	configName := "website-" + website.Spec.Title
	config := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: "default",
		},
		Data: map[string]string{
			"index.html": website.Spec.Content,
		},
	}
	return config
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebsiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webv1alpha1.Website{}).
		Complete(r)
}

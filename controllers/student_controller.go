/*
Copyright 2022.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1 "github.com/gitamit544/k8s_operator/api/v1"
)

// StudentReconciler reconciles a Student object
type StudentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=operator.my.domain,resources=students,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.my.domain,resources=students/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.my.domain,resources=students/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Student object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *StudentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	student := &operatorv1.Student{}
	err := r.Get(ctx, req.NamespacedName, student)
	if err != nil {
		logger.Error(err, "Error getting student spec")
		return ctrl.Result{}, err
	}

	marks := student.Spec.Marks
	if marks >= 90 && marks < 100 {
		student.Status.Level = "A"
	} else if marks >= 80 && marks < 90 {
		student.Status.Level = "B"
	} else {
		student.Status.Level = "C"
	}

	err = r.Status().Update(ctx, student)
	if err != nil {
		logger.Error(err, "Error updating student")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StudentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1.Student{}).
		Complete(r)
}

package clusterrolebindings

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type CreateOptions struct {
	RESTConfig       *rest.Config
	Name             string
	SubjectName      string
	SubjectNamespace string
}

func Create(ctx context.Context, opts CreateOptions) error {
	gvr := schema.GroupVersionResource{
		Group:    "rbac.authorization.k8s.io",
		Version:  "v1",
		Resource: "clusterrolebindings",
	}

	crb := rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				core.InstalledByLabel: core.InstalledByValue,
			},
			Name: opts.Name, //fmt.Sprintf(clusterRoleBindingNamePattern, provider),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      opts.SubjectName,
				Namespace: opts.SubjectNamespace,
			},
		},
	}

	dat, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&crb)
	if err != nil {
		return err
	}

	obj := unstructured.Unstructured{}
	obj.SetUnstructuredContent(dat)

	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return err
	}

	_, err = dc.Resource(gvr).Create(context.TODO(), &obj, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

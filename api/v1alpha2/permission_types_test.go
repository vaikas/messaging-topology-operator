package v1alpha2

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Permission spec", func() {
	var (
		namespace = "default"
		ctx       = context.Background()
	)

	It("creates a permission with no vhost permission configured", func() {
		permission := Permission{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-permission-0",
				Namespace: namespace,
			},
			Spec: PermissionSpec{
				User:        "test",
				Vhost:       "/test",
				Permissions: VhostPermissions{},
				RabbitmqClusterReference: RabbitmqClusterReference{
					Name: "some-cluster",
				},
			},
		}
		Expect(k8sClient.Create(ctx, &permission)).To(Succeed())
		fetchedPermission := &Permission{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{
			Name:      permission.Name,
			Namespace: permission.Namespace,
		}, fetchedPermission)).To(Succeed())
		Expect(fetchedPermission.Spec.User).To(Equal("test"))
		Expect(fetchedPermission.Spec.Vhost).To(Equal("/test"))
		Expect(fetchedPermission.Spec.RabbitmqClusterReference.Name).To(Equal("some-cluster"))

		Expect(fetchedPermission.Spec.Permissions.Configure).To(Equal(""))
		Expect(fetchedPermission.Spec.Permissions.Write).To(Equal(""))
		Expect(fetchedPermission.Spec.Permissions.Read).To(Equal(""))
	})

	It("creates a permission with vhost permissions all configured", func() {
		permission := Permission{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-permission-1",
				Namespace: namespace,
			},
			Spec: PermissionSpec{
				User:  "test",
				Vhost: "/test",
				Permissions: VhostPermissions{
					Configure: ".*",
					Read:      "^?",
					Write:     ".*",
				},
				RabbitmqClusterReference: RabbitmqClusterReference{
					Name: "some-cluster",
				},
			},
		}
		Expect(k8sClient.Create(ctx, &permission)).To(Succeed())
		fetchedPermission := &Permission{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{
			Name:      permission.Name,
			Namespace: permission.Namespace,
		}, fetchedPermission)).To(Succeed())
		Expect(fetchedPermission.Spec.User).To(Equal("test"))
		Expect(fetchedPermission.Spec.Vhost).To(Equal("/test"))
		Expect(fetchedPermission.Spec.RabbitmqClusterReference.Name).To(Equal("some-cluster"))

		Expect(fetchedPermission.Spec.Permissions.Configure).To(Equal(".*"))
		Expect(fetchedPermission.Spec.Permissions.Write).To(Equal(".*"))
		Expect(fetchedPermission.Spec.Permissions.Read).To(Equal("^?"))
	})
})

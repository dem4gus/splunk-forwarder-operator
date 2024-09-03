package proxy

import (
	"context"
	"testing"

	"golang.org/x/exp/maps"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/splunk-forwarder-operator/api/v1alpha1"
	"github.com/openshift/splunk-forwarder-operator/config"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	fakekubeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var caBundle = &corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "trusted-ca-bundle",
		Namespace: "openshift-config-managed",
	},
	Data: map[string]string{
		"ca-bundle.crt": "test-ca-bundle",
	},
}

func TestReconcileProxy_Reconcile(t *testing.T) {
	if err := configv1.AddToScheme(scheme.Scheme); err != nil {
		t.Errorf("ProxyReconciler.Reconcile() error = %v", err)
		return
	}
	type args struct {
		request reconcile.Request
	}
	tests := []struct {
		name         string
		args         args
		want         map[string]string
		wantErr      bool
		localObjects []runtime.Object
	}{
		{
			name: "no proxy",
			args: args{
				request: reconcile.Request{NamespacedName: types.NamespacedName{Name: "test", Namespace: "openshift-test"}},
			},
			want:    map[string]string{},
			wantErr: false,
			localObjects: []runtime.Object{
				&configv1.Proxy{
					ObjectMeta: v1.ObjectMeta{Name: "cluster"},
					Spec:       configv1.ProxySpec{},
				},
				caBundle,
			},
		},
		{
			name: "http proxy",
			args: args{
				request: reconcile.Request{NamespacedName: types.NamespacedName{Name: "test", Namespace: "openshift-test"}},
			},
			want: map[string]string{
				"httpProxy": "http://proxy.example.com:8080",
				"noProxy":   "localhost",
			},
			wantErr: false,
			localObjects: []runtime.Object{
				&v1alpha1.SplunkForwarder{
					ObjectMeta: v1.ObjectMeta{Name: "test", Namespace: "openshift-test"},
				},
				&configv1.Proxy{
					ObjectMeta: v1.ObjectMeta{Name: "cluster"},
					Spec: configv1.ProxySpec{
						HTTPProxy: "http://proxy.example.com:8080",
						NoProxy:   "localhost",
					},
				},
				caBundle,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v1alpha1.AddToScheme(scheme.Scheme); err != nil {
				t.Errorf("ProxyReconciler.Reconcile() error = %v", err)
			}
			cm := &corev1.ConfigMap{
				ObjectMeta: v1.ObjectMeta{
					Name:      config.ProxyConfigMapName,
					Namespace: "openshift-test",
				},
				Data: map[string]string{},
			}
			ctx := context.Background()
			tt.localObjects = append(tt.localObjects, runtime.Object(cm))
			fakeClient := fakekubeclient.NewClientBuilder().WithScheme(scheme.Scheme).WithRuntimeObjects(tt.localObjects...).Build()

			r := &ProxyReconciler{
				Client: fakeClient,
				Scheme: scheme.Scheme,
				Config: cm,
			}
			_, err := r.Reconcile(ctx, tt.args.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProxyReconciler.Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotConfig := &corev1.ConfigMap{}
			err = fakeClient.Get(ctx, types.NamespacedName{Namespace: cm.Namespace, Name: cm.Name}, gotConfig)
			if err != nil {
				t.Errorf("got unexpected error: %v", err)
			}

			if !(maps.Equal(gotConfig.Data, tt.want)) {
				t.Errorf("got = %v, want = %v", gotConfig.Data, tt.want)
			}
		})
	}
}

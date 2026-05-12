package secret

import (
	"context"
	"reflect"
	"testing"

	sfv1alpha1 "github.com/openshift/splunk-forwarder-operator/api/v1alpha1"
	"github.com/openshift/splunk-forwarder-operator/config"
	"github.com/openshift/splunk-forwarder-operator/internal/testutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	fakekubeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileSecret_Reconcile(t *testing.T) {
	if err := sfv1alpha1.AddToScheme(scheme.Scheme); err != nil {
		t.Errorf("SecretReconciler.Reconcile() error = %v", err)
		return
	}
	type args struct {
		request reconcile.Request
	}
	tests := []struct {
		name         string
		args         args
		want         reconcile.Result
		wantErr      bool
		localObjects []runtime.Object
	}{
		{
			name: "Reconcile succeeds when SplunkForwarder CRD is not found",
			args: args{
				request: reconcile.Request{},
			},
			want:         reconcile.Result{},
			wantErr:      false,
			localObjects: []runtime.Object{},
		},
		{
			name: "Reconcile succeeds when splunk-auth secret is missing",
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      config.SplunkAuthSecretName,
						Namespace: testutil.InstanceNamespace,
					},
				},
			},
			want:    reconcile.Result{},
			wantErr: false,
			localObjects: []runtime.Object{
				testutil.NewSplunkForwarderCR().Build(),
			},
		},
		{
			name: "Reconcile succeeds when DaemonSet does not exist yet",
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      config.SplunkAuthSecretName,
						Namespace: testutil.InstanceNamespace,
					},
				},
			},
			want:    reconcile.Result{},
			wantErr: false,
			localObjects: []runtime.Object{
				testutil.NewSplunkForwarderCR().Build(),
				testutil.NewSplunkAuthSecret(),
			},
		},
		{
			name: "Reconcile updates DaemonSet restart timestamp when secret changes",
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      config.SplunkAuthSecretName,
						Namespace: testutil.InstanceNamespace,
					},
				},
			},
			want:    reconcile.Result{},
			wantErr: false,
			localObjects: []runtime.Object{
				testutil.NewSplunkForwarderCR().Build(),
				testutil.NewSplunkAuthSecret(),
				testutil.NewSplunkForwarderDaemonSet(),
			},
		},
		{
			name: "Reconcile updates DaemonSet timestamp when HEC secret is present",
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      config.SplunkAuthSecretName,
						Namespace: testutil.InstanceNamespace,
					},
				},
			},
			want:    reconcile.Result{},
			wantErr: false,
			localObjects: []runtime.Object{
				testutil.NewSplunkForwarderCR().Build(),
				testutil.NewSplunkAuthSecret(),
				testutil.NewSplunkForwarderDaemonSet(),
				testutil.NewSplunkHECSecret(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := fakekubeclient.NewClientBuilder().WithScheme(scheme.Scheme).WithRuntimeObjects(tt.localObjects...).Build()
			r := &SecretReconciler{
				Client: fakeClient,
				Scheme: scheme.Scheme,
			}
			got, err := r.Reconcile(context.Background(), tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecretReconciler.Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecretReconciler.Reconcile() = %v, want %v", got, tt.want)
			}
		})
	}
}

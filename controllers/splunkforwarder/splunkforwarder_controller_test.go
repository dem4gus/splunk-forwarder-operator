package splunkforwarder

import (
	"context"
	"reflect"
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	sfv1alpha1 "github.com/openshift/splunk-forwarder-operator/api/v1alpha1"
	"github.com/openshift/splunk-forwarder-operator/internal/testutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	fakekubeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// TODO: tests should also check the reconciliation side-effects
// ie. making sure objects get created or modified properly

func TestReconcileSplunkForwarder_Reconcile(t *testing.T) {
	if err := sfv1alpha1.AddToScheme(scheme.Scheme); err != nil {
		t.Errorf("ReconcileSplunkForwarder.Reconcile() error = %v", err)
		return
	}
	if err := configv1.AddToScheme(scheme.Scheme); err != nil {
		t.Errorf("ReconcileSplunkForwarder.Reconcile() error = %v", err)
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
			name: "Reconcile succeeds when SplunkForwarder CR does not exist",
			args: args{
				request: reconcile.Request{},
			},
			want:         reconcile.Result{},
			wantErr:      false,
			localObjects: []runtime.Object{},
		},
		{
			name: "Reconcile fails when required splunk-auth secret is missing",
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      testutil.InstanceName,
						Namespace: testutil.InstanceNamespace,
					},
				},
			},
			want:    reconcile.Result{},
			wantErr: true,
			localObjects: []runtime.Object{
				testutil.NewSplunkForwarderCR().Build(),
			},
		},
		{
			name: "Reconcile requeues when HEC token secret is present",
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      testutil.InstanceName,
						Namespace: testutil.InstanceNamespace,
					},
				},
			},
			want: reconcile.Result{
				Requeue: true,
			},
			wantErr: false,
			localObjects: []runtime.Object{
				testutil.NewSplunkForwarderCR().Build(),
				testutil.NewSplunkForwarderService(),
				testutil.NewSplunkAuthSecret(),
				testutil.NewSplunkHECSecret(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := fakekubeclient.NewClientBuilder().WithScheme(scheme.Scheme).WithRuntimeObjects(tt.localObjects...).Build()
			r := &SplunkForwarderReconciler{
				Client:    fakeClient,
				Scheme:    scheme.Scheme,
				ReqLogger: log.WithValues(),
			}
			got, err := r.Reconcile(context.TODO(), tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileSplunkForwarder.Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReconcileSplunkForwarder.Reconcile() = %v, want %v", got, tt.want)
			}
		})
	}
}

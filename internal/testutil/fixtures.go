package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	sfv1alpha1 "github.com/openshift/splunk-forwarder-operator/api/v1alpha1"
	"github.com/openshift/splunk-forwarder-operator/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SplunkForwarderCRBuilder provides a fluent API for building test SplunkForwarder CRs.
// This builder merges functionality from the old SplunkForwarderInstance(useDigest bool) helper.
type SplunkForwarderCRBuilder struct {
	cr *sfv1alpha1.SplunkForwarder
}

// NewSplunkForwarderCR creates a new builder with sensible defaults.
// Use the fluent API methods to customize the CR before calling Build().
//
// Example:
//
//	cr := testutil.NewSplunkForwarderCR().
//		WithName("my-test").
//		WithImageTag("1.0.0").
//		Build()
func NewSplunkForwarderCR() *SplunkForwarderCRBuilder {
	return &SplunkForwarderCRBuilder{
		cr: &sfv1alpha1.SplunkForwarder{
			TypeMeta: metav1.TypeMeta{
				Kind:       "SplunkForwarder",
				APIVersion: "splunkforwarder.managed.openshift.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      InstanceName,
				Namespace: InstanceNamespace,
			},
			Spec: sfv1alpha1.SplunkForwarderSpec{
				SplunkLicenseAccepted: true,
				Image:                 Image,
				ImageTag:              ImageTag,
				SplunkInputs: []sfv1alpha1.SplunkForwarderInputs{
					{
						Path: TestLogPath,
					},
				},
			},
		},
	}
}

// WithName overrides the instance name.
func (b *SplunkForwarderCRBuilder) WithName(name string) *SplunkForwarderCRBuilder {
	b.cr.Name = name
	return b
}

// WithNamespace overrides the namespace.
func (b *SplunkForwarderCRBuilder) WithNamespace(namespace string) *SplunkForwarderCRBuilder {
	b.cr.Namespace = namespace
	return b
}

// WithGeneration sets the generation number.
func (b *SplunkForwarderCRBuilder) WithGeneration(gen int64) *SplunkForwarderCRBuilder {
	b.cr.Generation = gen
	return b
}

// WithSplunkInputs sets custom splunk inputs.
func (b *SplunkForwarderCRBuilder) WithSplunkInputs(inputs []sfv1alpha1.SplunkForwarderInputs) *SplunkForwarderCRBuilder {
	b.cr.Spec.SplunkInputs = inputs
	return b
}

// WithImageTag sets the ImageTag field.
// Note: ImageTag and ImageDigest can coexist in the spec; digest takes priority if both are present.
// This replaces the old SplunkForwarderInstance(false) pattern.
func (b *SplunkForwarderCRBuilder) WithImageTag(tag string) *SplunkForwarderCRBuilder {
	b.cr.Spec.ImageTag = tag
	return b
}

// WithImageDigest sets the ImageDigest field.
// Note: ImageTag and ImageDigest can coexist in the spec; digest takes priority if both are present.
// This replaces the old SplunkForwarderInstance(true) pattern.
func (b *SplunkForwarderCRBuilder) WithImageDigest(digest string) *SplunkForwarderCRBuilder {
	b.cr.Spec.ImageDigest = digest
	return b
}

// Build returns the constructed SplunkForwarder CR.
func (b *SplunkForwarderCRBuilder) Build() *sfv1alpha1.SplunkForwarder {
	return b.cr
}

// NewSplunkAuthSecret creates a splunk-auth secret with deterministic timestamp.
// Uses FixedTestTime instead of time.Now() for test reliability.
func NewSplunkAuthSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.SplunkAuthSecretName,
			Namespace: InstanceNamespace,
			CreationTimestamp: metav1.Time{
				Time: FixedTestTime,
			},
		},
	}
}

// NewSplunkHECSecret creates a splunk-hec-token secret with deterministic timestamp.
// Uses FixedTestTime instead of time.Now() for test reliability.
func NewSplunkHECSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.SplunkHECTokenSecretName,
			Namespace: InstanceNamespace,
			CreationTimestamp: metav1.Time{
				Time: FixedTestTime,
			},
		},
	}
}

// NewSplunkForwarderService creates a test service with deterministic timestamp.
// Uses FixedTestTime for test reliability.
func NewSplunkForwarderService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceName,
			Namespace: InstanceNamespace,
			CreationTimestamp: metav1.Time{
				Time: FixedTestTime,
			},
		},
	}
}

// NewSplunkForwarderDaemonSet creates a test DaemonSet with deterministic timestamp.
// Uses FixedTestTime for test reliability.
func NewSplunkForwarderDaemonSet() *appsv1.DaemonSet {
	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceName + "-ds",
			Namespace: InstanceNamespace,
			CreationTimestamp: metav1.Time{
				Time: FixedTestTime,
			},
		},
	}
}

// NewConfigMapVolume creates a ConfigMap volume for testing.
// This helper reduces verbosity when constructing volume arrays in tests.
func NewConfigMapVolume(name string) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: name,
				},
			},
		},
	}
}

// NewHostPathVolume creates a HostPath volume for testing.
// This helper reduces verbosity when constructing volume arrays in tests.
// The path type is always set to HostPathDirectory.
func NewHostPathVolume(name, path string) corev1.Volume {
	hostPathDirectory := corev1.HostPathDirectory
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: path,
				Type: &hostPathDirectory,
			},
		},
	}
}

// NewSecretVolume creates a Secret volume for testing.
// This helper reduces verbosity when constructing volume arrays in tests.
func NewSecretVolume(name, secretName string) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	}
}

// DeepEqualWithDiff compares two values of the same type and fails the test
// with a detailed diff if they differ.
//
// This generic function works with any type T, including runtime.Objects,
// slices, maps, and structs. It uses cmp.Diff for better error messages
// than reflect.DeepEqual.
//
// Example:
//
//	testutil.DeepEqualWithDiff(t, expectedService, actualService)
//	testutil.DeepEqualWithDiff(t, expectedVolumes, actualVolumes)
func DeepEqualWithDiff[T any](t *testing.T, expected, actual T) {
	t.Helper()
	diff := cmp.Diff(expected, actual)
	if diff != "" {
		t.Fatal("Objects differ: -expected, +actual\n", diff)
	}
}

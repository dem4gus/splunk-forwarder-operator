// Package kube test utilities.
//
// This file contains test helpers and fixtures for kube package tests.
// It should only be referenced from *_test.go files in this package.
package kube

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	sfv1alpha1 "github.com/openshift/splunk-forwarder-operator/api/v1alpha1"
)

const (
	instanceName      = "test"
	instanceNamespace = "openshift-test"
	image             = "test-image"
	imageTag          = "0.0.1"
	imageDigest       = "sha256:2452a3f01e840661ee1194777ed5a9185ceaaa9ec7329ed364fa2f02be22a701"
)

// splunkForwarderInstance returns a SplunkForwarder CR for testing.
//
// If useDigest is true, the returned CR will have ImageDigest set.
// If useDigest is false, the returned CR will have ImageTag set instead.
func splunkForwarderInstance(useDigest bool) *sfv1alpha1.SplunkForwarder {
	spec := sfv1alpha1.SplunkForwarderSpec{
		SplunkLicenseAccepted:  true,
		Image:                  image,
	}
	if useDigest {
		spec.ImageDigest = imageDigest
	} else {
		spec.ImageTag = imageTag
	}
	return &sfv1alpha1.SplunkForwarder{
		ObjectMeta: metav1.ObjectMeta{
			Name:       instanceName,
			Namespace:  instanceNamespace,
			Generation: 10,
		},
		Spec: spec,
	}
}

// deepEqualWithDiff compares two runtime.Object instances and fails the test
// with a detailed diff if they differ.
func deepEqualWithDiff(t *testing.T, expected, actual runtime.Object) {
	t.Helper()
	diff := cmp.Diff(expected, actual)
	if diff != "" {
		t.Fatal("Objects differ: -expected, +actual\n", diff)
	}
}

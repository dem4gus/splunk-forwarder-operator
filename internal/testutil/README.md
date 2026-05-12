# testutil Package

The `testutil` package provides centralized test utilities for the splunk-forwarder-operator test suite, including constants, fixture builders, and assertion helpers.

## Purpose

This package eliminates duplication across test files by providing:

* **Deterministic test fixtures** - All timestamps use `FixedTestTime` instead of `time.Now()`
* **Centralized constants** - Single source of truth for instance names, image references, ports, and paths
* **Reusable builders** - Fluent API for constructing test objects
* **Type-safe assertions** - Generic `DeepEqualWithDiff` for clear failure messages

## Package Structure

```
internal/testutil/
├── constants.go    # Shared test constants and fixed timestamp
├── fixtures.go     # Fixture builders and assertion helpers
└── README.md       # This file
```

## Usage

### Constants

```go
import "github.com/openshift/splunk-forwarder-operator/internal/testutil"

// Instance metadata
testutil.InstanceName      // "test"
testutil.InstanceNamespace // "openshift-test"

// Image references
testutil.Image       // "test-image"
testutil.ImageTag    // "0.0.1"
testutil.ImageDigest // "sha256:2452a3f01e840661ee1194777ed5a9185ceaaa9ec7329ed364fa2f02be22a701"

// Paths and ports
testutil.SplunkPort      // 9997
testutil.SplunkStatePath // "/var/lib/misc"
testutil.HostRootPath    // "/"
testutil.TestLogPath     // "/var/log/test"

// Deterministic timestamp
testutil.FixedTestTime // time.Date(2019, 12, 1, 12, 12, 0, 0, time.UTC)
```

### SplunkForwarder CR Builder

The builder pattern provides a fluent API for constructing test CRs:

```go
// Basic CR with image tag
cr := testutil.NewSplunkForwarderCR().
    WithImageTag(testutil.ImageTag).
    Build()

// CR with image digest
cr := testutil.NewSplunkForwarderCR().
    WithImageDigest(testutil.ImageDigest).
    WithGeneration(10).
    Build()

// Custom configuration
cr := testutil.NewSplunkForwarderCR().
    WithName("custom-name").
    WithNamespace("custom-namespace").
    WithSplunkInputs([]sfv1alpha1.SplunkForwarderInputs{
        {Path: "/custom/path"},
    }).
    Build()
```

### Kubernetes Object Fixtures

```go
// Secrets
authSecret := testutil.NewSplunkAuthSecret()
hecSecret := testutil.NewSplunkHECSecret()

// Service and DaemonSet
service := testutil.NewSplunkForwarderService()
daemonset := testutil.NewSplunkForwarderDaemonSet()
```

### Volume Builders

```go
// ConfigMap volume
configVol := testutil.NewConfigMapVolume("my-config")

// HostPath volume (always uses HostPathDirectory type)
hostVol := testutil.NewHostPathVolume("host", "/")

// Secret volume
secretVol := testutil.NewSecretVolume("splunk-auth", "splunk-auth")

// Use in test expectations
expectedVolumes := []corev1.Volume{
    testutil.NewConfigMapVolume("osd-monitored-logs-local"),
    testutil.NewHostPathVolume("host", testutil.HostRootPath),
    testutil.NewSecretVolume("splunk-auth", "splunk-auth"),
}
```

### Assertions

The generic `DeepEqualWithDiff` function works with any type and provides detailed diff output on failure:

```go
// Works with runtime.Objects
testutil.DeepEqualWithDiff(t, expectedService, actualService)

// Works with slices
testutil.DeepEqualWithDiff(t, expectedVolumes, actualVolumes)

// Works with any comparable type
testutil.DeepEqualWithDiff(t, expectedConfig, actualConfig)
```

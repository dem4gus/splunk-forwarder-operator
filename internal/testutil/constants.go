// Package testutil provides centralized test utilities, constants, and fixtures
// for the splunk-forwarder-operator test suite.
//
// This package eliminates code duplication across test files and provides
// deterministic test data through fixed constants and timestamps.
package testutil

import "time"

const (
	// Test instance metadata
	InstanceName      = "test"
	InstanceNamespace = "openshift-test"

	// Container image references
	Image       = "test-image"
	ImageTag    = "0.0.1"
	ImageDigest = "sha256:2452a3f01e840661ee1194777ed5a9185ceaaa9ec7329ed364fa2f02be22a701"

	// Test log path
	TestLogPath = "/var/log/test"

	// Splunk configuration
	SplunkPort      = 9997
	SplunkStatePath = "/var/lib/misc"
	HostRootPath    = "/"
)

var (
	// FixedTestTime provides deterministic timestamp for tests.
	// Using a fixed date: 2019-12-01 12:12:00 UTC.
	// This replaces non-deterministic time.Now() calls in test fixtures.
	FixedTestTime = time.Date(2019, 12, 1, 12, 12, 0, 0, time.UTC)
)

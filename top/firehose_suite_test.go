package top_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Top Suite")
}

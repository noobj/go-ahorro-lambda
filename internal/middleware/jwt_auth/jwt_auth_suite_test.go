package jwt_auth_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJwtAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JwtAuth Suite")
}

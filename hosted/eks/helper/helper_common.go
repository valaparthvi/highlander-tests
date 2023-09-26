package helper

import (
	. "github.com/onsi/gomega"
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	"github.com/rancher/rancher/tests/framework/extensions/cloudcredentials"
	"github.com/rancher/rancher/tests/framework/extensions/cloudcredentials/aws"
	"github.com/rancher/rancher/tests/framework/pkg/session"
)

type Context struct {
	CloudCred     *cloudcredentials.CloudCredential
	RancherClient *rancher.Client
	Session       *session.Session
}

// ContextOpts is set of options that can be set while creating a Context
type ContextOpts struct {
	// NoCreateCred if set to true, will not create cloud credentials
	NoCreateCred bool
}

func CommonBeforeEach(ctxopt ContextOpts) Context {
	testSession := session.NewSession()

	rancherClient, err := rancher.NewClient("", testSession)
	Expect(err).To(BeNil())

	var cloudCredential *cloudcredentials.CloudCredential
	if !ctxopt.NoCreateCred {
		cloudCredential, err = aws.CreateAWSCloudCredentials(rancherClient)
		Expect(err).To(BeNil())
	}

	return Context{
		CloudCred:     cloudCredential,
		RancherClient: rancherClient,
		Session:       testSession,
	}
}

func CommonAfterEach(ctx Context) {
	//	Delete created cloud creds
	//	Delete created cluster
}

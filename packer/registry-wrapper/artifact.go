package registrywrapper

import (
	"fmt"
)

const BuilderId = "packer.post-processor.packer-registry"

type Artifact struct {
	BucketSlug  string
	IterationID string
	BuildName   string
}

func (a *Artifact) BuilderId() string {
	return BuilderId
}

func (*Artifact) Id() string {
	return ""
}

func (a *Artifact) Files() []string {
	return []string{}
}

func (a *Artifact) String() string {
	return fmt.Sprintf("Published metadata to HCP Packer registry packer/%s/iterations/%s", a.BucketSlug, a.IterationID)
}

func (*Artifact) State(name string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	return nil
}

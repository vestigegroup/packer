package hcppackerregistry

import (
	packerregistry "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
)

const BuilderId = "packer.post-processor.hcp_packer_registry"

type Artifact struct {
	BuildName      string            `json:"name"`
	BuilderType    string            `json:"builder_type"`
	BuildTime      int64             `json:"build_time,omitempty"`
	Labels         map[string]string `json:"labels"`
	RegistryImages []packerregistry.Image
}

func (a *Artifact) BuilderId() string {
	return BuilderId
}

func (a *Artifact) Files() []string {
	return nil
}

func (a *Artifact) Id() string {
	return ""
}

func (a *Artifact) String() string {
	return ""
}

func (a *Artifact) State(name string) interface{} {
	if name != packerregistry.ArtifactStateURI {
		return nil
	}

	var images []packerregistry.Image
	for _, image := range a.RegistryImages {
		for k, v := range a.Labels {
			image.Labels[k] = v
		}
		images = append(images, image)
	}

	return images
}

func (a *Artifact) Destroy() error {
	return nil
}

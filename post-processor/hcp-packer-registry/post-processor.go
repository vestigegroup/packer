//go:generate packer-sdc mapstructure-to-hcl2 -type Config
//go:generate packer-sdc struct-markdown

package hcppackerregistry

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	packerregistry "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Labels              map[string]string `mapstructure:"labels"`
	ctx                 interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec { return p.config.FlatMapstructure().HCL2Spec() }

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		PluginType:         "packer.post-processor.hcp_packer_registry",
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	if len(p.config.Labels) == 0 {
		return fmt.Errorf("Error at least one key/value pair must be defined")
	}

	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packersdk.Ui, source packersdk.Artifact) (packersdk.Artifact, bool, bool, error) {
	generatedData := source.State("generated_data")
	if generatedData == nil {
		// Make sure it's not a nil map so we can assign to it later.
		generatedData = make(map[string]interface{})
	}
	p.config.ctx.Data = generatedData

	for key, data := range p.config.Labels {
		interpolatedData, err := createInterpolatedLabels(&p.config, data)
		if err != nil {
			return nil, false, false, err
		}
		p.config.Labels[key] = interpolatedData
	}

	var images []packerregistry.Image
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &images,
		WeaklyTypedInput: true,
		ErrorUnused:      true,
	})
	if err != nil {
		return source, true, true, nil
	}

	state := source.State(packerregistry.ArtifactStateURI)
	err = decoder.Decode(state)
	if err != nil {
		return source, true, true, nil
	}

	artifact := &Artifact{
		Labels:         p.config.Labels,
		RegistryImages: images,
	}

	return artifact, true, true, nil
}

func createInterpolatedLabels(config *Config, customData string) (string, error) {
	interpolatedCmd, err := interpolate.Render(customData, &config.ctx)
	if err != nil {
		return "", fmt.Errorf("Error interpolating custom data: %s", err)
	}
	return interpolatedCmd, nil
}

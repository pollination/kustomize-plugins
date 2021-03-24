package main

import (
	"context"
	"encoding/json"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
	"sigs.k8s.io/kustomize/api/kv"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

type SecretSource struct {
	ProjectID string `json:"projectId" yaml:"projectId"`
	Name      string `json:"name" yaml:"name"`
	Version   string `json:"version,omitempty" yaml:"version,omitempty"`
}

// A simple generator example.  Makes one service.
type plugin struct {
	h                     *resmap.PluginHelpers
	rf                    *resmap.Factory
	Metadata              *types.ObjectMeta
	DisableNameSuffixHash bool         `json:"disableNameSuffixHash,omitempty" yaml:"disableNameSuffixHash,omitempty"`
	Type                  string       `json:"type,omitempty" yaml:"type,omitempty"`
	Behavior              string       `json:"behavior,omitempty" yaml:"behavior,omitempty"`
	Source                SecretSource `json:"source" yaml:"source"`
	Keys                  []string     `json:"keys,omitempty" yaml:"keys,omitempty"`
}

var KustomizePlugin plugin

func (p *plugin) Config(h *resmap.PluginHelpers, config []byte) error {
	p.h = h
	p.rf = h.ResmapFactory()
	return yaml.Unmarshal(config, p)
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	secret, err := p.getSecret()
	if err != nil {
		return nil, err
	}
	return p.makeKubeSecret(secret)
}

func (p *plugin) getSecret() (secret []byte, err error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	secretVersion := "latest"

	if p.Source.Version != "" {
		secretVersion = p.Source.Version
	}

	secretFullName := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", p.Source.ProjectID, p.Source.Name, secretVersion)
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretFullName,
	}
	secretResponse, err := client.AccessSecretVersion(ctx, req)

	if err != nil {
		return nil, err
	}

	return secretResponse.GetPayload().Data, nil
}

func (p *plugin) stringInKeys(str string) bool {
	for _, k := range p.Keys {
		if str == k {
			return true
		}
	}
	return false
}

func (p *plugin) makeKubeSecret(secret []byte) (resmap.ResMap, error) {
	secretMap := make(map[string]string)

	err := yaml.Unmarshal(secret, &secretMap)
	if err != nil {
		err = json.Unmarshal(secret, &secretMap)
		if err != nil {
			return nil, fmt.Errorf("Could not unmarshal secret as a YAML or a JSON")
		}
	}

	options := &types.GeneratorOptions{
		DisableNameSuffixHash: p.DisableNameSuffixHash,
	}

	arg := types.SecretArgs{
		Type: p.Type,
		GeneratorArgs: types.GeneratorArgs{
			Name:      p.Metadata.Name,
			Namespace: p.Metadata.Namespace,
			Behavior:  p.Behavior,
			Options:   options,
		},
	}

	for key, value := range secretMap {
		if p.stringInKeys(key) {
			arg.LiteralSources = append(arg.LiteralSources, key+"="+value)
		}
	}

	return p.rf.FromSecretArgs(
		kv.NewLoader(p.h.Loader(), p.h.Validator()), arg)

}

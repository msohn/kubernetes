/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"github.com/spf13/cobra"

	"k8s.io/kubernetes/pkg/printers"
)

// kubectlConfigPrintFlags composes common printer flag structs
// used across all config commands, and provides a method
// of retrieving a known printer based on flag values provided.
type kubectlConfigPrintFlags struct {
	JSONYamlPrintFlags *printers.JSONYamlPrintFlags
	NamePrintFlags     *printers.NamePrintFlags
	TemplateFlags      *printers.KubeTemplatePrintFlags

	OutputFormat *string
}

func (f *kubectlConfigPrintFlags) Complete(successTemplate string) error {
	return f.NamePrintFlags.Complete(successTemplate)
}

func (f *kubectlConfigPrintFlags) ToPrinter() (printers.ResourcePrinter, error) {
	outputFormat := ""
	if f.OutputFormat != nil {
		outputFormat = *f.OutputFormat
	}

	if p, err := f.JSONYamlPrintFlags.ToPrinter(outputFormat); !printers.IsNoCompatiblePrinterError(err) {
		return p, err
	}

	if p, err := f.NamePrintFlags.ToPrinter(outputFormat); !printers.IsNoCompatiblePrinterError(err) {
		return p, err
	}

	if p, err := f.TemplateFlags.ToPrinter(outputFormat); !printers.IsNoCompatiblePrinterError(err) {
		return p, err
	}

	return nil, printers.NoCompatiblePrinterError{Options: f}
}

func (f *kubectlConfigPrintFlags) AddFlags(cmd *cobra.Command) {
	f.JSONYamlPrintFlags.AddFlags(cmd)
	f.NamePrintFlags.AddFlags(cmd)
	f.TemplateFlags.AddFlags(cmd)

	if f.OutputFormat != nil {
		cmd.Flags().StringVarP(f.OutputFormat, "output", "o", *f.OutputFormat, "Output format. One of: json|yaml|name|go-template=...|go-template-file=...|jsonpath=...|jsonpath-file=... See custom columns [http://kubernetes.io/docs/user-guide/kubectl-overview/#custom-columns], golang template [http://golang.org/pkg/text/template/#pkg-overview] and jsonpath template [http://kubernetes.io/docs/user-guide/jsonpath].")
	}
}

// WithDefaultOutput sets a default output format if one is not provided through a flag value
func (f *kubectlConfigPrintFlags) WithDefaultOutput(output string) *kubectlConfigPrintFlags {
	f.OutputFormat = &output
	return f
}

func newKubeConfigPrintFlags(operation string) *kubectlConfigPrintFlags {
	outputFormat := ""

	return &kubectlConfigPrintFlags{
		OutputFormat: &outputFormat,

		JSONYamlPrintFlags: printers.NewJSONYamlPrintFlags(),
		NamePrintFlags:     printers.NewNamePrintFlags(operation),
		TemplateFlags:      printers.NewKubeTemplatePrintFlags(),
	}
}

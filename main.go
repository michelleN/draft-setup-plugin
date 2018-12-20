package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const globalUsage = `
Set up Draft and all dependencies for developer environment.
This requires Minikube, Helm and, Draft to all be installed
on your machine.
	$ helm setup
`

const defaultKubernetesVersion = "v1.11.6"

type setupCmd struct {
	kubernetesVersion string
	helmVersion       string
}

func main() {

	s := &setupCmd{}
	cmd := &cobra.Command{
		Use:   "template [flags] CHART",
		Short: "set up developer environment for Draft",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&s.kubernetesVersion, "kubernetes-version", "k", "v1.11.6", "kubernetes version")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (s *setupCmd) run() error {
	kv := s.kubernetesVersion
	if kv == "" {
		kv = defaultKubernetesVersion
	}

	//start minikube
	var c *exec.Cmd
	c = exec.Command("minikube", "start", "--kubernetes-version", kv)
	c.Stderr = os.Stderr
	fmt.Println("Preparing a development cluster...")
	_, err := c.Output()
	if err != nil {
		return errors.New("Unable to provision local development cluster")
	}
	fmt.Println("Successfully prepared a development cluster.")

	//install helm
	c = exec.Command("helm", "init")
	c.Stderr = os.Stderr
	fmt.Println("Installing Helm...")
	_, err = c.Output()
	if err != nil {
		return errors.New("Unable to install Helm")
	}
	fmt.Println("Successfully installed Helm.")

	//run a draft init
	c = exec.Command("draft", "init")
	c.Stderr = os.Stderr
	fmt.Println("Configuring Draft...")
	_, err = c.Output()
	if err != nil {
		return err
	}
	fmt.Println("Successfully configured Draft")

	return nil
}

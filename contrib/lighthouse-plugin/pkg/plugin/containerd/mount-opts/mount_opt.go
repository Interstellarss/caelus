package mount_opts

import (
	"context"
	"github.com/Interstellarss/lighthouse-plugin/pkg/plugin"
	containerdPlugin "github.com/Interstellarss/lighthouse-plugin/pkg/plugin/containerd"
	"github.com/docker/docker/api/types/container"
	"k8s.io/klog"
	"log"

	"os/exec"
)

func init() {
	//
	containerdPlugin.PreHookUpdateContainer.RegisterSubPlugin(opt.handle)
}

var (
	opt = &mountOpt{}
)

type mountOpt struct {
}

func (p *mountOpt) handle(
	toolKits *containerdPlugin.ToolKits,
	//not sure if we need this part
	updateConfig *container.UpdateConfig,
	metadata *plugin.PodMetadata,
	containerStatus string) error {
	//plugin containerd mount opts
	//if !utilfeature.DefaultFeatureGate.Enabled(plugin.cont)

	if metadata.ContainerType != plugin.ContainerTypeLabelContainer {

	}

	//TODO: change this to json?
	containerSpec, err := toolKits.ContainerdClient.LoadContainer(context.TODO(), metadata.SandBoxID)

	if err != nil {

	}

	//info , err := containerSpec.Info()

	//mount new directory for the container :)
	//testing for mound

	//task, err := containerSpec.NewTask(context.Background(), cio.NewCreator(cio.WithStdio))

	//task, err := containerSpec

	//containerSpec.NewTask()

	out, err := exec.Command("/bin/bash", "mount.sh", containerSpec.ID(), metadata.FunctionName).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}

	klog.Infof("out put of the mounting %v", out)

	return nil

}

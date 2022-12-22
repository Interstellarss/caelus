package pause_opts

import (
	"context"
	"github.com/Interstellarss/lighthouse-plugin/pkg/plugin"
	containerdPlugin "github.com/Interstellarss/lighthouse-plugin/pkg/plugin/containerd"
	"github.com/containerd/containerd/cio"
	"syscall"

	//"github.com/containerd/containerd/namespaces"
	"github.com/docker/docker/api/types/container"
	//"github.com/opencontainers/runtime-spec/specs-go"
)

func init() {

	containerdPlugin.PreHookUpdateContainer.RegisterSubPlugin(opt.handle)
}

var (
	opt = &pauseOpt{}
)

type pauseOpt struct {
}

func (p *pauseOpt) handle(
	toolKits *containerdPlugin.ToolKits,
	updateConfig *container.UpdateConfig,
	metadata *plugin.PodMetadata,
	containerStatus string) error {
	if metadata.ContainerType != plugin.ContainerTypeLabelContainer {

	}

	//TODO: implementing the pause mechanism (send signals? or directly pause?)
	info, err := toolKits.ContainerdClient.LoadContainer(context.Background(), metadata.ContainerName)

	if err != nil {

	}

	spec, err := info.Spec(context.Background())

	//TODO:logic of checking whether need to pause or not

	//get the actual task of the container
	task, err := info.Task(context.Background(), cio.NewAttach(cio.WithStdio))
	if err != nil {

	}

	//if it is moved outside the pool
	//TODO: check this part
	if spec.Annotations["pre-pool"] == "true" {
		task.Pause(context.Background())
	} else {
		// Create a new stdio creator that uses the current process's stdio streams
		stdioAttatch := cio.NewAttach(cio.WithStdio)

		// Set the namespace for the stdio creator
		//stdioAttatch = namespaces.WithNamespace(context.Background(), "my-namespace")

		// Load a new process inside the task with the specified command and arguments
		// TODO: need to define a new default process id for the process
		process, err := task.LoadProcess(context.Background(),
			"<process_id>", stdioAttatch)
		if err != nil {
			// handle error
		}

		//TODO: Define different SIGNALs for different signals!
		if spec.Annotations["pre-pool"] == "true" {
			err = process.Kill(context.Background(), syscall.SIGSTOP)
			if err != nil {

			}
		} else if spec.Annotations["pre-pool"] == "false" {
			//TODO: defining new siganls maybe?
			err = process.Kill(context.Background(), syscall.SIGALRM)
		} else {
			return nil
		}

		//os.Signal()
		//syscall.si
	}

	//task.

	return nil
}

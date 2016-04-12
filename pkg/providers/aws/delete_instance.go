package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/emc-advanced-dev/unik/pkg/types"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func (p *AwsProvider) DeleteInstance(id string) error {
	instance, err := p.GetInstance(id)
	if err != nil {
		return lxerrors.New("retrieving instance "+id, err)
	}
	param := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instance.Id),
		},
	}
	_, err = p.newEC2().TerminateInstances(param)
	if err != nil {
		return lxerrors.New("failed to terminate instance "+instance.Id, err)
	}
	err = p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
		delete(instances, instance.Id)
		return nil
	})
	if err != nil {
		return lxerrors.New("modifying image map in state", err)
	}
	err = p.state.Save()
	if err != nil {
		return lxerrors.New("saving image to state", err)
	}
	return nil
}

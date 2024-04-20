package group

import (
	"context"
	"log"

	"github.com/turao/topics/metadata"
	apiV1 "github.com/turao/topics/users/api/v1"
	"github.com/turao/topics/users/entity/group"
)

type GroupRepository interface {
	Save(ctx context.Context, group group.Group) error
	FindByID(ctx context.Context, groupID group.ID) (group.Group, error)
}

type service struct {
	groupRepository GroupRepository
}

var _ apiV1.Groups = (*service)(nil)

func NewService(
	groupRepository GroupRepository,
) (*service, error) {
	return &service{
		groupRepository: groupRepository,
	}, nil
}

// RegisterGroup implements apiV1.Groups
func (svc *service) CreateGroup(ctx context.Context, req apiV1.CreateGroupRequest) (apiV1.CreateGroupResponse, error) {
	log.Println("creating group", req)
	group, err := group.NewGroup(
		group.WithName(req.Name),
		group.WithTenancy(metadata.Tenancy(req.Tenancy)),
	)
	if err != nil {
		return apiV1.CreateGroupResponse{}, err
	}

	err = svc.groupRepository.Save(ctx, group)
	if err != nil {
		return apiV1.CreateGroupResponse{}, err
	}

	log.Println("group registered succesfully")
	return apiV1.CreateGroupResponse{
		ID: group.ID().String(),
	}, nil
}

func (svc *service) DeleteGroup(ctx context.Context, req apiV1.DeleteGroupRequest) (apiV1.DeleteGroupResponse, error) {
	log.Println("deleting group", req)
	group, err := svc.groupRepository.FindByID(ctx, group.ID(req.ID))
	if err != nil {
		return apiV1.DeleteGroupResponse{}, err
	}

	group.Delete()
	err = svc.groupRepository.Save(ctx, group)
	if err != nil {
		return apiV1.DeleteGroupResponse{}, err
	}

	log.Println("group deleted succesfully")
	return apiV1.DeleteGroupResponse{}, nil
}

func (svc *service) GetGroup(ctx context.Context, req apiV1.GetGroupRequest) (apiV1.GetGroupResponse, error) {
	group, err := svc.groupRepository.FindByID(ctx, group.ID(req.ID))
	if err != nil {
		return apiV1.GetGroupResponse{}, err
	}

	groupInfo, err := groupMapper.ToGroupInfo(group)
	if err != nil {
		return apiV1.GetGroupResponse{}, nil
	}

	return apiV1.GetGroupResponse{
		Group: groupInfo,
	}, nil
}

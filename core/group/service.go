package group

import (
	"context"
	"fmt"
	"strings"

	"github.com/odpf/shield/core/relation"
	"github.com/odpf/shield/core/user"
	"github.com/odpf/shield/internal/schema"
	"github.com/odpf/shield/pkg/uuid"
)

type RelationService interface {
	Create(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error)
	DeleteSubjectRelations(ctx context.Context, resourceType, optionalResourceID string) error
	ListRelations(ctx context.Context, rel relation.RelationV2) ([]relation.RelationV2, error)
}

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
	GetByID(ctx context.Context, id string) (user.User, error)
	GetByIDs(ctx context.Context, userIDs []string) ([]user.User, error)
}

type Service struct {
	repository      Repository
	relationService RelationService
	userService     UserService
}

func NewService(repository Repository, relationService RelationService, userService UserService) *Service {
	return &Service{
		repository:      repository,
		relationService: relationService,
		userService:     userService,
	}
}

func (s Service) Create(ctx context.Context, grp Group) (Group, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Group{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	newGroup, err := s.repository.Create(ctx, grp)
	if err != nil {
		return Group{}, err
	}

	// add relationship between group to org
	if err = s.CreateRelation(ctx, newGroup, relation.Subject{
		ID:        newGroup.OrganizationID,
		Namespace: schema.OrganizationNamespace,
		RoleID:    schema.OrganizationRelationName,
	}); err != nil {
		return Group{}, err
	}

	// attach group to org
	if err = s.addGroupAsMember(ctx, newGroup); err != nil {
		return Group{}, err
	}

	// attach current user to group as admin
	if err = s.CreateRelation(ctx, newGroup, BuildUserGroupAdminSubject(currentUser)); err != nil {
		return Group{}, err
	}

	return newGroup, nil
}

func (s Service) Get(ctx context.Context, idOrSlug string) (Group, error) {
	if uuid.IsValid(idOrSlug) {
		return s.repository.GetByID(ctx, idOrSlug)
	}
	return s.repository.GetBySlug(ctx, idOrSlug)
}

func (s Service) GetByIDs(ctx context.Context, groupIDs []string) ([]Group, error) {
	return s.repository.GetByIDs(ctx, groupIDs)
}

func (s Service) List(ctx context.Context, flt Filter) ([]Group, error) {
	return s.repository.List(ctx, flt)
}

func (s Service) Update(ctx context.Context, grp Group) (Group, error) {
	if strings.TrimSpace(grp.ID) != "" {
		return s.repository.UpdateByID(ctx, grp)
	}
	return s.repository.UpdateBySlug(ctx, grp)
}

func (s Service) ListUserGroups(ctx context.Context, userId string, roleId string) ([]Group, error) {
	return s.repository.ListUserGroups(ctx, userId, roleId)
}

func (s Service) ListGroupRelations(ctx context.Context, objectId, subjectType, role string) ([]user.User, []Group, map[string][]string, map[string][]string, error) {
	relationList, err := s.repository.ListGroupRelations(ctx, objectId, subjectType, role)
	if err != nil {
		return []user.User{}, []Group{}, map[string][]string{}, map[string][]string{}, fmt.Errorf("%w: %s", ErrListingGroupRelations, err.Error())
	}

	userIDs := []string{}
	groupIDs := []string{}
	userIDRoleMap := map[string][]string{}
	groupIDRoleMap := map[string][]string{}
	users := []user.User{}
	groups := []Group{}

	for _, relation := range relationList {
		if relation.Subject.Namespace == schema.UserPrincipal {
			userIDs = append(userIDs, relation.Subject.ID)
			userIDRoleMap[relation.Subject.ID] = append(userIDRoleMap[relation.Subject.ID], relation.Subject.RoleID)
		} else if relation.Subject.Namespace == schema.GroupPrincipal {
			groupIDs = append(groupIDs, relation.Subject.ID)
			groupIDRoleMap[relation.Subject.ID] = append(groupIDRoleMap[relation.Subject.ID], relation.Subject.RoleID)
		}
	}

	if len(userIDs) > 0 {
		userList, err := s.userService.GetByIDs(ctx, userIDs)
		if err != nil {
			return []user.User{}, []Group{}, map[string][]string{}, map[string][]string{}, fmt.Errorf("%w: %s", ErrFetchingUsers, err.Error())
		}

		users = append(users, userList...)
	}

	if len(groupIDs) > 0 {
		groupList, err := s.repository.GetByIDs(ctx, groupIDs)
		if err != nil {
			return []user.User{}, []Group{}, map[string][]string{}, map[string][]string{}, fmt.Errorf("%w: %s", ErrFetchingGroups, err.Error())
		}

		groups = append(groups, groupList...)
	}

	return users, groups, userIDRoleMap, groupIDRoleMap, nil
}

func (s Service) CreateRelation(ctx context.Context, team Group, subject relation.Subject) error {
	rel := relation.RelationV2{
		Object: relation.Object{
			ID:        team.ID,
			Namespace: schema.GroupNamespace,
		},
		Subject: subject,
	}
	if _, err := s.relationService.Create(ctx, rel); err != nil {
		return err
	}
	return nil
}

func (s Service) addGroupAsMember(ctx context.Context, team Group) error {
	rel := relation.RelationV2{
		Object: relation.Object{
			ID:        team.OrganizationID,
			Namespace: schema.OrganizationNamespace,
		},
		Subject: relation.Subject{
			ID:        team.ID,
			Namespace: schema.GroupNamespace,
			RoleID:    schema.MemberRole,
		},
	}

	_, err := s.relationService.Create(ctx, rel)
	if err != nil {
		return err
	}

	return nil
}

// ListByOrganization will be useful for nested groups but we don't do that at the moment
// so it will not be directly used
func (s Service) ListByOrganization(ctx context.Context, id string) ([]Group, error) {
	relations, err := s.relationService.ListRelations(ctx, relation.RelationV2{
		Object: relation.Object{
			Namespace: schema.GroupNamespace,
		},
		Subject: relation.Subject{
			ID:        id,
			Namespace: schema.OrganizationNamespace,
			RoleID:    schema.OrganizationRelationName,
		},
	})
	if err != nil {
		return nil, err
	}

	var groupIDs []string
	for _, rel := range relations {
		groupIDs = append(groupIDs, rel.Object.ID)
	}
	if len(groupIDs) == 0 {
		// no groups
		return []Group{}, nil
	}
	return s.repository.GetByIDs(ctx, groupIDs)
}

func (s Service) Enable(ctx context.Context, id string) error {
	return s.repository.SetState(ctx, id, Enabled)
}

func (s Service) Disable(ctx context.Context, id string) error {
	return s.repository.SetState(ctx, id, Disabled)
}

func (s Service) Delete(ctx context.Context, id string) error {
	if err := s.relationService.DeleteSubjectRelations(ctx, schema.GroupPrincipal, id); err != nil {
		return err
	}
	return s.repository.Delete(ctx, id)
}

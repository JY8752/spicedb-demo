package spicedb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SpiceDBClient interface {
	CheckPermission(ctx context.Context, request *CheckPermissionRequest) (bool, error)
	LookupResources(ctx context.Context, request *LookupResourcesRequest) ([]string, error)
	LookupSubjects(ctx context.Context, request *LookupSubjectsRequest) ([]string, error)
}

type spiceDBClient struct {
	client *authzed.Client
}

func NewSpiceDBClient(host string, port int, token string) (*spiceDBClient, error) {
	client, err := authzed.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken(token),
	)
	if err != nil {
		return nil, err
	}
	return &spiceDBClient{client: client}, nil
}

type CheckPermissionRequest struct {
	resource   ObjectType
	resourceId string
	subject    ObjectType
	subjectId  string
	permission Permission
}

func NewCheckPermissionRequest(resource ObjectType, resourceId string, subject ObjectType, subjectId string, permission Permission) *CheckPermissionRequest {
	return &CheckPermissionRequest{
		resource:   resource,
		resourceId: resourceId,
		subject:    subject,
		subjectId:  subjectId,
		permission: permission,
	}
}

func (c *spiceDBClient) CheckPermission(ctx context.Context, request *CheckPermissionRequest) (bool, error) {
	resource := newObjectReference(request.resource, request.resourceId)
	subject := newSubjectReference(request.subject, request.subjectId)

	resp, err := c.client.CheckPermission(ctx, &v1.CheckPermissionRequest{
		Resource:   resource,
		Permission: request.permission.String(),
		Subject:    subject,
	})
	if err != nil {
		return false, err
	}

	return resp.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION, nil
}

type LookupResourcesRequest struct {
	resourceObjectType ObjectType
	permission         Permission
	subject            ObjectType
	subjectId          string
}

func NewLookupResourcesRequest(resourceObjectType ObjectType, permission Permission, subject ObjectType, subjectId string) *LookupResourcesRequest {
	return &LookupResourcesRequest{
		resourceObjectType: resourceObjectType,
		permission:         permission,
		subject:            subject,
		subjectId:          subjectId,
	}
}

func (c *spiceDBClient) LookupResources(ctx context.Context, request *LookupResourcesRequest) (iter.Seq[string], error) {
	subject := newSubjectReference(request.subject, request.subjectId)

	resp, err := c.client.LookupResources(ctx, &v1.LookupResourcesRequest{
		ResourceObjectType: request.resourceObjectType.String(),
		Permission:         request.permission.String(),
		Subject:            subject,
	})
	if err != nil {
		return nil, err
	}

	return receive(func() (string, error) {
		r, err := resp.Recv()
		if err != nil {
			return "", err
		}
		return r.ResourceObjectId, nil
	}), nil
}

func receive(recv func() (string, error)) iter.Seq[string] {
	return func(yield func(string) bool) {
		for {
			r, err := recv()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				continue
			}

			if !yield(r) {
				break
			}
		}
	}
}

type LookupSubjectsRequest struct {
	resourceObjectType ObjectType
	resourceId         string
	permission         Permission
	subjectObjectType  ObjectType
}

func NewLookupSubjectsRequest(resourceObjectType ObjectType, resourceId string, permission Permission, subjectObjectType ObjectType) *LookupSubjectsRequest {
	return &LookupSubjectsRequest{
		resourceObjectType: resourceObjectType,
		resourceId:         resourceId,
		permission:         permission,
		subjectObjectType:  subjectObjectType,
	}
}

func (c *spiceDBClient) LookupSubjects(ctx context.Context, request *LookupSubjectsRequest) (iter.Seq[string], error) {
	resource := newObjectReference(request.resourceObjectType, request.resourceId)

	resp, err := c.client.LookupSubjects(ctx, &v1.LookupSubjectsRequest{
		Resource:          resource,
		Permission:        request.permission.String(),
		SubjectObjectType: request.subjectObjectType.String(),
	})
	if err != nil {
		return nil, err
	}

	return receive(func() (string, error) {
		r, err := resp.Recv()
		if err != nil {
			return "", err
		}
		return r.Subject.SubjectObjectId, nil
	}), nil
}

type Permission string

func (p Permission) String() string {
	return string(p)
}

const (
	ReadPermission Permission = "read"
)

type ObjectType string

const (
	UserObjectType ObjectType = "user"
	PostObjectType ObjectType = "post"
)

func (o ObjectType) String() string {
	return string(o)
}

func newObjectReference(objectType ObjectType, id string) *v1.ObjectReference {
	return &v1.ObjectReference{
		ObjectType: objectType.String(),
		ObjectId:   id,
	}
}

func newSubjectReference(objectType ObjectType, id string) *v1.SubjectReference {
	return &v1.SubjectReference{
		Object: newObjectReference(objectType, id),
	}
}

package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
)

var (
	ErrIdInvalid = errors.New("invalid id")
)

var _ ActorService = (*Service)(nil)

type Service struct {
	repo ActorRepository
}

func NewService(ar ActorRepository) *Service {
	return &Service{
		repo: ar,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]*ActorResponse, error) {
	const op = "actor.Service.GetAll"

	actors, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("ERROR: failed to get actor records from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]*ActorResponse, 0, len(actors))
	for _, v := range actors {
		res = append(res, ToActorResponse(v))
	}

	return res, nil
}

func (s *Service) Add(ctx context.Context, req *ActorInfo) (*ActorResponse, error) {
	const op = "actor.Service.Add"

	vErr := ValidateEmptyActorInfo(req)
	if vErr != nil {
		log.Printf("ERROR: failed request empty validation\n")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}
	vErr = ValidateFormatActorInfo(req)
	if vErr != nil {
		log.Printf("ERROR: failed request format validation\n")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}
	actor := ToActor(req)

	actor, err := s.repo.Add(ctx, actor)
	if err != nil {
		log.Printf("ERROR: failed to create actor record in repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorResponse(actor)

	return res, nil
}

func (s *Service) Get(ctx context.Context, req *ActorIdRequest) (*ActorResponse, error) {
	const op = "actor.Service.Get"

	id, err := strconv.ParseUint(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	actor, err := s.repo.Get(ctx, int(id))
	if err != nil {
		log.Printf("ERROR: failed to get actor record from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorResponse(actor)

	return res, nil
}

func (s *Service) Update(ctx context.Context, req *ActorIdInfoRequest) (*ActorResponse, error) {
	const op = "actor.Service.Update"

	id, err := strconv.ParseUint(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	vErr := ValidateFormatActorInfo(&req.Info)
	if vErr != nil {
		log.Printf("ERROR: failed request format validation\n")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}

	actor := ToActor(&req.Info)
	actor.ID = int(id)

	err = s.repo.Update(ctx, actor)
	if err != nil {
		log.Printf("ERROR: failed to update actor record in repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	actor, err = s.repo.Get(ctx, int(id))
	if err != nil {
		log.Printf("ERROR: failed to get actor record from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorResponse(actor)

	return res, nil
}

func (s *Service) Delete(ctx context.Context, req *ActorIdRequest) (*ActorResponse, error) {
	const op = "actor.Service.Delete"

	id, err := strconv.ParseUint(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	actor, err := s.repo.Get(ctx, int(id))
	if err != nil {
		log.Printf("ERROR: failed to get actor record from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.repo.Delete(ctx, int(id))
	if err != nil {
		log.Printf("ERROR: failed to delete actor record in repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorResponse(actor)

	return res, nil
}

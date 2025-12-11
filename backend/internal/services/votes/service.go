package votes

import (
	"context"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/hub"
	voterepository "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/votes"
	"github.com/google/uuid"
)

type VoteService interface {
	Add(context.Context, rooms.Vote) (rooms.Vote, error)
	Get(context.Context, string) (rooms.Vote, error)
	GetForRoom(context.Context, string) ([]rooms.Vote, error)
	Delete(context.Context, string, string) error
}

type Service struct {
	repo voterepository.VoteRepository
	hub  hub.Hub
}

func NewService(repo voterepository.VoteRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetHub(h hub.Hub) {
	s.hub = h
}

func (s *Service) Add(ctx context.Context, vote rooms.Vote) (rooms.Vote, error) {
	id, err := uuid.Parse(vote.ID)
	if err != nil {
		return rooms.Vote{}, err
	}

	roomID, err := uuid.Parse(vote.RoomID)
	if err != nil {
		return rooms.Vote{}, err
	}

	gameID, err := uuid.Parse(vote.GameID)
	if err != nil {
		return rooms.Vote{}, err
	}

	userID, err := uuid.Parse(vote.UserID)
	if err != nil {
		return rooms.Vote{}, err
	}

	result, err := s.repo.Add(ctx, voterepository.AddParams{
		ID:     id,
		RoomID: roomID,
		GameID: gameID,
		UserID: userID,
	})
	if err == nil && s.hub != nil {
		s.hub.Broadcast(vote.RoomID, hub.RoomEvent{
			Type:   hub.EventVoteAdded,
			RoomID: vote.RoomID,
			Payload: map[string]any{
				"id":      result.ID,
				"game_id": result.GameID,
				"user_id": result.UserID,
			},
		})
	}
	return result, err
}

func (s *Service) Get(ctx context.Context, id string) (rooms.Vote, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return rooms.Vote{}, err
	}

	return s.repo.Get(ctx, uuidID)
}

func (s *Service) GetForRoom(ctx context.Context, roomID string) ([]rooms.Vote, error) {
	uuidRoomID, err := uuid.Parse(roomID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetForRoom(ctx, uuidRoomID)
}

func (s *Service) Delete(ctx context.Context, id string, roomID string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(ctx, uuidID)
	if err == nil && s.hub != nil {
		s.hub.Broadcast(roomID, hub.RoomEvent{
			Type:   hub.EventVoteDeleted,
			RoomID: roomID,
			Payload: map[string]any{
				"id": id,
			},
		})
	}
	return err
}

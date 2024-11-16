package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"

	"github.com/kitae0522/gommunity/internal/dto"
	"github.com/kitae0522/gommunity/internal/model"
	"github.com/kitae0522/gommunity/internal/repository"
	"github.com/kitae0522/gommunity/pkg/exception"
)

type ThreadService struct {
	threadRepo *repository.ThreadRepository
	redisCache *redis.Client
}

func NewThreadService(repo *repository.ThreadRepository, rdconn *redis.Client) *ThreadService {
	return &ThreadService{
		threadRepo: repo,
		redisCache: rdconn,
	}
}

func (s *ThreadService) CreateThread(ctx context.Context, req *dto.CreateThreadRequest) (*model.ThreadModel, *exception.ErrResponseCtx) {
	thread, err := s.threadRepo.CreateThread(ctx, req)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			return nil, exception.GenerateErrorCtx(fiber.StatusNotFound, "❌ 쓰레드 생성 실패. 존재하지 않는 사용자입니다.", err)
		default:
			return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 생성 실패. Repository에서 문제가 발생했습니다.", err)
		}
	}

	txns := make([]model.PrismaTransaction, 0)
	if req.ParentThread != nil {
		txns = append(txns, s.threadRepo.LinkParentThread(ctx, thread.ID, *req.ParentThread))
	} else if req.NextThread != nil {
		txns = append(txns, s.threadRepo.LinkNextThread(ctx, thread.ID, *req.NextThread))
	} else if req.PrevThread != nil {
		txns = append(txns, s.threadRepo.LinkPrevThread(ctx, thread.ID, *req.NextThread))
	}

	if err := s.threadRepo.LinkRelation(ctx, txns); err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 생성 실패. Reposioty에서 문제가 발생했습니다.", err)
	}

	return thread, nil
}

func (s *ThreadService) ListThread(ctx context.Context) ([]model.ThreadModel, *exception.ErrResponseCtx) {
	threadList, err := s.threadRepo.ListThread(ctx)
	if err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. Repository에서 문제가 발생했습니다.", err)
	}
	return threadList, nil
}

func (s *ThreadService) ListThreadByHandle(ctx context.Context, handle string) ([]model.ThreadModel, *exception.ErrResponseCtx) {
	threadList, err := s.threadRepo.ListThreadByHandle(ctx, handle)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			return nil, exception.GenerateErrorCtx(fiber.StatusNotFound, "❌ 쓰레드 조회 실패. 존재하지 않는 사용자입니다.", err)
		default:
			return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. Repository에서 문제가 발생했습니다.", err)
		}
	}
	return threadList, nil
}

func (s *ThreadService) GetThreadByID(ctx context.Context, threadID int) (*model.ThreadModel, *exception.ErrResponseCtx) {
	thread, err := s.getThreadFromCache(ctx, threadID)
	if err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. 캐시하는 과정에서 문제가 발생했습니다.", err)
	}
	if thread != nil {
		return thread, nil
	}

	thread, err = s.threadRepo.GetThreadByID(ctx, threadID)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			return nil, exception.GenerateErrorCtx(fiber.StatusNotFound, "❌ 쓰레드 조회 실패. 존재하지 않는 쓰레드입니다.", err)
		default:
			return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. Repository에서 문제가 발생했습니다.", err)
		}
	}

	if err := s.setThreadToCache(ctx, thread, 5*time.Minute); err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. 캐시에 저장하지 못했습니다.", err)
	}

	return thread, nil
}

func (s *ThreadService) CommentsByID(ctx context.Context, threadID int) ([]model.ThreadModel, *exception.ErrResponseCtx) {
	comments, err := s.threadRepo.CommentsByID(ctx, threadID)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			return nil, exception.GenerateErrorCtx(fiber.StatusNotFound, "❌ 쓰레드 조회 실패. 존재하지 않는 쓰레드입니다.", err)
		default:
			return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. Repository에서 문제가 발생했습니다.", err)
		}
	}
	return comments, nil
}

func (s *ThreadService) getFromCache(ctx context.Context, key string, data interface{}) error {
	cachedData, err := s.redisCache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	return json.Unmarshal(cachedData, data)
}

func (s *ThreadService) setToCache(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.redisCache.Set(ctx, key, jsonData, ttl).Err()
}

func (s *ThreadService) getThreadFromCache(ctx context.Context, threadID int) (*model.ThreadModel, error) {
	var thread *model.ThreadModel
	err := s.getFromCache(ctx, fmt.Sprintf("thread:%d", threadID), thread)
	return thread, err
}

func (s *ThreadService) setThreadToCache(ctx context.Context, thread *model.ThreadModel, ttl time.Duration) error {
	return s.setToCache(ctx, fmt.Sprintf("thread:%d", thread.ID), thread, ttl)
}

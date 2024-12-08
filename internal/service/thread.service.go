package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"

	"github.com/kitae0522/gommunity/internal/config"
	"github.com/kitae0522/gommunity/internal/dto"
	"github.com/kitae0522/gommunity/internal/model"
	"github.com/kitae0522/gommunity/internal/repository"
	"github.com/kitae0522/gommunity/pkg/exception"
	"github.com/kitae0522/gommunity/pkg/utils"
)

type ThreadService struct {
	threadRepo *repository.ThreadRepository
	redisCache *redis.Client
	txnsItr    []model.PrismaTransaction
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

	if err := s.threadRepo.RunTransaction(ctx, txns); err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 생성 실패. Reposioty에서 문제가 발생했습니다.", err)
	}

	utils.ClearCacheByPattern(s.redisCache, ctx, "thread:list:*")

	return thread, nil
}

func (s *ThreadService) ListThread(ctx context.Context, pageNumber, pageSize int) ([]dto.ThreadResponse, *exception.ErrResponseCtx) {
	threadList, err := s.listThreadFromCache(ctx, pageNumber, pageSize)
	if err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. 캐시하는 과정에서 문제가 발생했습니다.", err)
	}
	if threadList != nil {
		return threadList, nil
	}

	listThreadFromRepo, err := s.threadRepo.ListThread(ctx, pageNumber, pageSize)
	if err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. Repository에서 문제가 발생했습니다.", err)
	}

	var listThread []dto.ThreadResponse
	for _, thread := range listThreadFromRepo {
		threadDTO := dto.ThreadResponse{
			ID:        thread.ID,
			UserID:    thread.UserID,
			Handle:    "",
			Title:     thread.Title,
			Content:   thread.Content,
			Views:     thread.Views,
			Likes:     thread.Likes,
			Dislikes:  thread.Dislikes,
			CreatedAt: thread.CreatedAt,
			UpdatedAt: thread.UpdatedAt,
		}
		user, err := s.threadRepo.GetUserByID(ctx, threadDTO.UserID)
		if err != nil {
			return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. Repository에서 문제가 발생했습니다.", err)
		}
		threadDTO.Handle = user.Handle
		listThread = append(listThread, threadDTO)
	}

	if err := s.setListThreadToCache(ctx, pageNumber, pageSize, listThread, 5*time.Minute); err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. 캐시에 저장하지 못했습니다.", err)
	}

	return listThread, nil
}

func (s *ThreadService) ListThreadByHandle(ctx context.Context, handle string) ([]model.ThreadModel, *exception.ErrResponseCtx) {
	threadList, err := s.listThreadByHandleFromCache(ctx, handle)
	if err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. 캐시하는 과정에서 문제가 발생했습니다.", err)
	}
	if len(threadList) > 0 {
		return threadList, nil
	}

	threadList, err = s.threadRepo.ListThreadByHandle(ctx, handle)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			return nil, exception.GenerateErrorCtx(fiber.StatusNotFound, "❌ 쓰레드 조회 실패. 존재하지 않는 사용자입니다.", err)
		default:
			return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. Repository에서 문제가 발생했습니다.", err)
		}
	}

	if err := s.setListThreadByHandleToCache(ctx, handle, threadList, 5*time.Minute); err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. 캐시에 저장하지 못했습니다.", err)
	}

	return threadList, nil
}

func (s *ThreadService) GetThreadByID(ctx context.Context, threadID int) (*model.ThreadModel, *exception.ErrResponseCtx) {
	if err := s.IncrementViews(ctx, threadID); err != nil {
		return nil, err
	}

	thread, errs := s.getThreadFromCache(ctx, threadID)
	if len(errs) > 0 {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. 캐시하는 과정에서 문제가 발생했습니다.", errs)
	}
	if thread != nil {
		return thread, nil
	}

	thread, err := s.threadRepo.GetThreadByID(ctx, threadID)
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

func (s *ThreadService) RemoveThreadByID(ctx context.Context, userID string, threadID int) *exception.ErrResponseCtx {
	thread, err := s.threadRepo.GetThreadByID(ctx, threadID)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			return exception.GenerateErrorCtx(fiber.StatusNotFound, "❌ 쓰레드 삭제 실패. 존재하지 않는 쓰레드입니다.", err)
		default:
			return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 삭제 실패. Repository에서 문제가 발생했습니다.", err)
		}
	}

	if thread.UserID != userID {
		return exception.GenerateErrorCtx(fiber.StatusForbidden, "❌ 쓰레드 삭제 실패. 해당 쓰레드를 삭제할 권한이 없습니다.", err)
	}

	ok, err := s.threadRepo.RemoveThreadByID(ctx, userID, threadID)
	if err != nil {
		return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 삭제 실패. Repository에서 문제가 발생했습니다.", err)
	} else if !ok {
		return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 삭제 실패. 쓰레드를 삭제할 수 없습니다.", err)
	}

	utils.ClearCacheByPattern(s.redisCache, ctx, "thread:list:*")

	return nil
}

func (s *ThreadService) IncrementViews(ctx context.Context, threadID int) *exception.ErrResponseCtx {
	return s.incrementInteraction(ctx, threadID, "views")
}

func (s *ThreadService) IncrementLikes(ctx context.Context, threadID int) *exception.ErrResponseCtx {
	return s.incrementInteraction(ctx, threadID, "likes")
}

func (s *ThreadService) IncrementDislikes(ctx context.Context, threadID int) *exception.ErrResponseCtx {
	return s.incrementInteraction(ctx, threadID, "dislikes")
}

/*
이 함수 로직을 조금 까다롭게 짜서 주석을 조금 남겨봄..
쓰레드 인터렉션은 이렇게 정의됨. -> 조회수, 좋아요, 싫어요, 공유 등등....
물론 인터렉션 별로 함수를 만들어서 관리해도 좋지만, 재사용성을 위해서 이렇게 제작함.
또한, DB에 직접적인 접근은 지양하기 위해서 Redis를 Cache의 용도로 사용함.
Redis 캐시 전략은 Key 값을 각각 인터렉션 별로 구분해서 저장함. (ex. 조회수 -> thread:{id}:views, 좋아요 -> thread:{id}:likes, ...)
해당 값이 만약 일정 수준에 도달하면 DB에 바로 반영하는 것이 아니라, 트랜잭션으로 만들어 저장함.
이후 트랜잭션의 양이 일정 수준에 도달하면 DB에 반영함.
이렇게 구현한 이유는 재사용성을 위해 인터렉션들을 통합해서 하나의 함수로 만들었으므로, 트랜잭션의 상태도 동시에 관리하기 위함임.
*/
func (s *ThreadService) incrementInteraction(ctx context.Context, threadID int, interactionField string) *exception.ErrResponseCtx {
	cachePattern := fmt.Sprintf("thread:%d:%s", threadID, interactionField)
	threadItrAmount, err := s.redisCache.Incr(ctx, cachePattern).Result()
	if err != nil {
		return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 인터렉션 증가 실패. 캐시하는 과정에서 문제가 발생했습니다.", err)
	}

	// 일정 수준 인터렉션 값이 올라가면 DB 트랜잭션에 쿼리 하나 저장
	if threadItrAmount%config.Envs.RedisInteractionAmount == 0 {
		txn := s.incrementInteractionMapper(ctx, interactionField, threadID, int(threadItrAmount))
		s.txnsItr = append(s.txnsItr, txn)
		return nil
	}

	// 일정 트랜잭션 수에 도달하면 한번에 DB Push
	if len(s.txnsItr) >= int(config.Envs.RedisInteractionCount) {
		if err := s.threadRepo.RunTransaction(ctx, s.txnsItr); err != nil {
			return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 인터렉션 증가 실패. Reposioty에서 문제가 발생했습니다.", err)
		}

		// DB Push가 마무리 되었다면 트랜잭션 초기화
		s.txnsItr = make([]model.PrismaTransaction, 0)
	}

	return nil
}

func (s *ThreadService) incrementInteractionMapper(ctx context.Context, interaction string, threadID int, amount int) model.PrismaTransaction {
	userInteractions := []struct {
		field         string
		incrementFunc func(context.Context, int, int) model.ThreadUniqueTxResult
	}{
		{"views", s.threadRepo.IncrementViews},
		{"likes", s.threadRepo.IncrementLikes},
		{"dislikes", s.threadRepo.IncrementDislikes},
	}

	for _, itraction := range userInteractions {
		if itraction.field == interaction {
			return itraction.incrementFunc(ctx, threadID, amount)
		}
	}

	return nil
}

func (s *ThreadService) listThreadFromCache(ctx context.Context, pageNumber, pageSize int) ([]dto.ThreadResponse, error) {
	var threadList []dto.ThreadResponse
	err := utils.GetCache(s.redisCache, ctx, fmt.Sprintf("thread:list:page:%d:size:%d", pageNumber, pageSize), &threadList)
	return threadList, err
}

func (s *ThreadService) setListThreadToCache(ctx context.Context, pageNumber, pageSize int, threadList []dto.ThreadResponse, ttl time.Duration) error {
	return utils.SetCache(s.redisCache, ctx, fmt.Sprintf("thread:list:page:%d:size:%d", pageNumber, pageSize), threadList, ttl)
}

func (s *ThreadService) listThreadByHandleFromCache(ctx context.Context, handle string) ([]model.ThreadModel, error) {
	var threadList []model.ThreadModel
	err := utils.GetCache(s.redisCache, ctx, fmt.Sprintf("thread:list:handle:%s", handle), &threadList)
	return threadList, err
}

func (s *ThreadService) setListThreadByHandleToCache(ctx context.Context, handle string, threadList []model.ThreadModel, ttl time.Duration) error {
	return utils.SetCache(s.redisCache, ctx, fmt.Sprintf("thread:list:handle:%s", handle), threadList, ttl)
}

func (s *ThreadService) getThreadFromCache(ctx context.Context, threadID int) (*model.ThreadModel, []error) {
	var (
		thread    *model.ThreadModel
		errs      []error
		itrFields = []struct {
			fieldName string
			getter    func(*redis.Client, context.Context, string, interface{}) error
		}{
			{"views", utils.GetCache},
			{"likes", utils.GetCache},
			{"dislikes", utils.GetCache},
		}
		itrAmount = map[string]int{
			"views":    0,
			"likes":    0,
			"dislikes": 0,
		}
	)

	// Cache 값이 없다면,
	if err := utils.GetCache(s.redisCache, ctx, fmt.Sprintf("thread:%d", threadID), &thread); err == nil && thread == nil {
		return nil, nil
	}

	for _, field := range itrFields {
		if amount, exists := itrAmount[field.fieldName]; exists {
			if err := field.getter(s.redisCache, ctx, fmt.Sprintf("thread:%d:%s", threadID, field.fieldName), &amount); err == nil {
				updatedThread, err := s.threadConversion(ctx, thread, field.fieldName, amount)
				if err != nil {
					errs = append(errs, err)
				}
				thread = updatedThread
			} else {
				errs = append(errs, err)
			}
		}
	}

	return thread, errs
}

func (s *ThreadService) threadConversion(ctx context.Context, thread *model.ThreadModel, itrField string, amount int) (*model.ThreadModel, error) {
	fieldToSetter := map[string]func(context.Context, *model.ThreadModel, string, int) (*model.ThreadModel, error){
		"views":    s.setThreadField,
		"likes":    s.setThreadField,
		"dislikes": s.setThreadField,
	}

	if setter, exists := fieldToSetter[itrField]; exists {
		return setter(ctx, thread, itrField, amount)
	}

	return thread, exception.ErrInvalidParameter
}

func (s *ThreadService) setThreadField(ctx context.Context, thread *model.ThreadModel, field string, value int) (*model.ThreadModel, error) {
	if thread == nil {
		return nil, exception.ErrMissingParams
	}

	switch field {
	case "views":
		if thread.Views < value {
			thread.Views = value
		}
	case "likes":
		if thread.Likes < value {
			thread.Likes = value
		}
	case "dislikes":
		if thread.Dislikes < value {
			thread.Dislikes = value
		}
	default:
		return nil, exception.ErrInvalidParameter
	}

	return thread, nil
}

func (s *ThreadService) setThreadToCache(ctx context.Context, thread *model.ThreadModel, ttl time.Duration) error {
	return utils.SetCache(s.redisCache, ctx, fmt.Sprintf("thread:%d", thread.ID), thread, ttl)
}

package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kitae0522/gommunity/internal/dto"
	"github.com/kitae0522/gommunity/internal/model"
	"github.com/kitae0522/gommunity/internal/repository"
	"github.com/kitae0522/gommunity/pkg/exception"
)

type ThreadService struct {
	threadRepo *repository.ThreadRepository
}

func NewThreadService(repo *repository.ThreadRepository) *ThreadService {
	return &ThreadService{threadRepo: repo}
}

func (s *ThreadService) CreateThread(req *dto.CreateThreadRequest) (*model.ThreadModel, *exception.ErrResponseCtx) {
	thread, err := s.threadRepo.CreateThread(req)
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
		txns = append(txns, s.threadRepo.LinkParentThread(thread.ID, *req.ParentThread))
	} else if req.NextThread != nil {
		txns = append(txns, s.threadRepo.LinkNextThread(thread.ID, *req.NextThread))
	} else if req.PrevThread != nil {
		txns = append(txns, s.threadRepo.LinkPrevThread(thread.ID, *req.NextThread))
	}

	if err := s.threadRepo.LinkRelation(txns); err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 생성 실패. Reposioty에서 문제가 발생했습니다.", err)
	}

	return thread, nil
}

func (s *ThreadService) ListThread() ([]model.ThreadModel, *exception.ErrResponseCtx) {
	threadList, err := s.threadRepo.ListThread()
	if err != nil {
		return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. Repository에서 문제가 발생했습니다.", err)
	}
	return threadList, nil
}

func (s *ThreadService) ListThreadByHandle(handle string) ([]model.ThreadModel, *exception.ErrResponseCtx) {
	threadList, err := s.threadRepo.ListThreadByHandle(handle)
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

func (s *ThreadService) GetThreadByID(threadID int) (*model.ThreadModel, *exception.ErrResponseCtx) {
	thread, err := s.threadRepo.GetThreadByID(threadID)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			return nil, exception.GenerateErrorCtx(fiber.StatusNotFound, "❌ 쓰레드 조회 실패. 존재하지 않는 쓰레드입니다.", err)
		default:
			return nil, exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 쓰레드 조회 실패. Repository에서 문제가 발생했습니다.", err)
		}
	}
	return thread, nil
}

func (s *ThreadService) CommentsByID(threadID int) ([]model.ThreadModel, *exception.ErrResponseCtx) {
	comments, err := s.threadRepo.CommentsByID(threadID)
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

package repository

import (
	"context"

	"github.com/kitae0522/gommunity/internal/dto"
	"github.com/kitae0522/gommunity/internal/model"
)

type ThreadRepository struct {
	client *model.PrismaClient
}

func NewThreadRepository(prismaClient *model.PrismaClient) *ThreadRepository {
	return &ThreadRepository{client: prismaClient}
}

func (r *ThreadRepository) CreateThread(req *dto.CreateThreadRequest) (*model.ThreadModel, error) {
	thread, err := r.client.Thread.CreateOne(
		model.Thread.Content.Set(req.Content),
		model.Thread.User.Link(model.Users.ID.Equals(req.UserID)),
		model.Thread.Title.SetIfPresent(req.Title),
		model.Thread.ImgURL.SetIfPresent(req.ImgUrl),
	).Exec(context.Background())

	return thread, err
}

func (r *ThreadRepository) ListThread() ([]model.ThreadModel, error) {
	listThread, err := r.client.Thread.FindMany(
		model.Thread.ParentThread.IsNull(),
	).Exec(context.Background())
	return listThread, err
}

func (r *ThreadRepository) ListThreadByHandle(handle string) ([]model.ThreadModel, error) {
	user, err := r.getUserByHandle(handle)
	if err != nil {
		return nil, err
	}

	listThread, err := r.client.Thread.FindMany(
		model.Thread.UserID.Equals(user.ID),
	).Select(
		model.Thread.ID.Field(),
		model.Thread.Title.Field(),
		model.Thread.ImgURL.Field(),
		model.Thread.Content.Field(),
		model.Thread.ParentThread.Field(),
		model.Thread.NextThread.Field(),
		model.Thread.PrevThread.Field(),
		model.Thread.Views.Field(),
		model.Thread.Likes.Field(),
		model.Thread.Dislikes.Field(),
		model.Thread.CreatedAt.Field(),
		model.Thread.UpdatedAt.Field(),
	).Exec(context.Background())

	return listThread, err
}

func (r *ThreadRepository) GetThreadByID(threadID int) (*model.ThreadModel, error) {
	thread, err := r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Exec(context.Background())

	return thread, err
}

func (r *ThreadRepository) CommentsByID(threadID int) ([]model.ThreadModel, error) {
	commentThreads, err := r.client.Thread.FindMany(
		model.Thread.ParentThread.Equals(threadID),
	).Exec(context.Background())

	return commentThreads, err
}

func (r *ThreadRepository) LinkRelation(txns []model.PrismaTransaction) error {
	if err := r.client.Prisma.Transaction(txns...).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}

func (r *ThreadRepository) LinkParentThread(threadID, parentID int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Parent.Link(
			model.Thread.ID.Equals(parentID),
		),
	).Tx()
}

func (r *ThreadRepository) LinkNextThread(threadID, nextID int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Next.Link(
			model.Thread.ID.Equals(nextID),
		),
	).Tx()
}

func (r *ThreadRepository) LinkPrevThread(threadID, prevID int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Prev.Link(
			model.Thread.ID.Equals(prevID),
		),
	).Tx()
}

func (r *ThreadRepository) getUserByHandle(handle string) (*model.UsersModel, error) {
	user, err := r.client.Users.FindUnique(
		model.Users.Handle.Equals(handle),
	).Exec(context.Background())
	return user, err
}

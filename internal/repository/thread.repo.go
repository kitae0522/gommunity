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

func (r *ThreadRepository) CreateThread(ctx context.Context, req *dto.CreateThreadRequest) (*model.ThreadModel, error) {
	thread, err := r.client.Thread.CreateOne(
		model.Thread.Title.Set(req.Title),
		model.Thread.Content.Set(req.Content),
		model.Thread.User.Link(model.Users.ID.Equals(req.UserID)),
		model.Thread.ImgURL.SetIfPresent(req.ImgUrl),
	).Exec(ctx)

	return thread, err
}

func (r *ThreadRepository) ListThread(ctx context.Context, pageNumber int, pageSize int) ([]model.ThreadModel, error) {
	if pageNumber <= 0 {
		pageNumber = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNumber - 1) * pageSize
	listThread, err := r.client.Thread.FindMany(
		model.Thread.ParentThread.IsNull(),
	).Take(pageSize).Skip(offset).Exec(ctx)
	return listThread, err
}

func (r *ThreadRepository) ListThreadByHandle(ctx context.Context, handle string) ([]model.ThreadModel, error) {
	user, err := r.getUserByHandle(ctx, handle)
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
	).Exec(ctx)

	return listThread, err
}

func (r *ThreadRepository) GetThreadByID(ctx context.Context, threadID int) (*model.ThreadModel, error) {
	thread, err := r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Exec(ctx)

	return thread, err
}

func (r *ThreadRepository) CommentsByID(ctx context.Context, threadID int) ([]model.ThreadModel, error) {
	commentThreads, err := r.client.Thread.FindMany(
		model.Thread.ParentThread.Equals(threadID),
	).Exec(ctx)

	return commentThreads, err
}

func (r *ThreadRepository) RemoveThreadByID(ctx context.Context, userID string, threadID int) (bool, error) {
	_, err := r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Delete().Exec(ctx)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ThreadRepository) IncrementViews(ctx context.Context, threadID int, amount int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Views.Increment(amount),
	).Tx()
}

func (r *ThreadRepository) IncrementLikes(ctx context.Context, threadID int, amount int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Likes.Increment(amount),
	).Tx()
}

func (r *ThreadRepository) IncrementDislikes(ctx context.Context, threadID int, amount int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Dislikes.Increment(amount),
	).Tx()
}

func (r *ThreadRepository) DecrementLikes(ctx context.Context, threadID int, amount int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Likes.Decrement(amount),
	).Tx()
}

func (r *ThreadRepository) DecrementDislikes(ctx context.Context, threadID int, amount int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Dislikes.Decrement(amount),
	).Tx()
}

func (r *ThreadRepository) RunTransaction(ctx context.Context, txns []model.PrismaTransaction) error {
	if err := r.client.Prisma.Transaction(txns...).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *ThreadRepository) LinkParentThread(ctx context.Context, threadID, parentID int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Parent.Link(
			model.Thread.ID.Equals(parentID),
		),
	).Tx()
}

func (r *ThreadRepository) LinkNextThread(ctx context.Context, threadID, nextID int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Next.Link(
			model.Thread.ID.Equals(nextID),
		),
	).Tx()
}

func (r *ThreadRepository) LinkPrevThread(ctx context.Context, threadID, prevID int) model.ThreadUniqueTxResult {
	return r.client.Thread.FindUnique(
		model.Thread.ID.Equals(threadID),
	).Update(
		model.Thread.Prev.Link(
			model.Thread.ID.Equals(prevID),
		),
	).Tx()
}

func (r *ThreadRepository) GetUserByID(ctx context.Context, id string) (*model.UsersModel, error) {
	return r.client.Users.FindUnique(model.Users.ID.Equals(id)).Exec(ctx)
}

func (r *ThreadRepository) getUserByHandle(ctx context.Context, handle string) (*model.UsersModel, error) {
	user, err := r.client.Users.FindUnique(
		model.Users.Handle.Equals(handle),
	).Exec(ctx)
	return user, err
}

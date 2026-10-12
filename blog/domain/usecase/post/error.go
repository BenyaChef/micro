package postusecase

import (
	"context"
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const (
	ErrCodeOperationFailed = errors.ErrorCode("cce2da53-010")
	ErrCodePostNotFound    = errors.ErrorCode("cce2da53-020")
	ErrCodeNotAuthor       = errors.ErrorCode("cce2da53-030")
)

var (
	ErrPostRepositoryIsRequired = errors.NewError("SYS", "PostUseCase: PostRepository is required")
	ErrLogPublisherIsRequired   = errors.NewError("SYS", "PostUseCase: LogPublisher is required")
)

var (
	ErrInputIsRequired = errors.NewErrorWithLevel(
		"cce2da53-001", "Post: request input is required", errors.LevelInfo,
	)
	ErrPostIDIsRequired = errors.NewErrorWithLevel(
		"cce2da53-002", "Post: postID is required", errors.LevelInfo,
	)
	ErrAuthorIDIsRequired = errors.NewErrorWithLevel(
		"cce2da53-003", "Post: authorID is required", errors.LevelInfo,
	)
)

func ErrPostNotFound(postID string) error {
	msg := fmt.Sprintf("Post not found: postID=%s", postID)

	return errors.NewErrorWithLevel(ErrCodePostNotFound, msg, errors.LevelInfo)
}

func ErrNotAuthor(postID, authorID string) error {
	msg := fmt.Sprintf("Post belongs to another author: postID=%s, authorID=%s", postID, authorID)

	return errors.NewErrorWithLevel(ErrCodeNotAuthor, msg, errors.LevelInfo)
}

func ErrOperationFailed(cause error) error {
	return errors.NewErrorWithCause(ErrCodeOperationFailed, "Post operation failed", cause)
}

type errorProcessor struct {
	logPublisher loggerinterface.LogPublisher
}

func newErrorProcessor(logPublisher loggerinterface.LogPublisher) *errorProcessor {
	return &errorProcessor{logPublisher: logPublisher}
}

func (p *errorProcessor) LogAndReturn(ctx context.Context, err error) error {
	p.logPublisher.LogError(ctx, err)

	return err
}

// Нижний слой уже залогировал причину, поэтому здесь только оборачиваем:
// иначе одна авария даёт две строки в логе.
func (p *errorProcessor) WrapOperationError(cause error) error {
	return ErrOperationFailed(cause)
}

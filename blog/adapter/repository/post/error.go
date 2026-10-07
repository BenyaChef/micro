package postrepository

import (
	"context"
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const (
	ErrCodeInsertFailed = errors.ErrorCode("6dadc350-010")
	ErrCodeSelectFailed = errors.ErrorCode("6dadc350-011")
	ErrCodeUpdateFailed = errors.ErrorCode("6dadc350-012")
	ErrCodeDeleteFailed = errors.ErrorCode("6dadc350-013")
	ErrCodePostVanished = errors.ErrorCode("6dadc350-020")
)

var (
	ErrExecutorIsRequired     = errors.NewError("SYS", "PostRepository: Executor is required")
	ErrLogPublisherIsRequired = errors.NewError("SYS", "PostRepository: LogPublisher is required")
)

var ErrPostIsRequired = errors.NewErrorWithLevel(
	"6dadc350-001", "PostRepository: post is required", errors.LevelInfo,
)

func ErrInsertFailed(postID string, cause error) error {
	msg := fmt.Sprintf("Unable to insert post: postID=%s", postID)

	return errors.NewErrorWithCause(ErrCodeInsertFailed, msg, cause)
}

func ErrSelectFailed(cause error) error {
	return errors.NewErrorWithCause(ErrCodeSelectFailed, "Unable to select post", cause)
}

func ErrUpdateFailed(postID string, cause error) error {
	msg := fmt.Sprintf("Unable to update post: postID=%s", postID)

	return errors.NewErrorWithCause(ErrCodeUpdateFailed, msg, cause)
}

func ErrDeleteFailed(postID string, cause error) error {
	msg := fmt.Sprintf("Unable to delete post: postID=%s", postID)

	return errors.NewErrorWithCause(ErrCodeDeleteFailed, msg, cause)
}

func ErrPostVanished(postID string) error {
	msg := fmt.Sprintf("Post disappeared between read and write: postID=%s", postID)

	return errors.NewErrorWithLevel(ErrCodePostVanished, msg, errors.LevelWarn)
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

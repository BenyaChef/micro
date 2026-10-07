package postrestcontroller

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	requestrestmodel "github.com/BenyaChef/micro/blog/adapter/controller/rest/post/request"
	responserestmodel "github.com/BenyaChef/micro/blog/adapter/controller/rest/post/response"
	boundarymodel "github.com/BenyaChef/micro/blog/boundary/model"
	usecaseinterface "github.com/BenyaChef/micro/blog/boundary/usecase"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	restservercontroller "github.com/BenyaChef/micro/infrastructure/restserver/controller"
	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
)

const (
	paramPostID = "postID"
	queryLimit  = "limit"
	queryOffset = "offset"
)

type PostController struct {
	*restservercontroller.BaseController

	postUseCase usecaseinterface.PostUseCase
}

func (c *PostController) List(writer http.ResponseWriter, request *http.Request) {
	limit, err := intQueryParam(request, queryLimit)
	if err != nil {
		c.ErrorResponseWithLog(writer, request, err)

		return
	}

	offset, err := intQueryParam(request, queryOffset)
	if err != nil {
		c.ErrorResponseWithLog(writer, request, err)

		return
	}

	posts, err := c.postUseCase.List(request.Context(), &boundarymodel.ListPostsInput{Limit: limit, Offset: offset})
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	c.JSONResponse(writer, request, responserestmodel.FromPosts(posts), http.StatusOK)
}

func (c *PostController) GetByID(writer http.ResponseWriter, request *http.Request) {
	postID, err := postIDFromPath(request)
	if err != nil {
		c.ErrorResponseWithLog(writer, request, err)

		return
	}

	post, err := c.postUseCase.GetByID(request.Context(), postID)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	c.JSONResponse(writer, request, responserestmodel.FromPost(post), http.StatusOK)
}

func (c *PostController) Create(writer http.ResponseWriter, request *http.Request) {
	requestModel := &requestrestmodel.CreatePostRequest{}
	if err := c.FillAndValidateReqModel(request, requestModel); err != nil {
		c.ErrorResponseWithLog(writer, request, err)

		return
	}

	authorID, err := callerID(request)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	input := &boundarymodel.CreatePostInput{
		AuthorID: authorID,
		Title:    requestModel.Title,
		Body:     requestModel.Body,
	}

	post, err := c.postUseCase.Create(request.Context(), input)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	c.JSONResponse(writer, request, responserestmodel.FromPost(post), http.StatusCreated)
}

func (c *PostController) Update(writer http.ResponseWriter, request *http.Request) {
	requestModel := &requestrestmodel.UpdatePostRequest{}
	if err := c.FillAndValidateReqModel(request, requestModel); err != nil {
		c.ErrorResponseWithLog(writer, request, err)

		return
	}

	postID, err := postIDFromPath(request)
	if err != nil {
		c.ErrorResponseWithLog(writer, request, err)

		return
	}

	authorID, err := callerID(request)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	input := &boundarymodel.UpdatePostInput{
		PostID:   postID,
		AuthorID: authorID,
		Title:    requestModel.Title,
		Body:     requestModel.Body,
	}

	post, err := c.postUseCase.Update(request.Context(), input)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	c.JSONResponse(writer, request, responserestmodel.FromPost(post), http.StatusOK)
}

func (c *PostController) Delete(writer http.ResponseWriter, request *http.Request) {
	postID, err := postIDFromPath(request)
	if err != nil {
		c.ErrorResponseWithLog(writer, request, err)

		return
	}

	authorID, err := callerID(request)
	if err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	if err := c.postUseCase.Delete(request.Context(), postID, authorID); err != nil {
		c.ErrorResponse(writer, request, err)

		return
	}

	c.JSONResponse(writer, request, nil, http.StatusNoContent)
}

func postIDFromPath(request *http.Request) (postentity.PostID, error) {
	return postentity.PostIDFrom(chi.URLParam(request, paramPostID))
}

func callerID(request *http.Request) (commonuserentity.UserID, error) {
	rawCallerID := restservermodel.CallerIDFrom(request.Context())
	if rawCallerID == "" {
		return commonuserentity.UserID{}, ErrCallerIsRequired
	}

	return commonuserentity.UserIDFrom(rawCallerID)
}

func intQueryParam(request *http.Request, name string) (int, error) {
	rawValue := request.URL.Query().Get(name)
	if rawValue == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, ErrInvalidQueryParam(name, rawValue)
	}

	return value, nil
}

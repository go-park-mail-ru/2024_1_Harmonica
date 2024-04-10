package handler

import (
	"context"
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

// Create board.
//
//	@Summary		Create board
//	@Description	Create board by description
//	@Tags			Boards
//	@Produce		json
//	@Accept			json
//	@Param			Cookie	header		string			true	"session-token"	default(session-token=)
//	@Param			board	body		entity.Board	false	"Board information"
//	@Success		200		{object}	entity.FullBoard
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 3, 4, 5, 21."
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/boards [post]
func (h *APIHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := r.Context().Value("request_id").(string)

	board := entity.Board{}
	board.VisibilityType = "public"
	err := UnmarshalRequest(r, &board)

	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	if !ValidateBoard(board) {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrInvalidInputFormat))
		return
	}
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrTypeConversion))
		return
	}
	fullBoard, errInfo := h.service.CreateBoard(ctx, board, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, fullBoard)
}

// Get board.
//
//	@Summary		Get board
//	@Description	Get board by id
//	@Tags			Boards
//	@Produce		json
//	@Accept			json
//	@Param			Cookie		header		string	true	"session-token"	default(session-token=)
//	@Param			board_id	path		int		true	"Board ID"
//	@Success		200			{object}	entity.FullBoard
//	@Failure		400			{object}	errs.ErrorResponse	"Possible code responses: 4, 12, 21."
//	@Failure		401			{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403			{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500			{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/boards/{board_id}/ [get]
func (h *APIHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := r.Context().Value("request_id").(string)

	boardId, err := ReadInt64Slug(r, "board_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	userId, errInfo := CheckAuth(ctx)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
	}

	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	board, errInfo := h.service.GetBoardById(ctx, entity.BoardID(boardId), userId, limit, offset)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, board)
}

// Update board.
//
//	@Summary		Update board
//	@Description	Update board by board information
//	@Tags			Boards
//	@Produce		json
//	@Accept			multipart/form-data
//	@Param			Cookie		header		string	true	"session-token"	default(session-token=)
//	@Param			board_id	path		int		false	"Board ID"
//	@Param			image		formData	file	false	"Cover image"
//	@Param			board		formData	string	false	"Board information in json"
//	@Success		200			{object}	entity.FullBoard
//	@Failure		400			{object}	errs.ErrorResponse	"Possible code responses: 4, 12, 17, 18, 21."
//	@Failure		401			{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403			{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500			{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/boards/{board_id}/ [post]
func (h *APIHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := r.Context().Value("request_id").(string)

	boardId, err := ReadInt64Slug(r, "board_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	var newBoard entity.Board
	newBoard.VisibilityType = "public"

	boardParams := r.FormValue("board")
	err = json.Unmarshal([]byte(boardParams), &newBoard)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	if !ValidateBoard(newBoard) {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrInvalidInputFormat))
		return
	}
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrTypeConversion))
		return
	}
	newBoard.CoverURL = ""
	image, imageHeader, err := r.FormFile("image")
	if err == nil {
		name, errUploading := h.service.UploadImage(ctx, image, imageHeader)
		if errUploading != nil {
			WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrInvalidImg))
			return
		}
		newBoard.CoverURL = FormImgURL(name)
	}
	newBoard.BoardID = entity.BoardID(boardId)
	board, errInfo := h.service.UpdateBoard(ctx, newBoard, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, board)
}

// Add pin to board.
//
//	@Summary		Add pin to board
//	@Description	Add pin to board by pin id and board id.
//	@Tags			Boards
//	@Produce		json
//	@Accept			json
//	@Param			Cookie		header		string	true	"session-token"	default(session-token=)
//	@Param			board_id	path		int		true	"Board ID"
//	@Param			pin_id		path		int		true	"Board ID"
//	@Success		200			{object}	interface{}
//	@Failure		400			{object}	errs.ErrorResponse	"Possible code responses: 4, 12, 21."
//	@Failure		401			{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403			{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500			{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/boards/{board_id}/pins/{pin_id}/ [post]
func (h *APIHandler) AddPinToBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	boardId, pinId, userId, errInfo := GetInfoFromSlugAndContext(r)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	errInfo = h.service.AddPinToBoard(ctx, boardId, pinId, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

// Delete pin from board.
//
//	@Summary		Delete pin from board
//	@Description	Delete pin from board by pin id and board id.
//	@Tags			Boards
//	@Produce		json
//	@Accept			json
//	@Param			Cookie		header		string	true	"session-token"	default(session-token=)
//	@Param			board_id	path		int		true	"Board ID"
//	@Param			pin_id		path		int		true	"Board ID"
//	@Success		200			{object}	interface{}
//	@Failure		400			{object}	errs.ErrorResponse	"Possible code responses: 4, 12, 21."
//	@Failure		401			{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403			{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500			{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/boards/{board_id}/pins/{pin_id}/ [post]
func (h *APIHandler) DeletePinFromBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	boardId, pinId, userId, errInfo := GetInfoFromSlugAndContext(r)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	errInfo = h.service.DeletePinFromBoard(ctx, boardId, pinId, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

// Delete board.
//
//	@Summary		Delete board
//	@Description	Delete board by board id.
//	@Tags			Boards
//	@Produce		json
//	@Param			Cookie		header		string	true	"session-token"	default(session-token=)
//	@Param			board_id	path		int		true	"Board ID"
//	@Success		200			{object}	interface{}
//	@Failure		400			{object}	errs.ErrorResponse	"Possible code responses: 4, 12, 21."
//	@Failure		401			{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403			{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500			{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/boards/{board_id}/ [delete]
func (h *APIHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := r.Context().Value("request_id").(string)

	boardId, err := ReadInt64Slug(r, "board_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrTypeConversion))
		return
	}
	errInfo := h.service.DeleteBoard(ctx, entity.BoardID(boardId), userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

// Get boards created by user.
//
//	@Summary		Get boards created by user
//	@Description	Get boards created by user by user nickname.
//	@Tags			Boards
//	@Produce		json
//	@Param			Cookie		header		string	true	"session-token"	default(session-token=)
//	@Param			nickname	path		string	true	"user nickname"
//	@Success		200			{object}	entity.UserBoards
//	@Failure		400			{object}	errs.ErrorResponse	"Possible code responses: 4, 12, 21."
//	@Failure		401			{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403			{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500			{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/boards/created/{nickname}/ [get]
func (h *APIHandler) UserBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := r.Context().Value("request_id").(string)

	authorNickname := r.PathValue("nickname")
	if !ValidateNickname(authorNickname) {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrInvalidSlug))
		return
	}
	userId, errInfo := CheckAuth(ctx)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
	}

	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	boards, errInfo := h.service.GetUserBoards(ctx, authorNickname, userId, limit, offset)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, boards)
}

func GetInfoFromSlugAndContext(r *http.Request) (entity.BoardID, entity.PinID, entity.UserID, errs.ErrorInfo) {
	ctx := r.Context()
	boardStringId, err := ReadInt64Slug(r, "board_id")
	if err != nil {
		return 0, 0, 0, MakeErrorInfo(err, errs.ErrInvalidSlug)
	}
	pinStringId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		return 0, 0, 0, MakeErrorInfo(err, errs.ErrInvalidSlug)
	}
	boardId := entity.BoardID(boardStringId)
	pinId := entity.PinID(pinStringId)
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		return 0, 0, 0, MakeErrorInfo(nil, errs.ErrTypeConversion)
	}
	return boardId, pinId, userId, errs.ErrorInfo{}
}

func CheckAuth(ctx context.Context) (entity.UserID, errs.ErrorInfo) {
	userIdString := ctx.Value("user_id")
	var (
		userId entity.UserID
		ok     bool
	)
	if userIdString != nil {
		userId, ok = userIdString.(entity.UserID)
		if !ok {
			return 0, errs.ErrorInfo{LocalErr: errs.ErrTypeConversion}
		}
	}
	return userId, emptyErrorInfo
}

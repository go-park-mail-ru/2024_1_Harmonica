package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

// CreateBoard Create board.
//
//	@Summary		Create board
//	@Description	Create board by description
//	@Tags			Pins
//	@Produce		json
//	@Accept			multipart/form-data
//	@Param			Cookie	header		string	true	"session-token"	default(session-token=)
//	@Param			board		body  entity.Board	string	false	"Board information in json"
//	@Param			image	formData	file	true	"Board cover"
//	@Success		200		{object}	entity.FullBoard
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: ."
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: ."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/boards [post]
func (h *APIHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	board := entity.Board{}
	err := UnmarshalRequest(r, &board)
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	// нужна нормальная валидация ?
	if !ValidateBoard(board) {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrInvalidInputFormat))
		return
	}
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrTypeConversion))
	}
	res, errInfo := h.service.CreateBoard(ctx, board, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, res)
}

func (h *APIHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	boardId, err := ReadInt64Slug(r, "board_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	var (
		userId entity.UserID
		ok     bool
	)
	if ctx.Value("is_auth") == true {
		userId, ok = ctx.Value("user_id").(entity.UserID)
		if !ok {
			WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrTypeConversion))
		}
	}
	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	board, errInfo := h.service.GetBoardById(ctx, entity.BoardID(boardId), userId, limit, offset)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, board)
}

func (h *APIHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	boardId, err := ReadInt64Slug(r, "board_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	var newBoard entity.Board
	err = UnmarshalRequest(r, &newBoard)
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	if !ValidateBoard(newBoard) {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrInvalidInputFormat))
		return
	}
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrTypeConversion))
	}
	newBoard.BoardID = entity.BoardID(boardId)
	board, errInfo := h.service.UpdateBoard(ctx, newBoard, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, board)
}

func (h *APIHandler) AddPinToBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	boardId, err := ReadInt64Slug(r, "board_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	pinId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrTypeConversion))
	}
	errInfo := h.service.AddPinToBoard(ctx, entity.BoardID(boardId), entity.PinID(pinId), userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

func (h *APIHandler) DeletePinFromBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	boardId, err := ReadInt64Slug(r, "board_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	pinId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrTypeConversion))
	}
	errInfo := h.service.DeletePinFromBoard(ctx, entity.BoardID(boardId), entity.PinID(pinId), userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

func (h *APIHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	boardId, err := ReadInt64Slug(r, "board_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrTypeConversion))
	}
	errInfo := h.service.DeleteBoard(ctx, entity.BoardID(boardId), userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

func (h *APIHandler) UserBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorNickname := r.PathValue("nickname")
	if !ValidateNickname(authorNickname) {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrInvalidSlug))
		return
	}
	var (
		userId entity.UserID
		ok     bool
	)
	if ctx.Value("is_auth") == true {
		userId, ok = ctx.Value("user_id").(entity.UserID)
		if !ok {
			WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrTypeConversion))
		}
	}
	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	boards, errInfo := h.service.GetUserBoards(ctx, authorNickname, userId, limit, offset)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, boards)
}

package handler

import (
	"encoding/json"
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
//	@Accept			json
//	@Param			Cookie	header		string	true	"session-token"	default(session-token=)
//	@Param			board		body  entity.Board	string	false	"Board information in json"
//	@Success		200		{object}	entity.FullBoard
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: ."
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: ."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/boards [post]
func (h *APIHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	board := entity.Board{}
	board.VisibilityType = "public"
	err := UnmarshalRequest(r, &board)

	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	if !ValidateBoard(board) {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrInvalidInputFormat))
		return
	}
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrTypeConversion))
		return
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
	userStringId := ctx.Value("user_id")
	if userStringId == nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrUnauthorized))
		return
	}
	userId, ok := userStringId.(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrTypeConversion))
		return
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
	newBoard.VisibilityType = "public"

	boardParams := r.FormValue("board")
	err = json.Unmarshal([]byte(boardParams), &newBoard)
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
		return
	}
	newBoard.CoverURL = ""
	image, imageHeader, err := r.FormFile("image")
	if err == nil {
		name, errUploading := h.service.UploadImage(ctx, image, imageHeader)
		if errUploading != nil {
			WriteErrorResponse(w, h.logger, errs.ErrorInfo{
				GeneralErr: err,
				LocalErr:   errs.ErrInvalidImg,
			})
			return
		}
		newBoard.CoverURL = FormImgURL(name)
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
		return
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
		return
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
		return
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
	userStringId := ctx.Value("user_id")
	if userStringId == nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrUnauthorized))
		return
	}
	userId, ok := userStringId.(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrTypeConversion))
		return
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

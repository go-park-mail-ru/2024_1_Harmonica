package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

func (h *APIHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	board := entity.Board{}
	err := UnmarshalRequest(r, &board)
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	//TODO сделать нормальную валидацию
	if !ValidateBoard(board) {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrInvalidInputFormat))
		return
	}

	userId := ctx.Value("user_id").(entity.UserID)
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

	userId := entity.UserID(0)
	if ctx.Value("is_auth") == true {
		userIdFromCtx := ctx.Value("user_id")
		userId = userIdFromCtx.(entity.UserID)
	}
	pin, errInfo := h.service.GetBoardById(ctx, entity.BoardID(boardId), userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, pin)
}

func (h *APIHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {

}

func (h *APIHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {}

func (h *APIHandler) UserBoards(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	userNickname := r.PathValue("nickname")
	if !ValidateNickname(userNickname) {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(nil, errs.ErrInvalidSlug))
		return
	}
	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	boards, errInfo := h.service.GetUserBoards(r.Context(), userNickname, limit, offset)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, boards)
}

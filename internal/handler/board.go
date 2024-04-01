package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

func (h *APIHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	board := entity.FullBoard{}
	err := UnmarshalRequest(r, &board)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	board.BoardAuthor.UserId = ctx.Value("user_id").(entity.UserID)
	res, errInfo := h.service.CreateBoard(ctx, board)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, res)
}

func (h *APIHandler) UserBoards(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	userNickname := r.PathValue("nickname")
	if !ValidateNickname(userNickname) {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrInvalidSlug,
		})
		return
	}
	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	boards, errInfo := h.service.GetUserBoards(r.Context(), userNickname, limit, offset)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, boards)
}

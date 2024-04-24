package test_service

//
//import (
//	"context"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//	"harmonica/internal/entity"
//	"harmonica/internal/entity/errs"
//	"harmonica/internal/service"
//	mock_repository "harmonica/mocks/repository"
//	"testing"
//)
//
//func TestCreateBoard(t *testing.T) {
//	type mockArgs struct {
//		Ctx    context.Context
//		Board  entity.Board
//		UserId entity.UserID
//	}
//	type mockReturn struct {
//		BoardCreate    entity.Board
//		ErrBoardCreate error
//		Author         entity.User
//		AuthorErr      error
//	}
//	type funcArgs struct {
//		Ctx    context.Context
//		Board  entity.Board
//		UserId entity.UserID
//	}
//	type funcReturn struct {
//		Board entity.FullBoard
//		Err   errs.ErrorInfo
//	}
//	type test struct {
//		Name                   string
//		MockArgs               mockArgs
//		MockReturn             mockReturn
//		FuncArgs               funcArgs
//		ExpectedFuncReturn     funcReturn
//		ExcpectGetUserByIdCall bool
//	}
//	tests := []test{
//		{
//			Name: "Correct work test 1",
//			MockArgs: mockArgs{
//				Ctx:    context.Background(),
//				Board:  entity.Board{},
//				UserId: entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				BoardCreate:    entity.Board{},
//				ErrBoardCreate: nil,
//				Author:         entity.User{},
//				AuthorErr:      nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:    context.Background(),
//				Board:  entity.Board{},
//				UserId: entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{
//					BoardAuthors: []entity.BoardAuthor{
//						{
//							UserId: entity.UserID(1),
//						},
//					},
//				},
//			},
//			ExcpectGetUserByIdCall: true,
//		},
//		{
//			Name: "Uncorrect work test 1",
//			MockArgs: mockArgs{
//				Ctx:    context.Background(),
//				Board:  entity.Board{},
//				UserId: entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				BoardCreate:    entity.Board{},
//				ErrBoardCreate: nil,
//				Author:         entity.User{},
//				AuthorErr:      errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:    context.Background(),
//				Board:  entity.Board{},
//				UserId: entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{},
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//			ExcpectGetUserByIdCall: true,
//		},
//		{
//			Name: "Uncorrect work test 2",
//			MockArgs: mockArgs{
//				Ctx:    context.Background(),
//				Board:  entity.Board{},
//				UserId: entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				BoardCreate:    entity.Board{},
//				ErrBoardCreate: errs.ErrDBInternal,
//				Author:         entity.User{},
//				AuthorErr:      errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:    context.Background(),
//				Board:  entity.Board{},
//				UserId: entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{},
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//			ExcpectGetUserByIdCall: false,
//		},
//	}
//	ctrl := gomock.NewController(t)
//	repo := mock_repository.NewMockIRepository(ctrl)
//	for _, test := range tests {
//		repo.EXPECT().CreateBoard(test.MockArgs.Ctx, test.MockArgs.Board, test.MockArgs.UserId).Return(
//			test.MockReturn.BoardCreate, test.MockReturn.ErrBoardCreate)
//		if test.ExcpectGetUserByIdCall {
//			repo.EXPECT().GetUserById(test.MockArgs.Ctx, test.MockArgs.UserId).Return(
//				test.MockReturn.Author, test.MockReturn.AuthorErr)
//		}
//		service := service.NewService(repo)
//		board, err := service.CreateBoard(test.FuncArgs.Ctx, test.FuncArgs.Board, test.FuncArgs.UserId)
//		assert.Equal(t, test.ExpectedFuncReturn.Board, board)
//		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
//	}
//}
//
//func TestGetBoardById(t *testing.T) {
//	type mockArgs struct {
//		Ctx     context.Context
//		BoardId entity.BoardID
//		UserId  entity.UserID
//		Limit   int
//		Offset  int
//	}
//	type mockReturn struct {
//		BoardByID                    entity.Board
//		ErrBoardByID                 error
//		CheckBoardAuthorExistence    bool
//		ErrCheckBoardAuthorExistence error
//		GetBoardAuthors              []entity.BoardAuthor
//		ErrGetBoardAuthors           error
//		GetBoardPins                 []entity.BoardPinResponse
//		ErrGetBoardPins              error
//	}
//	type funcArgs struct {
//		Ctx     context.Context
//		BoardId entity.BoardID
//		UserId  entity.UserID
//		Limit   int
//		Offset  int
//	}
//	type funcReturn struct {
//		Board entity.FullBoard
//		Err   errs.ErrorInfo
//	}
//	type test struct {
//		Name               string
//		MockArgs           mockArgs
//		MockReturn         mockReturn
//		FuncArgs           funcArgs
//		ExpectedFuncReturn funcReturn
//	}
//	tests := []test{
//		{
//			Name: "Correct work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  10,
//			},
//			MockReturn: mockReturn{
//				BoardByID: entity.Board{BoardID: entity.BoardID(1)},
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{
//					Board: entity.Board{BoardID: entity.BoardID(1)},
//				},
//			},
//		},
//		{
//			Name: "Correct work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  10,
//			},
//			MockReturn: mockReturn{
//				BoardByID:       entity.Board{BoardID: entity.BoardID(1)},
//				ErrGetBoardPins: errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{},
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//		},
//		{
//			Name: "Uncorrect work test 2",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  10,
//			},
//			MockReturn: mockReturn{
//				BoardByID:          entity.Board{BoardID: entity.BoardID(1)},
//				ErrGetBoardAuthors: errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{},
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//		},
//		{
//			Name: "Uncorrect work test 3",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  10,
//			},
//			MockReturn: mockReturn{
//				BoardByID:                    entity.Board{BoardID: entity.BoardID(1)},
//				ErrGetBoardAuthors:           errs.ErrDBInternal,
//				ErrCheckBoardAuthorExistence: errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{},
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//		},
//		{
//			Name: "Uncorrect work test 4",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  10,
//			},
//			MockReturn: mockReturn{
//				BoardByID:                    entity.Board{BoardID: entity.BoardID(1)},
//				ErrBoardByID:                 errs.ErrDBInternal,
//				ErrGetBoardAuthors:           errs.ErrDBInternal,
//				ErrCheckBoardAuthorExistence: errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{},
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//		},
//	}
//	ctrl := gomock.NewController(t)
//	repo := mock_repository.NewMockIRepository(ctrl)
//	for _, test := range tests {
//		repo.EXPECT().GetBoardById(test.MockArgs.Ctx, test.MockArgs.BoardId).Return(
//			test.MockReturn.BoardByID, test.MockReturn.ErrBoardByID)
//		if test.MockReturn.ErrBoardByID == nil {
//			repo.EXPECT().CheckBoardAuthorExistence(test.MockArgs.Ctx, test.MockArgs.UserId, test.MockArgs.BoardId).Return(
//				test.MockReturn.CheckBoardAuthorExistence, test.MockReturn.ErrCheckBoardAuthorExistence)
//		}
//		if test.MockReturn.ErrCheckBoardAuthorExistence == nil {
//			repo.EXPECT().GetBoardAuthors(test.MockArgs.Ctx, test.MockArgs.BoardId).Return(
//				test.MockReturn.GetBoardAuthors, test.MockReturn.ErrGetBoardAuthors)
//		}
//		if test.MockReturn.ErrGetBoardAuthors == nil {
//			repo.EXPECT().GetBoardPins(test.MockArgs.Ctx, test.MockArgs.BoardId, test.MockArgs.Limit, test.MockArgs.Offset).Return(
//				test.MockReturn.GetBoardPins, test.MockReturn.ErrGetBoardPins)
//		}
//
//		service := service.NewService(repo)
//		board, err := service.GetBoardById(test.FuncArgs.Ctx, test.FuncArgs.BoardId, test.FuncArgs.UserId,
//			test.FuncArgs.Limit, test.FuncArgs.Offset)
//		assert.Equal(t, test.ExpectedFuncReturn.Board, board)
//		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
//	}
//}
//
//func TestUpdateBoard(t *testing.T) {
//	type mockArgs struct {
//		Ctx     context.Context
//		BoardId entity.BoardID
//		Board   entity.Board
//		UserId  entity.UserID
//		Limit   int
//		Offset  int
//	}
//	type mockReturn struct {
//		UpdateBoard                  entity.Board
//		ErrUpdateBoard               error
//		CheckBoardAuthorExistence    bool
//		ErrCheckBoardAuthorExistence error
//		GetBoardAuthors              []entity.BoardAuthor
//		ErrGetBoardAuthors           error
//		GetBoardPins                 []entity.BoardPinResponse
//		ErrGetBoardPins              error
//	}
//	type funcArgs struct {
//		Ctx     context.Context
//		BoardId entity.BoardID
//		UserId  entity.UserID
//		Board   entity.Board
//		Limit   int
//		Offset  int
//	}
//	type funcReturn struct {
//		Board entity.FullBoard
//		Err   errs.ErrorInfo
//	}
//	type test struct {
//		Name               string
//		MockArgs           mockArgs
//		MockReturn         mockReturn
//		FuncArgs           funcArgs
//		ExpectedFuncReturn funcReturn
//	}
//	tests := []test{
//		{
//			Name: "Correct work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				Board:   entity.Board{BoardID: entity.BoardID(1)},
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  0,
//			},
//			MockReturn: mockReturn{
//				UpdateBoard:               entity.Board{BoardID: entity.BoardID(1)},
//				CheckBoardAuthorExistence: true,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Board:   entity.Board{BoardID: entity.BoardID(1)},
//				Limit:   10,
//				Offset:  10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{
//					Board: entity.Board{BoardID: entity.BoardID(1)},
//				},
//			},
//		},
//		{
//			Name: "Uncorrect work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				Board:   entity.Board{BoardID: entity.BoardID(1)},
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  0,
//			},
//			MockReturn: mockReturn{
//				UpdateBoard:               entity.Board{BoardID: entity.BoardID(1)},
//				CheckBoardAuthorExistence: true,
//				ErrGetBoardPins:           errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Board:   entity.Board{BoardID: entity.BoardID(1)},
//				Limit:   10,
//				Offset:  10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{},
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//		},
//		{
//			Name: "Uncorrect work test 2",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				Board:   entity.Board{BoardID: entity.BoardID(1)},
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  0,
//			},
//			MockReturn: mockReturn{
//				UpdateBoard:               entity.Board{BoardID: entity.BoardID(1)},
//				CheckBoardAuthorExistence: true,
//				ErrGetBoardPins:           errs.ErrDBInternal,
//				ErrGetBoardAuthors:        errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Board:   entity.Board{BoardID: entity.BoardID(1)},
//				Limit:   10,
//				Offset:  10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{},
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//		},
//		{
//			Name: "Uncorrect work test 3",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				Board:   entity.Board{BoardID: entity.BoardID(1)},
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  0,
//			},
//			MockReturn: mockReturn{
//				UpdateBoard:               entity.Board{BoardID: entity.BoardID(1)},
//				CheckBoardAuthorExistence: true,
//				ErrGetBoardPins:           errs.ErrDBInternal,
//				ErrGetBoardAuthors:        errs.ErrDBInternal,
//				ErrUpdateBoard:            errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Board:   entity.Board{BoardID: entity.BoardID(1)},
//				Limit:   10,
//				Offset:  10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{},
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//		},
//		{
//			Name: "Uncorrect work test 4",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				Board:   entity.Board{BoardID: entity.BoardID(1)},
//				UserId:  entity.UserID(1),
//				Limit:   10,
//				Offset:  0,
//			},
//			MockReturn: mockReturn{
//				UpdateBoard:                  entity.Board{BoardID: entity.BoardID(1)},
//				CheckBoardAuthorExistence:    true,
//				ErrGetBoardPins:              errs.ErrDBInternal,
//				ErrGetBoardAuthors:           errs.ErrDBInternal,
//				ErrUpdateBoard:               errs.ErrDBInternal,
//				ErrCheckBoardAuthorExistence: errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//				Board:   entity.Board{BoardID: entity.BoardID(1)},
//				Limit:   10,
//				Offset:  10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Board: entity.FullBoard{},
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//		},
//	}
//	ctrl := gomock.NewController(t)
//	repo := mock_repository.NewMockIRepository(ctrl)
//	for _, test := range tests {
//		repo.EXPECT().CheckBoardAuthorExistence(test.MockArgs.Ctx, test.MockArgs.UserId, test.MockArgs.BoardId).Return(
//			test.MockReturn.CheckBoardAuthorExistence, test.MockReturn.ErrCheckBoardAuthorExistence)
//
//		if test.MockReturn.ErrUpdateBoard == nil {
//			repo.EXPECT().GetBoardAuthors(test.MockArgs.Ctx, test.MockArgs.BoardId).Return(
//				test.MockReturn.GetBoardAuthors, test.MockReturn.ErrGetBoardAuthors)
//		}
//		if test.MockReturn.ErrCheckBoardAuthorExistence == nil {
//			repo.EXPECT().UpdateBoard(test.MockArgs.Ctx, test.MockArgs.Board).Return(
//				test.MockReturn.UpdateBoard, test.MockReturn.ErrUpdateBoard)
//		}
//		if test.MockReturn.ErrGetBoardAuthors == nil {
//			repo.EXPECT().GetBoardPins(test.MockArgs.Ctx, test.MockArgs.BoardId, test.MockArgs.Limit, test.MockArgs.Offset).Return(
//				test.MockReturn.GetBoardPins, test.MockReturn.ErrGetBoardPins)
//		}
//
//		service := service.NewService(repo)
//		board, err := service.UpdateBoard(test.FuncArgs.Ctx, test.FuncArgs.Board, test.FuncArgs.UserId)
//		assert.Equal(t, test.ExpectedFuncReturn.Board, board)
//		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
//	}
//}
//
//func TestGetUserBoards(t *testing.T) {
//	type mockArgs struct {
//		Ctx      context.Context
//		Nickname string
//		AuthorId entity.UserID
//		BoardId  entity.BoardID
//		UserId   entity.UserID
//		Limit    int
//		Offset   int
//	}
//	type mockReturn struct {
//		GetUserByNickname    entity.User
//		ErrGetUserByNickname error
//		GetUserBoards        entity.UserBoards
//		ErrGetUserBoards     error
//	}
//	type funcArgs struct {
//		Ctx      context.Context
//		Nickname string
//		AuthorId entity.UserID
//		Limit    int
//		Offset   int
//	}
//	type funcReturn struct {
//		Boards entity.UserBoards
//		Err    errs.ErrorInfo
//	}
//	type test struct {
//		Name                  string
//		MockArgs              mockArgs
//		MockReturn            mockReturn
//		FuncArgs              funcArgs
//		ExpectedFuncReturn    funcReturn
//		WaitGetUserBoardsCall bool
//	}
//	tests := []test{
//		{
//			Name: "Correct work test 1",
//			MockArgs: mockArgs{
//				Ctx:      context.Background(),
//				Nickname: "NICKNAMEMEM",
//				BoardId:  entity.BoardID(1),
//				UserId:   entity.UserID(1),
//				Limit:    10,
//				Offset:   10,
//			},
//			MockReturn: mockReturn{
//				GetUserByNickname:    entity.User{UserID: entity.UserID(1)},
//				ErrGetUserByNickname: nil,
//				GetUserBoards: entity.UserBoards{
//					Boards: []entity.Board{
//						{
//							BoardID:        entity.BoardID(1),
//							VisibilityType: "private",
//						},
//						{
//							BoardID:        entity.BoardID(2),
//							VisibilityType: "public",
//						},
//					},
//				},
//				ErrGetUserBoards: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:      context.Background(),
//				Nickname: "NICKNAMEMEM",
//				AuthorId: entity.UserID(1),
//				Limit:    10,
//				Offset:   10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				Boards: entity.UserBoards{
//					Boards: []entity.Board{
//						{
//							BoardID:        entity.BoardID(1),
//							VisibilityType: "private",
//						},
//						{
//							BoardID:        entity.BoardID(2),
//							VisibilityType: "public",
//						},
//					},
//				},
//				Err: errs.ErrorInfo{},
//			},
//			WaitGetUserBoardsCall: true,
//		},
//		{
//			Name: "Correct work test 2",
//			MockArgs: mockArgs{
//				Ctx:      context.Background(),
//				Nickname: "NICKNAMEMEM",
//				BoardId:  entity.BoardID(1),
//				UserId:   entity.UserID(4),
//				Limit:    10,
//				Offset:   10,
//			},
//			MockReturn: mockReturn{
//				GetUserByNickname:    entity.User{UserID: entity.UserID(4)},
//				ErrGetUserByNickname: nil,
//				GetUserBoards:        entity.UserBoards{},
//				ErrGetUserBoards:     nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:      context.Background(),
//				Nickname: "NICKNAMEMEM",
//				AuthorId: entity.UserID(1),
//				Limit:    10,
//				Offset:   10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				entity.UserBoards{}, errs.ErrorInfo{},
//			},
//			WaitGetUserBoardsCall: true,
//		},
//		{
//			Name: "Uncorrect work test 1",
//			MockArgs: mockArgs{
//				Ctx:      context.Background(),
//				Nickname: "NICKNAMEMEM",
//				BoardId:  entity.BoardID(1),
//				UserId:   entity.UserID(1),
//				Limit:    10,
//				Offset:   10,
//			},
//			MockReturn: mockReturn{
//				GetUserByNickname:    entity.User{},
//				ErrGetUserByNickname: nil,
//				GetUserBoards:        entity.UserBoards{},
//				ErrGetUserBoards:     nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:      context.Background(),
//				Nickname: "NICKNAMEMEM",
//				AuthorId: entity.UserID(1),
//				Limit:    10,
//				Offset:   10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				entity.UserBoards{}, errs.ErrorInfo{
//					LocalErr: errs.ErrUserNotExist,
//				},
//			},
//			WaitGetUserBoardsCall: false,
//		},
//		{
//			Name: "Uncorrect work test 2",
//			MockArgs: mockArgs{
//				Ctx:      context.Background(),
//				Nickname: "NICKNAMEMEM",
//				BoardId:  entity.BoardID(1),
//				UserId:   entity.UserID(1),
//				Limit:    10,
//				Offset:   10,
//			},
//			MockReturn: mockReturn{
//				GetUserByNickname:    entity.User{},
//				ErrGetUserByNickname: errs.ErrDBInternal,
//				GetUserBoards:        entity.UserBoards{},
//				ErrGetUserBoards:     nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:      context.Background(),
//				Nickname: "NICKNAMEMEM",
//				AuthorId: entity.UserID(1),
//				Limit:    10,
//				Offset:   10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				entity.UserBoards{}, errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//			WaitGetUserBoardsCall: false,
//		},
//		{
//			Name: "Uncorrect work test 3",
//			MockArgs: mockArgs{
//				Ctx:      context.Background(),
//				Nickname: "NICKNAMEMEM",
//				BoardId:  entity.BoardID(1),
//				UserId:   entity.UserID(1),
//				Limit:    10,
//				Offset:   10,
//			},
//			MockReturn: mockReturn{
//				GetUserByNickname:    entity.User{UserID: entity.UserID(1)},
//				ErrGetUserByNickname: nil,
//				GetUserBoards:        entity.UserBoards{},
//				ErrGetUserBoards:     errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:      context.Background(),
//				Nickname: "NICKNAMEMEM",
//				AuthorId: entity.UserID(1),
//				Limit:    10,
//				Offset:   10,
//			},
//			ExpectedFuncReturn: funcReturn{
//				entity.UserBoards{}, errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//			WaitGetUserBoardsCall: true,
//		},
//	}
//	ctrl := gomock.NewController(t)
//	repo := mock_repository.NewMockIRepository(ctrl)
//	for _, test := range tests {
//		repo.EXPECT().GetUserByNickname(test.MockArgs.Ctx, test.MockArgs.Nickname).Return(
//			test.MockReturn.GetUserByNickname, test.MockReturn.ErrGetUserByNickname)
//		if test.WaitGetUserBoardsCall {
//			repo.EXPECT().GetUserBoards(test.MockArgs.Ctx, test.MockArgs.UserId, test.MockArgs.Limit, test.MockArgs.Offset).Return(
//				test.MockReturn.GetUserBoards, test.MockReturn.ErrGetUserBoards)
//		}
//		service := service.NewService(repo)
//		boards, err := service.GetUserBoards(test.FuncArgs.Ctx, test.FuncArgs.Nickname, test.FuncArgs.AuthorId,
//			test.FuncArgs.Limit, test.FuncArgs.Offset)
//		assert.Equal(t, test.ExpectedFuncReturn.Boards, boards)
//		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
//	}
//}
//
//func TestDeleteBoard(t *testing.T) {
//	type mockArgs struct {
//		Ctx     context.Context
//		BoardId entity.BoardID
//		UserId  entity.UserID
//	}
//	type mockReturn struct {
//		ErrDelete                    error
//		IsAuthor                     bool
//		ErrCheckBoardAuthorExistence error
//	}
//	type funcArgs struct {
//		Ctx     context.Context
//		BoardId entity.BoardID
//		UserId  entity.UserID
//	}
//	type funcReturn struct {
//		Err errs.ErrorInfo
//	}
//	type test struct {
//		Name               string
//		MockArgs           mockArgs
//		MockReturn         mockReturn
//		FuncArgs           funcArgs
//		ExpectedFuncReturn funcReturn
//		WaitDeleteCall     bool
//	}
//	tests := []test{
//		{
//			Name: "Correct work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				IsAuthor:                     true,
//				ErrDelete:                    nil,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{},
//			},
//			WaitDeleteCall: true,
//		},
//		{
//			Name: "Uncorrect work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				IsAuthor:                     false,
//				ErrDelete:                    nil,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr: errs.ErrPermissionDenied,
//				},
//			},
//			WaitDeleteCall: false,
//		},
//		{
//			Name: "Uncorrect work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				IsAuthor:                     true,
//				ErrDelete:                    errs.ErrDBInternal,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//			WaitDeleteCall: true,
//		},
//		{
//			Name: "Uncorrect work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				IsAuthor:                     false,
//				ErrDelete:                    nil,
//				ErrCheckBoardAuthorExistence: errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					GeneralErr: errs.ErrDBInternal,
//					LocalErr:   errs.ErrDBInternal,
//				},
//			},
//			WaitDeleteCall: false,
//		},
//	}
//	ctrl := gomock.NewController(t)
//	repo := mock_repository.NewMockIRepository(ctrl)
//	for _, test := range tests {
//		repo.EXPECT().CheckBoardAuthorExistence(test.MockArgs.Ctx, test.MockArgs.UserId, test.MockArgs.BoardId).Return(
//			test.MockReturn.IsAuthor, test.MockReturn.ErrCheckBoardAuthorExistence)
//		if test.WaitDeleteCall {
//			repo.EXPECT().DeleteBoard(test.MockArgs.Ctx, test.MockArgs.BoardId).Return(
//				test.MockReturn.ErrDelete)
//		}
//		service := service.NewService(repo)
//		err := service.DeleteBoard(test.FuncArgs.Ctx, test.FuncArgs.BoardId, test.FuncArgs.UserId)
//		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
//	}
//}
//func TestAddPinToBoard(t *testing.T) {
//	type mockArgs struct {
//		Ctx     context.Context
//		BoardId entity.BoardID
//		PinId   entity.PinID
//		UserId  entity.UserID
//	}
//	type mockReturn struct {
//		Exists                       bool
//		ErrCheckExist                error
//		ErrAdding                    error
//		IsAuthor                     bool
//		ErrCheckBoardAuthorExistence error
//	}
//	type funcArgs struct {
//		Ctx     context.Context
//		BoardId entity.BoardID
//		PinId   entity.PinID
//		UserId  entity.UserID
//	}
//	type funcReturn struct {
//		Err errs.ErrorInfo
//	}
//	type test struct {
//		Name                string
//		MockArgs            mockArgs
//		MockReturn          mockReturn
//		FuncArgs            funcArgs
//		ExpectedFuncReturn  funcReturn
//		WaitAddCall         bool
//		WaitCheckAuthorCall bool
//	}
//	tests := []test{
//		{
//			Name: "Correct work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       true,
//				ErrCheckExist:                nil,
//				ErrAdding:                    nil,
//				IsAuthor:                     true,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{},
//			},
//			WaitAddCall:         true,
//			WaitCheckAuthorCall: true,
//		},
//		{
//			Name: "Uncorrect work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       false,
//				ErrCheckExist:                nil,
//				ErrAdding:                    nil,
//				IsAuthor:                     true,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr: errs.ErrElementNotExist,
//				},
//			},
//			WaitAddCall:         false,
//			WaitCheckAuthorCall: false,
//		},
//		{
//			Name: "Uncorrect work test 2",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       false,
//				ErrCheckExist:                errs.ErrDBInternal,
//				ErrAdding:                    nil,
//				IsAuthor:                     true,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr:   errs.ErrDBInternal,
//					GeneralErr: errs.ErrDBInternal,
//				},
//			},
//			WaitAddCall:         false,
//			WaitCheckAuthorCall: false,
//		},
//		{
//			Name: "Uncorrect work test 3",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       true,
//				ErrCheckExist:                nil,
//				ErrAdding:                    nil,
//				IsAuthor:                     false,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr: errs.ErrPermissionDenied,
//				},
//			},
//			WaitAddCall:         false,
//			WaitCheckAuthorCall: true,
//		},
//		{
//			Name: "Uncorrect work test 4",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       true,
//				ErrCheckExist:                nil,
//				ErrAdding:                    nil,
//				IsAuthor:                     false,
//				ErrCheckBoardAuthorExistence: errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr:   errs.ErrDBInternal,
//					GeneralErr: errs.ErrDBInternal,
//				},
//			},
//			WaitAddCall:         false,
//			WaitCheckAuthorCall: true,
//		},
//		{
//			Name: "Uncorrect work test 5",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       true,
//				ErrCheckExist:                nil,
//				ErrAdding:                    errs.ErrDBInternal,
//				IsAuthor:                     true,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr:   errs.ErrDBInternal,
//					GeneralErr: errs.ErrDBInternal,
//				},
//			},
//			WaitAddCall:         true,
//			WaitCheckAuthorCall: true,
//		},
//	}
//	ctrl := gomock.NewController(t)
//	repo := mock_repository.NewMockIRepository(ctrl)
//	for _, test := range tests {
//		repo.EXPECT().CheckPinExistence(test.MockArgs.Ctx, test.MockArgs.PinId).Return(
//			test.MockReturn.Exists, test.MockReturn.ErrCheckExist)
//		if test.WaitCheckAuthorCall {
//			repo.EXPECT().CheckBoardAuthorExistence(test.MockArgs.Ctx, test.MockArgs.UserId, test.MockArgs.BoardId).Return(
//				test.MockReturn.IsAuthor, test.MockReturn.ErrCheckBoardAuthorExistence)
//		}
//		if test.WaitAddCall {
//			repo.EXPECT().AddPinToBoard(test.MockArgs.Ctx, test.MockArgs.BoardId, test.MockArgs.PinId).Return(
//				test.MockReturn.ErrAdding)
//		}
//
//		service := service.NewService(repo)
//		err := service.AddPinToBoard(test.FuncArgs.Ctx, test.FuncArgs.BoardId, test.FuncArgs.PinId, test.MockArgs.UserId)
//		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
//	}
//}
//
//func TestDeletePinFromBoard(t *testing.T) {
//	type mockArgs struct {
//		Ctx     context.Context
//		BoardId entity.BoardID
//		PinId   entity.PinID
//		UserId  entity.UserID
//	}
//	type mockReturn struct {
//		Exists                       bool
//		ErrCheckExist                error
//		ErrDeleting                  error
//		IsAuthor                     bool
//		ErrCheckBoardAuthorExistence error
//	}
//	type funcArgs struct {
//		Ctx     context.Context
//		BoardId entity.BoardID
//		PinId   entity.PinID
//		UserId  entity.UserID
//	}
//	type funcReturn struct {
//		Err errs.ErrorInfo
//	}
//	type test struct {
//		Name                string
//		MockArgs            mockArgs
//		MockReturn          mockReturn
//		FuncArgs            funcArgs
//		ExpectedFuncReturn  funcReturn
//		WaitDeleteCall      bool
//		WaitCheckAuthorCall bool
//	}
//	tests := []test{
//		{
//			Name: "Correct work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       true,
//				ErrCheckExist:                nil,
//				ErrDeleting:                  nil,
//				IsAuthor:                     true,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{},
//			},
//			WaitDeleteCall:      true,
//			WaitCheckAuthorCall: true,
//		},
//		{
//			Name: "Uncorrect work test 1",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       false,
//				ErrCheckExist:                nil,
//				ErrDeleting:                  nil,
//				IsAuthor:                     true,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr: errs.ErrElementNotExist,
//				},
//			},
//			WaitDeleteCall:      false,
//			WaitCheckAuthorCall: false,
//		},
//		{
//			Name: "Uncorrect work test 2",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       false,
//				ErrCheckExist:                errs.ErrDBInternal,
//				ErrDeleting:                  nil,
//				IsAuthor:                     true,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr:   errs.ErrDBInternal,
//					GeneralErr: errs.ErrDBInternal,
//				},
//			},
//			WaitDeleteCall:      false,
//			WaitCheckAuthorCall: false,
//		},
//		{
//			Name: "Uncorrect work test 3",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       true,
//				ErrCheckExist:                nil,
//				ErrDeleting:                  nil,
//				IsAuthor:                     false,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr: errs.ErrPermissionDenied,
//				},
//			},
//			WaitDeleteCall:      false,
//			WaitCheckAuthorCall: true,
//		},
//		{
//			Name: "Uncorrect work test 4",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       true,
//				ErrCheckExist:                nil,
//				ErrDeleting:                  nil,
//				IsAuthor:                     false,
//				ErrCheckBoardAuthorExistence: errs.ErrDBInternal,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr:   errs.ErrDBInternal,
//					GeneralErr: errs.ErrDBInternal,
//				},
//			},
//			WaitDeleteCall:      false,
//			WaitCheckAuthorCall: true,
//		},
//		{
//			Name: "Uncorrect work test 5",
//			MockArgs: mockArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			MockReturn: mockReturn{
//				Exists:                       true,
//				ErrCheckExist:                nil,
//				ErrDeleting:                  errs.ErrDBInternal,
//				IsAuthor:                     true,
//				ErrCheckBoardAuthorExistence: nil,
//			},
//			FuncArgs: funcArgs{
//				Ctx:     context.Background(),
//				PinId:   entity.PinID(1),
//				BoardId: entity.BoardID(1),
//				UserId:  entity.UserID(1),
//			},
//			ExpectedFuncReturn: funcReturn{
//				Err: errs.ErrorInfo{
//					LocalErr:   errs.ErrDBInternal,
//					GeneralErr: errs.ErrDBInternal,
//				},
//			},
//			WaitDeleteCall:      true,
//			WaitCheckAuthorCall: true,
//		},
//	}
//	ctrl := gomock.NewController(t)
//	repo := mock_repository.NewMockIRepository(ctrl)
//	for _, test := range tests {
//		repo.EXPECT().CheckPinExistence(test.MockArgs.Ctx, test.MockArgs.PinId).Return(
//			test.MockReturn.Exists, test.MockReturn.ErrCheckExist)
//		if test.WaitCheckAuthorCall {
//			repo.EXPECT().CheckBoardAuthorExistence(test.MockArgs.Ctx, test.MockArgs.UserId, test.MockArgs.BoardId).Return(
//				test.MockReturn.IsAuthor, test.MockReturn.ErrCheckBoardAuthorExistence)
//		}
//		if test.WaitDeleteCall {
//			repo.EXPECT().DeletePinFromBoard(test.MockArgs.Ctx, test.MockArgs.BoardId, test.MockArgs.PinId).Return(
//				test.MockReturn.ErrDeleting)
//		}
//
//		service := service.NewService(repo)
//		err := service.DeletePinFromBoard(test.FuncArgs.Ctx, test.FuncArgs.BoardId, test.FuncArgs.PinId, test.MockArgs.UserId)
//		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
//	}
//}
//
//func TestAuthorContains(t *testing.T) {
//	authors := []entity.BoardAuthor{
//		{UserId: entity.UserID(1)},
//		{UserId: entity.UserID(2)},
//		{UserId: entity.UserID(3)},
//		{UserId: entity.UserID(4)},
//		{UserId: entity.UserID(6)},
//	}
//	userIdVals := []entity.UserID{entity.UserID(3), entity.UserID(5)}
//	reses := []bool{true, false}
//
//	for testNumber, target := range userIdVals {
//		assert.Equal(t, reses[testNumber], service.AuthorContains(authors, target))
//	}
//}

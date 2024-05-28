package test_service

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/service"
	mock_proto "harmonica/mocks/microservices/like/proto"
	mock_repository "harmonica/mocks/repository"
	"testing"
)

func TestService_CreateBoard(t *testing.T) {
	type Args struct {
		Board  entity.Board
		UserId entity.UserID
	}
	type ExpectedReturn struct {
		FullBoard entity.FullBoard
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		Board  entity.Board
		UserId entity.UserID
	}
	type ExpectedMockReturn struct {
		Board  entity.Board
		User   entity.User
		Error1 error
		Error2 error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().CreateBoard(ctx, mockArgs.Board, mockArgs.UserId).Return(mockReturn.Board, mockReturn.Error1)
		repo.EXPECT().GetUserById(ctx, mockArgs.UserId).Return(mockReturn.User, mockReturn.Error2).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: Args{
				Board:  entity.Board{Title: "board 1"},
				UserId: 1,
			},
			expectedReturn: ExpectedReturn{
				FullBoard: entity.FullBoard{
					Board: entity.Board{Title: "board 1"},
					BoardAuthors: []entity.BoardAuthor{
						{UserId: 1},
					},
				},
			},
			expectedMockArgs: ExpectedMockArgs{
				Board:  entity.Board{Title: "board 1"},
				UserId: 1,
			},
			expectedMockReturn: ExpectedMockReturn{
				Board: entity.Board{Title: "board 1"},
				User:  entity.User{UserID: 1},
			},
		},
		{
			name: "Error test case 1",
			args: Args{
				Board: entity.Board{Title: "board 1"},
			},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{
				Board: entity.Board{Title: "board 1"},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error1: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 2",
			args: Args{
				Board: entity.Board{Title: "board 1"},
			},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{
				Board: entity.Board{Title: "board 1"},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error2: errs.ErrDBInternal,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			likeClient := mock_proto.NewMockLikeClient(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			fullBoard, errInfo := s.CreateBoard(context.Background(), testCase.args.Board, testCase.args.UserId)
			assert.Equal(t, testCase.expectedReturn.FullBoard, fullBoard)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetBoardById(t *testing.T) {
	type Args struct {
		BoardId entity.BoardID
		UserId  entity.UserID
	}
	type ExpectedReturn struct {
		FullBoard entity.FullBoard
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		BoardId entity.BoardID
		UserId  entity.UserID
	}
	type ExpectedMockReturn struct {
		Board        entity.Board
		IsAuthor     bool
		BoardAuthors []entity.BoardAuthor
		BoardPins    []entity.BoardPinResponse
		Error1       error
		Error2       error
		Error3       error
		Error4       error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetBoardById(ctx, mockArgs.BoardId).Return(mockReturn.Board, mockReturn.Error1)
		repo.EXPECT().CheckBoardAuthorExistence(ctx, mockArgs.UserId, mockArgs.BoardId).Return(mockReturn.IsAuthor, mockReturn.Error2).MaxTimes(1)
		repo.EXPECT().GetBoardAuthors(ctx, mockArgs.BoardId).Return(mockReturn.BoardAuthors, mockReturn.Error3).MaxTimes(1)
		repo.EXPECT().GetBoardPins(ctx, mockArgs.BoardId, Limit, Offset).Return(mockReturn.BoardPins, mockReturn.Error4).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: Args{
				BoardId: 100,
				UserId:  1,
			},
			expectedReturn: ExpectedReturn{
				FullBoard: entity.FullBoard{
					Board: entity.Board{
						BoardID: 100,
						IsOwner: true,
					},
					BoardAuthors: []entity.BoardAuthor{
						{UserId: 1},
					},
					Pins: make([]entity.BoardPinResponse, 0),
				},
			},
			expectedMockArgs: ExpectedMockArgs{
				BoardId: 100,
				UserId:  1,
			},
			expectedMockReturn: ExpectedMockReturn{
				Board:    entity.Board{BoardID: 100},
				IsAuthor: true,
				BoardAuthors: []entity.BoardAuthor{
					{UserId: 1},
				},
				BoardPins: []entity.BoardPinResponse{},
			},
		},
		{
			name: "Error test case 1",
			args: Args{
				BoardId: 100,
			},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: sql.ErrNoRows, LocalErr: errs.ErrElementNotExist},
			},
			expectedMockArgs: ExpectedMockArgs{
				BoardId: 100,
			},
			expectedMockReturn: ExpectedMockReturn{
				Error1: sql.ErrNoRows,
			},
		},
		{
			name: "Error test case 2",
			args: Args{
				BoardId: 100,
				UserId:  1,
			},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{LocalErr: errs.ErrPermissionDenied},
			},
			expectedMockArgs: ExpectedMockArgs{
				BoardId: 100,
				UserId:  1,
			},
			expectedMockReturn: ExpectedMockReturn{
				IsAuthor: false,
				Board:    entity.Board{BoardID: 100, VisibilityType: entity.VisibilityPrivate},
			},
		},
		{
			name: "Error test case 3",
			args: Args{
				BoardId: 100,
				UserId:  1,
			},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{
				BoardId: 100,
				UserId:  1,
			},
			expectedMockReturn: ExpectedMockReturn{
				IsAuthor: false,
				Board:    entity.Board{BoardID: 100, VisibilityType: entity.VisibilityPrivate},
				Error2:   errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 4",
			args: Args{
				BoardId: 100,
				UserId:  1,
			},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{
				BoardId: 100,
				UserId:  1,
			},
			expectedMockReturn: ExpectedMockReturn{
				Board:  entity.Board{BoardID: 100},
				Error3: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 5",
			args: Args{
				BoardId: 100,
				UserId:  1,
			},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{
				BoardId: 100,
				UserId:  1,
			},
			expectedMockReturn: ExpectedMockReturn{
				Board:  entity.Board{BoardID: 100},
				Error3: errs.ErrDBInternal,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			likeClient := mock_proto.NewMockLikeClient(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			fullBoard, errInfo := s.GetBoardById(context.Background(), testCase.args.BoardId, testCase.args.UserId, Limit, Offset)
			assert.Equal(t, testCase.expectedReturn.FullBoard, fullBoard)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_UpdateBoard(t *testing.T) {
	type mockArgs struct {
		Ctx     context.Context
		BoardId entity.BoardID
		Board   entity.Board
		UserId  entity.UserID
		Limit   int
		Offset  int
	}
	type mockReturn struct {
		UpdateBoard                  entity.Board
		ErrUpdateBoard               error
		CheckBoardAuthorExistence    bool
		ErrCheckBoardAuthorExistence error
		GetBoardAuthors              []entity.BoardAuthor
		ErrGetBoardAuthors           error
		GetBoardPins                 []entity.BoardPinResponse
		ErrGetBoardPins              error
	}
	type funcArgs struct {
		Ctx     context.Context
		BoardId entity.BoardID
		UserId  entity.UserID
		Board   entity.Board
		Limit   int
		Offset  int
	}
	type funcReturn struct {
		Board entity.FullBoard
		Err   errs.ErrorInfo
	}
	type test struct {
		Name               string
		MockArgs           mockArgs
		MockReturn         mockReturn
		FuncArgs           funcArgs
		ExpectedFuncReturn funcReturn
	}
	tests := []test{
		{
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				Board:   entity.Board{BoardID: entity.BoardID(1)},
				UserId:  entity.UserID(1),
				Limit:   10,
				Offset:  0,
			},
			MockReturn: mockReturn{
				UpdateBoard:               entity.Board{BoardID: entity.BoardID(1)},
				CheckBoardAuthorExistence: true,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
				Board:   entity.Board{BoardID: entity.BoardID(1)},
				Limit:   10,
				Offset:  10,
			},
			ExpectedFuncReturn: funcReturn{
				Board: entity.FullBoard{
					Board: entity.Board{BoardID: entity.BoardID(1)},
				},
			},
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				Board:   entity.Board{BoardID: entity.BoardID(1)},
				UserId:  entity.UserID(1),
				Limit:   10,
				Offset:  0,
			},
			MockReturn: mockReturn{
				UpdateBoard:               entity.Board{BoardID: entity.BoardID(1)},
				CheckBoardAuthorExistence: true,
				ErrGetBoardPins:           errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
				Board:   entity.Board{BoardID: entity.BoardID(1)},
				Limit:   10,
				Offset:  10,
			},
			ExpectedFuncReturn: funcReturn{
				Board: entity.FullBoard{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
		{
			Name: "Uncorrect work test 2",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				Board:   entity.Board{BoardID: entity.BoardID(1)},
				UserId:  entity.UserID(1),
				Limit:   10,
				Offset:  0,
			},
			MockReturn: mockReturn{
				UpdateBoard:               entity.Board{BoardID: entity.BoardID(1)},
				CheckBoardAuthorExistence: true,
				ErrGetBoardPins:           errs.ErrDBInternal,
				ErrGetBoardAuthors:        errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
				Board:   entity.Board{BoardID: entity.BoardID(1)},
				Limit:   10,
				Offset:  10,
			},
			ExpectedFuncReturn: funcReturn{
				Board: entity.FullBoard{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
		{
			Name: "Uncorrect work test 3",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				Board:   entity.Board{BoardID: entity.BoardID(1)},
				UserId:  entity.UserID(1),
				Limit:   10,
				Offset:  0,
			},
			MockReturn: mockReturn{
				UpdateBoard:               entity.Board{BoardID: entity.BoardID(1)},
				CheckBoardAuthorExistence: true,
				ErrGetBoardPins:           errs.ErrDBInternal,
				ErrGetBoardAuthors:        errs.ErrDBInternal,
				ErrUpdateBoard:            errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
				Board:   entity.Board{BoardID: entity.BoardID(1)},
				Limit:   10,
				Offset:  10,
			},
			ExpectedFuncReturn: funcReturn{
				Board: entity.FullBoard{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
		{
			Name: "Uncorrect work test 4",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				Board:   entity.Board{BoardID: entity.BoardID(1)},
				UserId:  entity.UserID(1),
				Limit:   10,
				Offset:  0,
			},
			MockReturn: mockReturn{
				UpdateBoard:                  entity.Board{BoardID: entity.BoardID(1)},
				CheckBoardAuthorExistence:    true,
				ErrGetBoardPins:              errs.ErrDBInternal,
				ErrGetBoardAuthors:           errs.ErrDBInternal,
				ErrUpdateBoard:               errs.ErrDBInternal,
				ErrCheckBoardAuthorExistence: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
				Board:   entity.Board{BoardID: entity.BoardID(1)},
				Limit:   10,
				Offset:  10,
			},
			ExpectedFuncReturn: funcReturn{
				Board: entity.FullBoard{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().CheckBoardAuthorExistence(test.MockArgs.Ctx, test.MockArgs.UserId, test.MockArgs.BoardId).Return(
			test.MockReturn.CheckBoardAuthorExistence, test.MockReturn.ErrCheckBoardAuthorExistence)

		if test.MockReturn.ErrUpdateBoard == nil {
			repo.EXPECT().GetBoardAuthors(test.MockArgs.Ctx, test.MockArgs.BoardId).Return(
				test.MockReturn.GetBoardAuthors, test.MockReturn.ErrGetBoardAuthors)
		}
		if test.MockReturn.ErrCheckBoardAuthorExistence == nil {
			repo.EXPECT().UpdateBoard(test.MockArgs.Ctx, test.MockArgs.Board).Return(
				test.MockReturn.UpdateBoard, test.MockReturn.ErrUpdateBoard)
		}
		if test.MockReturn.ErrGetBoardAuthors == nil {
			repo.EXPECT().GetBoardPins(test.MockArgs.Ctx, test.MockArgs.BoardId, test.MockArgs.Limit, test.MockArgs.Offset).Return(
				test.MockReturn.GetBoardPins, test.MockReturn.ErrGetBoardPins)
		}
		likeClient := mock_proto.NewMockLikeClient(ctrl)
		s := service.NewService(repo, likeClient)
		board, err := s.UpdateBoard(test.FuncArgs.Ctx, test.FuncArgs.Board, test.FuncArgs.UserId)
		assert.Equal(t, test.ExpectedFuncReturn.Board, board)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestGetUserBoards(t *testing.T) {
	type mockArgs struct {
		Ctx      context.Context
		Nickname string
		AuthorId entity.UserID
		BoardId  entity.BoardID
		UserId   entity.UserID
		Limit    int
		Offset   int
	}
	type mockReturn struct {
		GetUserByNickname    entity.User
		ErrGetUserByNickname error
		GetUserBoards        entity.UserBoards
		ErrGetUserBoards     error
	}
	type funcArgs struct {
		Ctx      context.Context
		Nickname string
		AuthorId entity.UserID
		Limit    int
		Offset   int
	}
	type funcReturn struct {
		Boards entity.UserBoards
		Err    errs.ErrorInfo
	}
	type test struct {
		Name                  string
		MockArgs              mockArgs
		MockReturn            mockReturn
		FuncArgs              funcArgs
		ExpectedFuncReturn    funcReturn
		WaitGetUserBoardsCall bool
	}
	tests := []test{
		{
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:      context.Background(),
				Nickname: "NICKNAMEMEM",
				BoardId:  entity.BoardID(1),
				UserId:   entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			MockReturn: mockReturn{
				GetUserByNickname:    entity.User{UserID: entity.UserID(1)},
				ErrGetUserByNickname: nil,
				GetUserBoards: entity.UserBoards{
					Boards: []entity.UserBoard{
						{
							BoardID:        entity.BoardID(1),
							VisibilityType: "private",
						},
						{
							BoardID:        entity.BoardID(2),
							VisibilityType: "public",
						},
					},
				},
				ErrGetUserBoards: nil,
			},
			FuncArgs: funcArgs{
				Ctx:      context.Background(),
				Nickname: "NICKNAMEMEM",
				AuthorId: entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			ExpectedFuncReturn: funcReturn{
				Boards: entity.UserBoards{
					Boards: []entity.UserBoard{
						{
							BoardID:        entity.BoardID(1),
							VisibilityType: "private",
						},
						{
							BoardID:        entity.BoardID(2),
							VisibilityType: "public",
						},
					},
				},
				Err: errs.ErrorInfo{},
			},
			WaitGetUserBoardsCall: true,
		},
		{
			Name: "Correct work test 2",
			MockArgs: mockArgs{
				Ctx:      context.Background(),
				Nickname: "NICKNAMEMEM",
				BoardId:  entity.BoardID(1),
				UserId:   entity.UserID(4),
				Limit:    10,
				Offset:   10,
			},
			MockReturn: mockReturn{
				GetUserByNickname:    entity.User{UserID: entity.UserID(4)},
				ErrGetUserByNickname: nil,
				GetUserBoards:        entity.UserBoards{},
				ErrGetUserBoards:     nil,
			},
			FuncArgs: funcArgs{
				Ctx:      context.Background(),
				Nickname: "NICKNAMEMEM",
				AuthorId: entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			ExpectedFuncReturn: funcReturn{
				entity.UserBoards{}, errs.ErrorInfo{},
			},
			WaitGetUserBoardsCall: true,
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:      context.Background(),
				Nickname: "NICKNAMEMEM",
				BoardId:  entity.BoardID(1),
				UserId:   entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			MockReturn: mockReturn{
				GetUserByNickname:    entity.User{},
				ErrGetUserByNickname: nil,
				GetUserBoards:        entity.UserBoards{},
				ErrGetUserBoards:     nil,
			},
			FuncArgs: funcArgs{
				Ctx:      context.Background(),
				Nickname: "NICKNAMEMEM",
				AuthorId: entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			ExpectedFuncReturn: funcReturn{
				entity.UserBoards{}, errs.ErrorInfo{
					LocalErr: errs.ErrUserNotExist,
				},
			},
			WaitGetUserBoardsCall: false,
		},
		{
			Name: "Uncorrect work test 2",
			MockArgs: mockArgs{
				Ctx:      context.Background(),
				Nickname: "NICKNAMEMEM",
				BoardId:  entity.BoardID(1),
				UserId:   entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			MockReturn: mockReturn{
				GetUserByNickname:    entity.User{},
				ErrGetUserByNickname: errs.ErrDBInternal,
				GetUserBoards:        entity.UserBoards{},
				ErrGetUserBoards:     nil,
			},
			FuncArgs: funcArgs{
				Ctx:      context.Background(),
				Nickname: "NICKNAMEMEM",
				AuthorId: entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			ExpectedFuncReturn: funcReturn{
				entity.UserBoards{}, errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			WaitGetUserBoardsCall: false,
		},
		{
			Name: "Uncorrect work test 3",
			MockArgs: mockArgs{
				Ctx:      context.Background(),
				Nickname: "NICKNAMEMEM",
				BoardId:  entity.BoardID(1),
				UserId:   entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			MockReturn: mockReturn{
				GetUserByNickname:    entity.User{UserID: entity.UserID(1)},
				ErrGetUserByNickname: nil,
				GetUserBoards:        entity.UserBoards{},
				ErrGetUserBoards:     errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:      context.Background(),
				Nickname: "NICKNAMEMEM",
				AuthorId: entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			ExpectedFuncReturn: funcReturn{
				entity.UserBoards{}, errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			WaitGetUserBoardsCall: true,
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().GetUserByNickname(test.MockArgs.Ctx, test.MockArgs.Nickname).Return(
			test.MockReturn.GetUserByNickname, test.MockReturn.ErrGetUserByNickname)
		if test.WaitGetUserBoardsCall {
			repo.EXPECT().GetUserBoards(test.MockArgs.Ctx, gomock.Any(), gomock.Any(), test.MockArgs.Limit, test.MockArgs.Offset).Return(
				test.MockReturn.GetUserBoards, test.MockReturn.ErrGetUserBoards)
		}
		likeClient := mock_proto.NewMockLikeClient(ctrl)
		s := service.NewService(repo, likeClient)
		boards, err := s.GetUserBoards(test.FuncArgs.Ctx, test.FuncArgs.Nickname, test.FuncArgs.AuthorId,
			test.FuncArgs.Limit, test.FuncArgs.Offset)
		assert.Equal(t, test.ExpectedFuncReturn.Boards, boards)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestDeleteBoard(t *testing.T) {
	type mockArgs struct {
		Ctx     context.Context
		BoardId entity.BoardID
		UserId  entity.UserID
	}
	type mockReturn struct {
		ErrDelete                    error
		IsAuthor                     bool
		ErrCheckBoardAuthorExistence error
	}
	type funcArgs struct {
		Ctx     context.Context
		BoardId entity.BoardID
		UserId  entity.UserID
	}
	type funcReturn struct {
		Err errs.ErrorInfo
	}
	type test struct {
		Name               string
		MockArgs           mockArgs
		MockReturn         mockReturn
		FuncArgs           funcArgs
		ExpectedFuncReturn funcReturn
		WaitDeleteCall     bool
	}
	tests := []test{
		{
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				IsAuthor:                     true,
				ErrDelete:                    nil,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{},
			},
			WaitDeleteCall: true,
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				IsAuthor:                     false,
				ErrDelete:                    nil,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr: errs.ErrPermissionDenied,
				},
			},
			WaitDeleteCall: false,
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				IsAuthor:                     true,
				ErrDelete:                    errs.ErrDBInternal,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			WaitDeleteCall: true,
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				IsAuthor:                     false,
				ErrDelete:                    nil,
				ErrCheckBoardAuthorExistence: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			WaitDeleteCall: false,
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().CheckBoardAuthorExistence(test.MockArgs.Ctx, test.MockArgs.UserId, test.MockArgs.BoardId).Return(
			test.MockReturn.IsAuthor, test.MockReturn.ErrCheckBoardAuthorExistence)
		if test.WaitDeleteCall {
			repo.EXPECT().DeleteBoard(test.MockArgs.Ctx, test.MockArgs.BoardId).Return(
				test.MockReturn.ErrDelete)
		}
		likeClient := mock_proto.NewMockLikeClient(ctrl)
		s := service.NewService(repo, likeClient)
		err := s.DeleteBoard(test.FuncArgs.Ctx, test.FuncArgs.BoardId, test.FuncArgs.UserId)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}
func TestAddPinToBoard(t *testing.T) {
	type mockArgs struct {
		Ctx     context.Context
		BoardId entity.BoardID
		PinId   entity.PinID
		UserId  entity.UserID
	}
	type mockReturn struct {
		Exists                       bool
		ErrCheckExist                error
		ErrAdding                    error
		IsAuthor                     bool
		ErrCheckBoardAuthorExistence error
	}
	type funcArgs struct {
		Ctx     context.Context
		BoardId entity.BoardID
		PinId   entity.PinID
		UserId  entity.UserID
	}
	type funcReturn struct {
		Err errs.ErrorInfo
	}
	type test struct {
		Name                string
		MockArgs            mockArgs
		MockReturn          mockReturn
		FuncArgs            funcArgs
		ExpectedFuncReturn  funcReturn
		WaitAddCall         bool
		WaitCheckAuthorCall bool
	}
	tests := []test{
		{
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       true,
				ErrCheckExist:                nil,
				ErrAdding:                    nil,
				IsAuthor:                     true,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{},
			},
			WaitAddCall:         true,
			WaitCheckAuthorCall: true,
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       false,
				ErrCheckExist:                nil,
				ErrAdding:                    nil,
				IsAuthor:                     true,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr: errs.ErrElementNotExist,
				},
			},
			WaitAddCall:         false,
			WaitCheckAuthorCall: false,
		},
		{
			Name: "Uncorrect work test 2",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       false,
				ErrCheckExist:                errs.ErrDBInternal,
				ErrAdding:                    nil,
				IsAuthor:                     true,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr:   errs.ErrDBInternal,
					GeneralErr: errs.ErrDBInternal,
				},
			},
			WaitAddCall:         false,
			WaitCheckAuthorCall: false,
		},
		{
			Name: "Uncorrect work test 3",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       true,
				ErrCheckExist:                nil,
				ErrAdding:                    nil,
				IsAuthor:                     false,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr: errs.ErrPermissionDenied,
				},
			},
			WaitAddCall:         false,
			WaitCheckAuthorCall: true,
		},
		{
			Name: "Uncorrect work test 4",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       true,
				ErrCheckExist:                nil,
				ErrAdding:                    nil,
				IsAuthor:                     false,
				ErrCheckBoardAuthorExistence: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr:   errs.ErrDBInternal,
					GeneralErr: errs.ErrDBInternal,
				},
			},
			WaitAddCall:         false,
			WaitCheckAuthorCall: true,
		},
		{
			Name: "Uncorrect work test 5",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       true,
				ErrCheckExist:                nil,
				ErrAdding:                    errs.ErrDBInternal,
				IsAuthor:                     true,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr:   errs.ErrDBInternal,
					GeneralErr: errs.ErrDBInternal,
				},
			},
			WaitAddCall:         true,
			WaitCheckAuthorCall: true,
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().CheckPinExistence(test.MockArgs.Ctx, test.MockArgs.PinId).Return(
			test.MockReturn.Exists, test.MockReturn.ErrCheckExist)
		if test.WaitCheckAuthorCall {
			repo.EXPECT().CheckBoardAuthorExistence(test.MockArgs.Ctx, test.MockArgs.UserId, test.MockArgs.BoardId).Return(
				test.MockReturn.IsAuthor, test.MockReturn.ErrCheckBoardAuthorExistence)
		}
		if test.WaitAddCall {
			repo.EXPECT().AddPinToBoard(test.MockArgs.Ctx, test.MockArgs.BoardId, test.MockArgs.PinId).Return(
				test.MockReturn.ErrAdding)
		}
		likeClient := mock_proto.NewMockLikeClient(ctrl)
		s := service.NewService(repo, likeClient)
		err := s.AddPinToBoard(test.FuncArgs.Ctx, test.FuncArgs.BoardId, test.FuncArgs.PinId, test.MockArgs.UserId)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestDeletePinFromBoard(t *testing.T) {
	type mockArgs struct {
		Ctx     context.Context
		BoardId entity.BoardID
		PinId   entity.PinID
		UserId  entity.UserID
	}
	type mockReturn struct {
		Exists                       bool
		ErrCheckExist                error
		ErrDeleting                  error
		IsAuthor                     bool
		ErrCheckBoardAuthorExistence error
	}
	type funcArgs struct {
		Ctx     context.Context
		BoardId entity.BoardID
		PinId   entity.PinID
		UserId  entity.UserID
	}
	type funcReturn struct {
		Err errs.ErrorInfo
	}
	type test struct {
		Name                string
		MockArgs            mockArgs
		MockReturn          mockReturn
		FuncArgs            funcArgs
		ExpectedFuncReturn  funcReturn
		WaitDeleteCall      bool
		WaitCheckAuthorCall bool
	}
	tests := []test{
		{
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       true,
				ErrCheckExist:                nil,
				ErrDeleting:                  nil,
				IsAuthor:                     true,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{},
			},
			WaitDeleteCall:      true,
			WaitCheckAuthorCall: true,
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       false,
				ErrCheckExist:                nil,
				ErrDeleting:                  nil,
				IsAuthor:                     true,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr: errs.ErrElementNotExist,
				},
			},
			WaitDeleteCall:      false,
			WaitCheckAuthorCall: false,
		},
		{
			Name: "Uncorrect work test 2",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       false,
				ErrCheckExist:                errs.ErrDBInternal,
				ErrDeleting:                  nil,
				IsAuthor:                     true,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr:   errs.ErrDBInternal,
					GeneralErr: errs.ErrDBInternal,
				},
			},
			WaitDeleteCall:      false,
			WaitCheckAuthorCall: false,
		},
		{
			Name: "Uncorrect work test 3",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       true,
				ErrCheckExist:                nil,
				ErrDeleting:                  nil,
				IsAuthor:                     false,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr: errs.ErrPermissionDenied,
				},
			},
			WaitDeleteCall:      false,
			WaitCheckAuthorCall: true,
		},
		{
			Name: "Uncorrect work test 4",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       true,
				ErrCheckExist:                nil,
				ErrDeleting:                  nil,
				IsAuthor:                     false,
				ErrCheckBoardAuthorExistence: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr:   errs.ErrDBInternal,
					GeneralErr: errs.ErrDBInternal,
				},
			},
			WaitDeleteCall:      false,
			WaitCheckAuthorCall: true,
		},
		{
			Name: "Uncorrect work test 5",
			MockArgs: mockArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Exists:                       true,
				ErrCheckExist:                nil,
				ErrDeleting:                  errs.ErrDBInternal,
				IsAuthor:                     true,
				ErrCheckBoardAuthorExistence: nil,
			},
			FuncArgs: funcArgs{
				Ctx:     context.Background(),
				PinId:   entity.PinID(1),
				BoardId: entity.BoardID(1),
				UserId:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					LocalErr:   errs.ErrDBInternal,
					GeneralErr: errs.ErrDBInternal,
				},
			},
			WaitDeleteCall:      true,
			WaitCheckAuthorCall: true,
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().CheckPinExistence(test.MockArgs.Ctx, test.MockArgs.PinId).Return(
			test.MockReturn.Exists, test.MockReturn.ErrCheckExist)
		if test.WaitCheckAuthorCall {
			repo.EXPECT().CheckBoardAuthorExistence(test.MockArgs.Ctx, test.MockArgs.UserId, test.MockArgs.BoardId).Return(
				test.MockReturn.IsAuthor, test.MockReturn.ErrCheckBoardAuthorExistence)
		}
		if test.WaitDeleteCall {
			repo.EXPECT().DeletePinFromBoard(test.MockArgs.Ctx, test.MockArgs.BoardId, test.MockArgs.PinId).Return(
				test.MockReturn.ErrDeleting)
		}
		likeClient := mock_proto.NewMockLikeClient(ctrl)
		s := service.NewService(repo, likeClient)
		err := s.DeletePinFromBoard(test.FuncArgs.Ctx, test.FuncArgs.BoardId, test.FuncArgs.PinId, test.MockArgs.UserId)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestService_GetUserBoardsWithoutPin(t *testing.T) {
	type ExpectedReturn struct {
		Boards    entity.UserBoardsWithoutPin
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Boards entity.UserBoardsWithoutPin
		Error  error
	}
	userId := entity.UserID(2)
	pinId := entity.PinID(3)
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserBoardsWithoutPin(ctx, pinId, userId).Return(mockReturn.Boards, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			expectedReturn: ExpectedReturn{
				Boards: entity.UserBoardsWithoutPin{
					Boards: []entity.UserBoardWithoutPin{
						{BoardID: 4},
					},
				},
			},
			expectedMockReturn: ExpectedMockReturn{
				Boards: entity.UserBoardsWithoutPin{
					Boards: []entity.UserBoardWithoutPin{
						{BoardID: 4},
					},
				},
			},
		},
		{
			name: "Error test case 1",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errs.ErrDBInternal,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			likeClient := mock_proto.NewMockLikeClient(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			boards, errInfo := s.GetUserBoardsWithoutPin(context.Background(), pinId, userId)
			assert.Equal(t, testCase.expectedReturn.Boards, boards)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

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

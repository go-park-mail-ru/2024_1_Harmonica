package handler

//
//import (
//	"fmt"
//	"harmonica/internal/entity"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//var pageToLimitAndOffsetTests = []struct {
//	name      string
//	in        int
//	excpected []int
//}{
//	{
//		name:      "Success test 1",
//		in:        0,
//		excpected: []int{Limit, 0},
//	},
//	{
//		name:      "Success test 2",
//		in:        5,
//		excpected: []int{Limit, Limit * 5},
//	},
//}
//
//func TestPageToLimitAndOffset(t *testing.T) {
//	for _, test := range pageToLimitAndOffsetTests {
//		limit, offset := PageToLimitAndOffset(test.in)
//		assert.Equal(t, test.excpected[0], limit)
//		assert.Equal(t, test.excpected[1], offset)
//	}
//}
//
//func TestPinsList(t *testing.T) {
//	t.Run("valid page test 1", func(t *testing.T) {
//		req, err := http.NewRequest("GET", "/pins?page=1", nil)
//		if err != nil {
//			t.Errorf("error: %v", err)
//		}
//		mock.On("GetPins", Limit, Limit*1).Return(entity.Pins{}, nil)
//		recorder := httptest.NewRecorder()
//		handler.PinsList(recorder, req)
//		assert.Equal(t, http.StatusOK, recorder.Code)
//		assert.Equal(t, recorder.Body.String(), `{"pins":null}`)
//	})
//
//	t.Run("invalid page test 1", func(t *testing.T) {
//		req, err := http.NewRequest("GET", "/pins?page=sth_bad", nil)
//		if err != nil {
//			t.Errorf("error: %v", err)
//		}
//		recorder := httptest.NewRecorder()
//		handler.PinsList(recorder, req)
//		assert.Equal(t, ErrorCodes[ErrReadingRequestBody].HttpCode, recorder.Code)
//		assert.Equal(t, recorder.Body.String(), fmt.Sprintf(`{"code":%d,"message":"%s"}`,
//			ErrorCodes[ErrReadingRequestBody].LocalCode, ErrReadingRequestBody.Error()))
//	})
//
//	t.Run("Empty page test", func(t *testing.T) {
//		req, err := http.NewRequest("GET", "/pins", nil)
//		if err != nil {
//			t.Errorf("error: %v", err)
//		}
//		mock.On("GetPins", Limit, Limit*0).Return(entity.Pins{}, ErrDBInternal) // тут не нужен
//		recorder := httptest.NewRecorder()
//		handler.PinsList(recorder, req)
//		assert.Equal(t, ErrorCodes[ErrDBInternal].HttpCode, recorder.Code)
//		assert.Equal(t, recorder.Body.String(), fmt.Sprintf(`{"code":%d,"message":"%s"}`,
//			ErrorCodes[ErrDBInternal].LocalCode, ErrDBInternal.Error()))
//	})
//}

package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"harmonica/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPinsList(t *testing.T) {
	t.Run("valid page", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/pins?page=1", nil)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		mock.On("GetPins", Limit, Limit*1).Return(models.Pins{}, nil)
		recorder := httptest.NewRecorder()
		handler.PinsList(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, recorder.Body.String(), `{"pins":null}`)
	})

	t.Run("invalid page", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/pins?page=sth_bad", nil)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		//mock.On("GetPins", 10, 10 * 1).Return(models.Pins{}, nil) // тут не нужен
		recorder := httptest.NewRecorder()
		handler.PinsList(recorder, req)
		assert.Equal(t, ErrorCodes[ErrReadingRequestBody].HttpCode, recorder.Code)
		assert.Equal(t, recorder.Body.String(), fmt.Sprintf(`{"code":%d,"message":"%s"}`,
			ErrorCodes[ErrReadingRequestBody].LocalCode, ErrReadingRequestBody.Error()))
	})

	t.Run("empty page", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/pins", nil)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		mock.On("GetPins", Limit, Limit*0).Return(models.Pins{}, ErrDBInternal) // тут не нужен
		recorder := httptest.NewRecorder()
		handler.PinsList(recorder, req)
		assert.Equal(t, ErrorCodes[ErrDBInternal].HttpCode, recorder.Code)
		assert.Equal(t, recorder.Body.String(), fmt.Sprintf(`{"code":%d,"message":"%s"}`,
			ErrorCodes[ErrDBInternal].LocalCode, ErrDBInternal.Error()))
	})
}

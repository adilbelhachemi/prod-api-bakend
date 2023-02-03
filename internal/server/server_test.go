package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pratbacknd/internal/category"
	"pratbacknd/internal/product"
	"pratbacknd/internal/storage"
	"pratbacknd/internal/utils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_CreateProduct(t *testing.T) {
	// GIVEN

	// storage mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedStorage := storage.NewMockStorage(ctrl)
	mockedStorage.EXPECT().CreateProduct(gomock.Any()).Return(nil)

	// uuid gen mock
	uuidctrl := gomock.NewController(t)
	defer uuidctrl.Finish()

	mockedUUID := utils.NewMockUUIDGenerator(uuidctrl)
	mockedUUID.EXPECT().Generate().Return("ABC123")

	// server
	testServer, err := New(Config{
		AllowedOrigins: "*",
		Storage:        mockedStorage,
		UUIDGen:        mockedUUID,
	})
	assert.NoError(t, err, "building a server should not return an error")

	recorder := httptest.NewRecorder()
	inputProduct := product.Product{
		Name:             "test",
		ShortDescription: "short description",
	}
	jsonProduct, err := json.Marshal(inputProduct)
	assert.NoError(t, err, "building a server should not return an error")

	req := httptest.NewRequest("POST", "/admin/products", bytes.NewReader(jsonProduct))

	// WHEN
	testServer.Mux.ServeHTTP(recorder, req)

	// THEN
	assert.Equal(t, http.StatusOK, recorder.Code)

	expectedPayload := `{"id":"ABC123","name":"test","image":"","shortDescription":"short description","description":"","priceVatExcluded":{"money":null,"display":""},"vat":{"money":null,"display":""},"totalPrice":{"money":null,"display":""}}`
	assert.Equal(
		t,
		expectedPayload,
		recorder.Body.String(),
	)
}

func Test_CreateCategory(t *testing.T) {
	// GIVEN

	// storage mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedStorage := storage.NewMockStorage(ctrl)
	mockedStorage.EXPECT().CreateCategory(gomock.Any()).Return(nil)

	// uuid gen mock
	uuidctrl := gomock.NewController(t)
	defer uuidctrl.Finish()

	mockedUUID := utils.NewMockUUIDGenerator(uuidctrl)
	mockedUUID.EXPECT().Generate().Return("ABC123")

	// server
	testServer, err := New(Config{
		AllowedOrigins: "*",
		Storage:        mockedStorage,
		UUIDGen:        mockedUUID,
	})
	assert.NoError(t, err, "building a server should not return an error")

	recorder := httptest.NewRecorder()
	inputCategory := category.Category{
		Name:        "test category",
		Description: "category description",
	}
	jsonProduct, err := json.Marshal(inputCategory)
	assert.NoError(t, err, "building a server should not return an error")

	req := httptest.NewRequest("POST", "/admin/categories", bytes.NewReader(jsonProduct))

	// WHEN
	testServer.Mux.ServeHTTP(recorder, req)

	// THEN
	assert.Equal(t, http.StatusOK, recorder.Code)

	expectedPayload, err := json.Marshal(category.Category{
		ID:          "ABC123",
		Name:        inputCategory.Name,
		Description: inputCategory.Description,
	})
	assert.NoError(t, err, "no error should be fired when marchalling category")

	assert.Equal(
		t,
		string(expectedPayload),
		recorder.Body.String(),
	)
}

func TestServer_Categories(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedStorage := storage.NewMockStorage(ctrl)
	mockedResp := []category.Category{
		{
			ID:          "11",
			Name:        "Test",
			Description: "this the first category",
		},
	}
	mockedStorage.EXPECT().Categories().Return(mockedResp, nil)

	// server
	testServer, err := New(Config{
		AllowedOrigins: "*",
		Storage:        mockedStorage,
		UUIDGen:        nil,
	})
	assert.NoError(t, err, "building a server should not return an error")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/categories", nil)
	assert.NoError(t, err, "no error should when building a request")

	// When
	testServer.Mux.ServeHTTP(recorder, req)

	// Then
	assert.Equal(t, http.StatusOK, recorder.Code)

	expectedBody, err := json.Marshal(mockedResp)
	assert.NoError(t, err, "no error should happen when marshalling the response")
	assert.Equal(t, expectedBody, recorder.Body.Bytes())

}

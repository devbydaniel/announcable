package testutil

import (
	"context"
	"io"

	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	"github.com/devbydaniel/announcable/internal/objstore"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog"
)

// MockObjStore is a mock implementation of ObjStore for testing
type MockObjStore struct {
	GetImageUrlFunc   func(bucket, path string) (string, error)
	UpdateImageFunc   func(bucket, path string, img *io.Reader) error
	DeleteImageFunc   func(bucket, path string) error
	GetImageUrlCalls  []MockObjStoreCall
	UpdateImageCalls  []MockObjStoreCall
	DeleteImageCalls  []MockObjStoreCall
}

// MockObjStoreCall records a method call
type MockObjStoreCall struct {
	Bucket string
	Path   string
	Img    *io.Reader
}

// GetImageUrl mocks the GetImageUrl method
func (m *MockObjStore) GetImageUrl(bucket, path string) (string, error) {
	m.GetImageUrlCalls = append(m.GetImageUrlCalls, MockObjStoreCall{
		Bucket: bucket,
		Path:   path,
	})
	if m.GetImageUrlFunc != nil {
		return m.GetImageUrlFunc(bucket, path)
	}
	return "https://example.com/mock-image.jpg", nil
}

// UpdateImage mocks the UpdateImage method
func (m *MockObjStore) UpdateImage(bucket, path string, img *io.Reader) error {
	m.UpdateImageCalls = append(m.UpdateImageCalls, MockObjStoreCall{
		Bucket: bucket,
		Path:   path,
		Img:    img,
	})
	if m.UpdateImageFunc != nil {
		return m.UpdateImageFunc(bucket, path, img)
	}
	return nil
}

// DeleteImage mocks the DeleteImage method
func (m *MockObjStore) DeleteImage(bucket, path string) error {
	m.DeleteImageCalls = append(m.DeleteImageCalls, MockObjStoreCall{
		Bucket: bucket,
		Path:   path,
	})
	if m.DeleteImageFunc != nil {
		return m.DeleteImageFunc(bucket, path)
	}
	return nil
}

// NewMockObjStore creates a new MockObjStore
func NewMockObjStore() *MockObjStore {
	return &MockObjStore{
		GetImageUrlCalls:  []MockObjStoreCall{},
		UpdateImageCalls:  []MockObjStoreCall{},
		DeleteImageCalls:  []MockObjStoreCall{},
	}
}

// MockDependencies creates mock dependencies for testing handlers
type MockDependencies struct {
	DB          *database.DB
	MockObjStore *MockObjStore
	ObjStore    *objstore.ObjStore
	Log         *zerolog.Logger
	Decoder     *schema.Decoder
}

// NewMockDependencies creates a new MockDependencies with test database
func NewMockDependencies(db *database.DB) *MockDependencies {
	logger := log
	mockObjStore := NewMockObjStore()
	// Create a nil ObjStore for testing - methods won't be called directly
	// The release notes service will use the repository which accepts *objstore.ObjStore
	objStore := &objstore.ObjStore{Client: nil}

	return &MockDependencies{
		DB:           db,
		MockObjStore: mockObjStore,
		ObjStore:     objStore,
		Log:          &logger,
		Decoder:      schema.NewDecoder(),
	}
}

// MockEmailSender is a mock implementation for email sending in tests
type MockEmailSender struct {
	SentEmails []MockEmail
	SendFunc   func(to, subject, body string) error
}

// MockEmail records an email that was sent
type MockEmail struct {
	To      string
	Subject string
	Body    string
}

// Send mocks sending an email
func (m *MockEmailSender) Send(to, subject, body string) error {
	m.SentEmails = append(m.SentEmails, MockEmail{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	if m.SendFunc != nil {
		return m.SendFunc(to, subject, body)
	}
	return nil
}

// NewMockEmailSender creates a new MockEmailSender
func NewMockEmailSender() *MockEmailSender {
	return &MockEmailSender{
		SentEmails: []MockEmail{},
	}
}

// ToSharedDependencies converts MockDependencies to shared.Dependencies
func (m *MockDependencies) ToSharedDependencies() *shared.Dependencies {
	return &shared.Dependencies{
		DB:       m.DB,
		ObjStore: m.ObjStore,
		Log:      m.Log,
		Decoder:  m.Decoder,
	}
}

// MockContext returns a basic context for testing
func MockContext() context.Context {
	return context.Background()
}

// AsObjStore converts MockObjStore to objstore.ObjStore interface
// Note: This is a helper for type conversion in tests
func AsObjStore(mock *MockObjStore) *objstore.ObjStore {
	// This function exists to help with type conversion in tests
	// In real usage, you would use the mock directly where the interface is expected
	// For now, return nil as this is meant to be used via interface matching
	return nil
}

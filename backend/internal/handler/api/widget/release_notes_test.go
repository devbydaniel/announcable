package widget

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devbydaniel/announcable/internal/domain/organisation"
	releasenotes "github.com/devbydaniel/announcable/internal/domain/release-notes"
	"github.com/devbydaniel/announcable/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleReleaseNotesServe_Success(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	err := testDB.DB.Client.AutoMigrate(&organisation.Organisation{}, &releasenotes.ReleaseNote{})
	require.NoError(t, err)

	// Create test organization
	testOrg, err := organisation.New("Test Org")
	require.NoError(t, err)
	err = testDB.DB.Client.Create(testOrg).Error
	require.NoError(t, err)

	// Create test release notes
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(testDB.DB, deps.ObjStore))
	testUserID := uuid.New()
	releaseDate := "2024-01-15"

	testReleaseNote := &releasenotes.ReleaseNote{
		OrganisationID:   testOrg.ID,
		Title:            "Test Release Note",
		DescriptionShort: "This is a test release note",
		DescriptionLong:  "This is a longer description for the test release note",
		ReleaseDate:      &releaseDate,
		CreatedBy:        testUserID,
		LastUpdatedBy:    testUserID,
	}

	releaseNoteID, err := releaseNotesService.Create(testReleaseNote, nil)
	require.NoError(t, err)

	// Publish the release note
	err = releaseNotesService.ChangePublishedStatus(releaseNoteID, true)
	require.NoError(t, err)

	// Update testReleaseNote with the created ID for assertions
	testReleaseNote.ID = releaseNoteID

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request with chi context for URL params
	req := httptest.NewRequest(http.MethodGet, "/api/widget/"+testOrg.ExternalID.String()+"/release-notes", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orgId", testOrg.ExternalID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Parse response body
	var response serveReleaseNotesWidgetResponseBody
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// Verify response contains the release note
	assert.Len(t, response.Data, 1)
	assert.Equal(t, releaseNoteID.String(), response.Data[0].ID)
	assert.Equal(t, "Test Release Note", response.Data[0].Title)
	assert.Equal(t, "This is a test release note", response.Data[0].Text)
	assert.Equal(t, "15.01.2024", response.Data[0].Date)
}

func TestHandleReleaseNotesServe_WithPagination(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	err := testDB.DB.Client.AutoMigrate(&organisation.Organisation{}, &releasenotes.ReleaseNote{})
	require.NoError(t, err)

	// Create test organization
	testOrg, err := organisation.New("Test Org")
	require.NoError(t, err)
	err = testDB.DB.Client.Create(testOrg).Error
	require.NoError(t, err)

	// Create multiple test release notes
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(testDB.DB, deps.ObjStore))
	testUserID := uuid.New()

	for i := 1; i <= 5; i++ {
		releaseDate := "2024-01-15"
		rn := &releasenotes.ReleaseNote{
			OrganisationID:   testOrg.ID,
			Title:            "Test Release Note",
			DescriptionShort: "Description",
			DescriptionLong:  "Long description",
			ReleaseDate:      &releaseDate,
			CreatedBy:        testUserID,
			LastUpdatedBy:    testUserID,
		}
		rnID, err := releaseNotesService.Create(rn, nil)
		require.NoError(t, err)

		// Publish the release note
		err = releaseNotesService.ChangePublishedStatus(rnID, true)
		require.NoError(t, err)
	}

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request with pagination params
	req := httptest.NewRequest(http.MethodGet, "/api/widget/"+testOrg.ExternalID.String()+"/release-notes?page=1&pageSize=2", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orgId", testOrg.ExternalID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response body
	var response serveReleaseNotesWidgetResponseBody
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// Verify pagination worked - should only get 2 items
	assert.Len(t, response.Data, 2)
}

func TestHandleReleaseNotesServe_WithForWidget(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	err := testDB.DB.Client.AutoMigrate(&organisation.Organisation{}, &releasenotes.ReleaseNote{})
	require.NoError(t, err)

	// Create test organization
	testOrg, err := organisation.New("Test Org")
	require.NoError(t, err)
	err = testDB.DB.Client.Create(testOrg).Error
	require.NoError(t, err)

	// Create release note that is hidden on widget
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(testDB.DB, deps.ObjStore))
	testUserID := uuid.New()
	releaseDate := "2024-01-15"

	rn := &releasenotes.ReleaseNote{
		OrganisationID:   testOrg.ID,
		Title:            "Test Release Note",
		DescriptionShort: "Description",
		DescriptionLong:  "Long description",
		ReleaseDate:      &releaseDate,
		CreatedBy:        testUserID,
		LastUpdatedBy:    testUserID,
	}
	rnID, err := releaseNotesService.Create(rn, nil)
	require.NoError(t, err)

	// Publish and hide on widget
	err = releaseNotesService.ChangePublishedStatus(rnID, true)
	require.NoError(t, err)

	// Update to hide on widget
	testDB.DB.Client.Model(&releasenotes.ReleaseNote{}).
		Where("id = ?", rnID).
		Update("hide_on_widget", true)

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request with for=widget param
	req := httptest.NewRequest(http.MethodGet, "/api/widget/"+testOrg.ExternalID.String()+"/release-notes?for=widget", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orgId", testOrg.ExternalID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response body
	var response serveReleaseNotesWidgetResponseBody
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// Should not include the release note hidden on widget
	assert.Len(t, response.Data, 0)
}

func TestHandleReleaseNotesServe_WithForWebsite(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	err := testDB.DB.Client.AutoMigrate(&organisation.Organisation{}, &releasenotes.ReleaseNote{})
	require.NoError(t, err)

	// Create test organization
	testOrg, err := organisation.New("Test Org")
	require.NoError(t, err)
	err = testDB.DB.Client.Create(testOrg).Error
	require.NoError(t, err)

	// Create release note that is hidden on release page
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(testDB.DB, deps.ObjStore))
	testUserID := uuid.New()
	releaseDate := "2024-01-15"

	rn := &releasenotes.ReleaseNote{
		OrganisationID:   testOrg.ID,
		Title:            "Test Release Note",
		DescriptionShort: "Description",
		DescriptionLong:  "Long description",
		ReleaseDate:      &releaseDate,
		CreatedBy:        testUserID,
		LastUpdatedBy:    testUserID,
	}
	rnID, err := releaseNotesService.Create(rn, nil)
	require.NoError(t, err)

	// Publish and hide on release page
	err = releaseNotesService.ChangePublishedStatus(rnID, true)
	require.NoError(t, err)

	// Update to hide on release page
	testDB.DB.Client.Model(&releasenotes.ReleaseNote{}).
		Where("id = ?", rnID).
		Update("hide_on_release_page", true)

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request with for=website param
	req := httptest.NewRequest(http.MethodGet, "/api/widget/"+testOrg.ExternalID.String()+"/release-notes?for=website", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orgId", testOrg.ExternalID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response body
	var response serveReleaseNotesWidgetResponseBody
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// Should not include the release note hidden on release page
	assert.Len(t, response.Data, 0)
}

func TestHandleReleaseNotesServe_InvalidPage(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	err := testDB.DB.Client.AutoMigrate(&organisation.Organisation{})
	require.NoError(t, err)

	// Create test organization
	testOrg, err := organisation.New("Test Org")
	require.NoError(t, err)
	err = testDB.DB.Client.Create(testOrg).Error
	require.NoError(t, err)

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request with invalid page param
	req := httptest.NewRequest(http.MethodGet, "/api/widget/"+testOrg.ExternalID.String()+"/release-notes?page=invalid", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orgId", testOrg.ExternalID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error getting release notes")
}

func TestHandleReleaseNotesServe_InvalidPageSize(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	err := testDB.DB.Client.AutoMigrate(&organisation.Organisation{})
	require.NoError(t, err)

	// Create test organization
	testOrg, err := organisation.New("Test Org")
	require.NoError(t, err)
	err = testDB.DB.Client.Create(testOrg).Error
	require.NoError(t, err)

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request with invalid pageSize param
	req := httptest.NewRequest(http.MethodGet, "/api/widget/"+testOrg.ExternalID.String()+"/release-notes?pageSize=invalid", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orgId", testOrg.ExternalID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error getting release notes")
}

func TestHandleReleaseNotesServe_MissingOrgId(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request without orgId in URL params
	req := httptest.NewRequest(http.MethodGet, "/api/widget/release-notes", nil)
	rctx := chi.NewRouteContext()
	// Don't add orgId to URLParams
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error getting release notes")
}

func TestHandleReleaseNotesServe_OrganisationNotFound(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	err := testDB.DB.Client.AutoMigrate(&organisation.Organisation{})
	require.NoError(t, err)

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request with non-existent orgId
	nonExistentOrgId := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/widget/"+nonExistentOrgId.String()+"/release-notes", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orgId", nonExistentOrgId.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error getting widget config")
}

func TestHandleReleaseNotesServe_OnlyPublishedReleaseNotes(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	err := testDB.DB.Client.AutoMigrate(&organisation.Organisation{}, &releasenotes.ReleaseNote{})
	require.NoError(t, err)

	// Create test organization
	testOrg, err := organisation.New("Test Org")
	require.NoError(t, err)
	err = testDB.DB.Client.Create(testOrg).Error
	require.NoError(t, err)

	// Create published and unpublished release notes
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(testDB.DB, deps.ObjStore))
	testUserID := uuid.New()
	releaseDate := "2024-01-15"

	// Create published release note
	publishedRN := &releasenotes.ReleaseNote{
		OrganisationID:   testOrg.ID,
		Title:            "Published Release Note",
		DescriptionShort: "Published description",
		DescriptionLong:  "Published long description",
		ReleaseDate:      &releaseDate,
		CreatedBy:        testUserID,
		LastUpdatedBy:    testUserID,
	}
	publishedRNID, err := releaseNotesService.Create(publishedRN, nil)
	require.NoError(t, err)
	err = releaseNotesService.ChangePublishedStatus(publishedRNID, true)
	require.NoError(t, err)

	// Create unpublished release note
	unpublishedRN := &releasenotes.ReleaseNote{
		OrganisationID:   testOrg.ID,
		Title:            "Unpublished Release Note",
		DescriptionShort: "Unpublished description",
		DescriptionLong:  "Unpublished long description",
		ReleaseDate:      &releaseDate,
		CreatedBy:        testUserID,
		LastUpdatedBy:    testUserID,
	}
	_, err = releaseNotesService.Create(unpublishedRN, nil)
	require.NoError(t, err)

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/api/widget/"+testOrg.ExternalID.String()+"/release-notes", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orgId", testOrg.ExternalID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response body
	var response serveReleaseNotesWidgetResponseBody
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// Should only include published release note
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "Published Release Note", response.Data[0].Title)
}

func TestHandleReleaseNotesServe_EmptyReleaseDate(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	err := testDB.DB.Client.AutoMigrate(&organisation.Organisation{}, &releasenotes.ReleaseNote{})
	require.NoError(t, err)

	// Create test organization
	testOrg, err := organisation.New("Test Org")
	require.NoError(t, err)
	err = testDB.DB.Client.Create(testOrg).Error
	require.NoError(t, err)

	// Create release note without release date
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(testDB.DB, deps.ObjStore))
	testUserID := uuid.New()

	rn := &releasenotes.ReleaseNote{
		OrganisationID:   testOrg.ID,
		Title:            "Test Release Note",
		DescriptionShort: "Description",
		DescriptionLong:  "Long description",
		ReleaseDate:      nil, // No release date
		CreatedBy:        testUserID,
		LastUpdatedBy:    testUserID,
	}
	rnID, err := releaseNotesService.Create(rn, nil)
	require.NoError(t, err)

	// Publish the release note
	err = releaseNotesService.ChangePublishedStatus(rnID, true)
	require.NoError(t, err)

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/api/widget/"+testOrg.ExternalID.String()+"/release-notes", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orgId", testOrg.ExternalID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response body
	var response serveReleaseNotesWidgetResponseBody
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// Verify empty date is handled correctly
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "", response.Data[0].Date)
}

func TestHandleReleaseNotesServe_DefaultPagination(t *testing.T) {
	cleanup := testutil.SetupTest()
	defer cleanup()

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	err := testDB.DB.Client.AutoMigrate(&organisation.Organisation{}, &releasenotes.ReleaseNote{})
	require.NoError(t, err)

	// Create test organization
	testOrg, err := organisation.New("Test Org")
	require.NoError(t, err)
	err = testDB.DB.Client.Create(testOrg).Error
	require.NoError(t, err)

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	// Create test request without page/pageSize params
	req := httptest.NewRequest(http.MethodGet, "/api/widget/"+testOrg.ExternalID.String()+"/release-notes", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orgId", testOrg.ExternalID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute handler
	handlers.HandleReleaseNotesServe(rr, req)

	// Assert response - should use defaults (page=1, pageSize=10)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response body
	var response serveReleaseNotesWidgetResponseBody
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// Should return empty array with default pagination
	assert.NotNil(t, response.Data)
}

// Benchmark tests to ensure performance
func BenchmarkHandleReleaseNotesServe(b *testing.B) {
	testutil.DisableLogging()

	// Setup test database
	testDB := testutil.SetupTestDB(&testing.T{})
	defer testDB.Cleanup(&testing.T{})

	// Create test dependencies
	deps := testutil.NewMockDependencies(testDB.DB)

	// Auto-migrate tables
	testDB.DB.Client.AutoMigrate(&organisation.Organisation{}, &releasenotes.ReleaseNote{})

	// Create test organization
	testOrg, _ := organisation.New("Bench Org")
	testDB.DB.Client.Create(testOrg)

	// Create some release notes
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(testDB.DB, deps.ObjStore))
	testUserID := uuid.New()
	releaseDate := "2024-01-15"

	for i := 0; i < 10; i++ {
		rn := &releasenotes.ReleaseNote{
			OrganisationID:   testOrg.ID,
			Title:            "Test Release Note",
			DescriptionShort: "Description",
			DescriptionLong:  "Long description",
			ReleaseDate:      &releaseDate,
			CreatedBy:        testUserID,
			LastUpdatedBy:    testUserID,
		}
		rnID, _ := releaseNotesService.Create(rn, nil)
		releaseNotesService.ChangePublishedStatus(rnID, true)
	}

	// Create handler
	handlers := New(deps.ToSharedDependencies())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/widget/"+testOrg.ExternalID.String()+"/release-notes", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("orgId", testOrg.ExternalID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()
		handlers.HandleReleaseNotesServe(rr, req)
	}
}

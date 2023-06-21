package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"gitea.slauson.io/blog/blog-ms/constants"
	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/blog/blog-ms/handler"
	"gitea.slauson.io/blog/blog-ms/testing/mockdata"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	createdPosts    []*dto.PostResponse
	createdComments []*dto.CommentResponse
)

func runTestServer() *httptest.Server {
	return httptest.NewServer(handler.CreateRouter())
}

func TestCreatePost(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	tests := []struct {
		name string
		req  *dto.CreatePostRequest
		res  *dto.PostResponse
	}{
		{
			name: "success",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				Body:        "Test Body",
				CreatedByID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := sioUtils.NewTokenClient().CreateToken()
			if err != nil {
				t.Fatal(err)
			}

			rJSON, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest("POST", ts.URL+"/api/blog/v1/task", sr)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()
			pr := parsePostResponse(t, resp)
			createdPosts = append(createdPosts, pr)
			require.Equal(t, tt.res.Title, pr.Title)
			require.Equal(t, tt.res.Body, pr.Body)
			require.Equal(t, tt.res.CreatedByID, pr.CreatedByID)
		})
	}
}

func TestCreatePost_Err(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	tests := []struct {
		name   string
		req    *dto.CreatePostRequest
		res    *dto.PostResponse
		status int
		err    string
	}{
		{
			name: "Already Exists",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				Body:        "Test Body",
				CreatedByID: 1,
			},
			status: http.StatusOK,
			err:    "Post already exists",
		},
		{
			name: "Bad Title",
			req: &dto.CreatePostRequest{
				Title:       mockdata.LongTitle,
				Body:        "Test Body",
				CreatedByID: 1,
			},
			status: http.StatusBadRequest,
			err:    "Post already exists",
		},
		{
			name: "No Title",
			req: &dto.CreatePostRequest{
				Body:        "Test Body",
				CreatedByID: 1,
			},
			status: http.StatusBadRequest,
			err:    "Post already exists",
		},
		{
			name: "Bad Body",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				Body:        mockdata.LongBody,
				CreatedByID: 1,
			},
			status: http.StatusBadRequest,
			err:    "Post already exists",
		},
		{
			name: "No Body",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				CreatedByID: 1,
			},
			status: http.StatusBadRequest,
			err:    "Post already exists",
		},
		{
			name: "Missing CreatedByID",
			req: &dto.CreatePostRequest{
				Title: "Title",
				Body:  "Test Body",
			},
			status: http.StatusBadRequest,
			err:    "Post already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := sioUtils.NewTokenClient().CreateToken()
			if err != nil {
				t.Fatal(err)
			}

			defer ts.Close()

			rJSON, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest("POST", ts.URL+"/api/blog/v1/task", sr)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()
			checkIfCorrectError(t, resp, tt.err, tt.status)
		})
	}
}

func TestUpdatePost(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	tests := []struct {
		name string
		req  *dto.UpdatePostRequest
		res  *dto.PostResponse
	}{
		{
			name: "success Update Both",
			req: &dto.UpdatePostRequest{
				Title:       "Title",
				Body:        "Test Body",
				UpdatedByID: 1,
			},
		},
		{
			name: "success Update Title",
			req: &dto.UpdatePostRequest{
				Title:       "Title",
				UpdatedByID: 1,
			},
		},
		{
			name: "success Update Body",
			req: &dto.UpdatePostRequest{
				Body:        "Test Body",
				UpdatedByID: 1,
			},
		},
	}
	id := createdPosts[0].ID

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := sioUtils.NewTokenClient().CreateToken()
			if err != nil {
				t.Fatal(err)
			}

			rJSON, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest(
				"PATCH",
				ts.URL+"/api/blog/v1/task/"+strconv.Itoa(int(id)),
				sr,
			)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()
			pr := parsePostResponse(t, resp)
			require.Equal(t, tt.res.Title, pr.Title)
			require.Equal(t, tt.res.Body, pr.Body)
			require.Equal(t, tt.req.UpdatedByID, pr.UpdatedByID)
		})
	}
}

func TestUpdatePost_Err(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	tests := []struct {
		name   string
		req    *dto.UpdatePostRequest
		status int
		err    string
	}{
		{
			name: "Bad ID",
			req: &dto.UpdatePostRequest{
				Body:        "Test Body",
				UpdatedByID: 1,
			},
			status: http.StatusBadRequest,
			err:    "1sdfsdfe",
		},
		{
			name: "Bad Title",
			req: &dto.UpdatePostRequest{
				Title:       mockdata.LongTitle,
				Body:        "Test Body",
				UpdatedByID: 1,
			},
			status: http.StatusBadRequest,
			err:    "1",
		},
		{
			name: "Bad Body",
			req: &dto.UpdatePostRequest{
				Title:       "Title",
				Body:        mockdata.LongBody,
				UpdatedByID: 1,
			},
			status: http.StatusBadRequest,
			err:    "1",
		},
		{
			name: "Bad - No updates passed",
			req: &dto.UpdatePostRequest{
				UpdatedByID: 1,
			},
			status: http.StatusBadRequest,
			err:    "1",
		},
		{
			name: "Missing UpdatedByID",
			req: &dto.UpdatePostRequest{
				Title: "Title",
				Body:  "Test Body",
			},
			status: http.StatusBadRequest,
			err:    "1",
		},
	}
	id := createdPosts[0].ID

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := sioUtils.NewTokenClient().CreateToken()
			if err != nil {
				t.Fatal(err)
			}

			rJSON, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest(
				"PATCH",
				ts.URL+"/api/blog/v1/task/"+strconv.Itoa(int(id)),
				sr,
			)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()
			checkIfCorrectError(t, resp, tt.err, tt.status)
		})
	}
}

func TestAddComment(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	tests := []struct {
		name string
		req  *dto.AddCommentRequest
	}{
		{
			name: "success",
			req: &dto.AddCommentRequest{
				Content: "Test Content",
				UserID:  1,
				PostID:  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := sioUtils.NewTokenClient().CreateToken()
			if err != nil {
				t.Fatal(err)
			}

			defer ts.Close()

			rJSON, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest("POST", ts.URL+"/api/blog/v1/task/comment", sr)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			cr := parseCommentResponse(t, resp)
			createdComments = append(createdComments, cr)
			require.Equal(t, tt.req.Content, cr.Content)
			require.Equal(t, tt.req.UserID, cr.UserID)
			require.Equal(t, tt.req.PostID, cr.PostID)
		})
	}
}

func TestAddComment_Err(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	tests := []struct {
		name   string
		req    *dto.AddCommentRequest
		status int
		err    string
	}{
		{
			name: "Bad Content",
			req: &dto.AddCommentRequest{
				Content: mockdata.LongComment,
				UserID:  1,
				PostID:  1,
			},
			status: http.StatusBadRequest,
			err:    constants.COMMENT_TOO_LONG,
		},
		{
			name: "Missing Content",
			req: &dto.AddCommentRequest{
				UserID: 1,
				PostID: 1,
			},
			status: http.StatusBadRequest,
			err:    constants.COMMENT_TOO_LONG,
		},
		{
			name: "Missing UserID",
			req: &dto.AddCommentRequest{
				Content: "Test Content",
				PostID:  1,
			},
			status: http.StatusBadRequest,
			err:    constants.COMMENT_TOO_LONG,
		},
		{
			name: "Missing PostID",
			req: &dto.AddCommentRequest{
				Content: "Test Content",
				UserID:  1,
			},
			status: http.StatusBadRequest,
			err:    constants.COMMENT_TOO_LONG,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := sioUtils.NewTokenClient().CreateToken()
			if err != nil {
				t.Fatal(err)
			}

			defer ts.Close()

			rJSON, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest("POST", ts.URL+"/api/blog/v1/task/comment", sr)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()
			checkIfCorrectError(t, resp, tt.err, tt.status)
		})
	}
}

func TestUpdateComment(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	id := createdComments[0].ID

	tests := []struct {
		name string
		req  *dto.UpdateCommentRequest
		res  *dto.CommentResponse
	}{
		{
			name: "success",
			req: &dto.UpdateCommentRequest{
				Content: "Test Content",
			},
		},
	}

	errTests := []struct {
		name   string
		req    *dto.UpdateCommentRequest
		err    string
		status int
	}{
		{
			name: "Bad ID",
			req: &dto.UpdateCommentRequest{
				Content: "Test Content",
			},
			status: http.StatusBadRequest,
			err:    "1asdf",
		},
		{
			name: "Bad Content",
			req: &dto.UpdateCommentRequest{
				Content: mockdata.LongComment,
			},
			status: http.StatusBadRequest,
			err:    constants.COMMENT_TOO_LONG,
		},
		{
			name:   "No Content",
			req:    &dto.UpdateCommentRequest{},
			status: http.StatusBadRequest,
			err:    "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := sioUtils.NewTokenClient().CreateToken()
			if err != nil {
				t.Fatal(err)
			}

			rJSON, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest(
				"PATCH",
				ts.URL+"/api/blog/v1/task/comment/"+strconv.Itoa(int(id)),
				sr,
			)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			cr := parseCommentResponse(t, resp)
			require.Equal(t, tt.req.Content, cr.Content)
		})
	}

	for _, et := range errTests {
		t.Run(et.name, func(t *testing.T) {
			t.Parallel()
			token, err := sioUtils.NewTokenClient().CreateToken()
			if err != nil {
				t.Fatal(err)
			}

			rJSON, err := json.Marshal(et.req)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest(
				"PATCH",
				ts.URL+"/api/blog/v1/task/comment/"+strconv.Itoa(int(id)),
				sr,
			)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()
			checkIfCorrectError(t, resp, et.err, et.status)
		})
	}
}

func TestUpdateComment_Err(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	errTests := []struct {
		name   string
		ID     string
		req    *dto.UpdateCommentRequest
		res    *dto.CommentResponse
		status int
	}{
		{
			name: "success",
			req: &dto.UpdateCommentRequest{
				Content: "Test Content",
			},
			status: http.StatusOK,
			ID:     "1",
		},
		{
			name: "success Bad ID",
			req: &dto.UpdateCommentRequest{
				Content: "Test Content",
			},
			status: http.StatusBadRequest,
			ID:     "1asdf",
		},
		{
			name: "success Bad Content",
			req: &dto.UpdateCommentRequest{
				Content: mockdata.LongComment,
			},
			status: http.StatusBadRequest,
			ID:     "1",
		},
		{
			name:   "success No Content",
			req:    &dto.UpdateCommentRequest{},
			status: http.StatusBadRequest,
			ID:     "1",
		},
	}
	id := createdComments[0].ID

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := sioUtils.NewTokenClient().CreateToken()
			if err != nil {
				t.Fatal(err)
			}

			rJSON, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest(
				"PATCH",
				ts.URL+"/api/blog/v1/task/comment/"+strconv.Itoa(int(id)),
				sr,
			)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			cr := parseCommentResponse(t, resp)
			require.Equal(t, tt.req.Content, cr.Content)
		})
	}
}

func TestGetPost(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	tests := []struct {
		name    string
		request string
		status  int
		result  *dto.PostResponse
	}{
		{name: "happy", request: "123", status: http.StatusOK, result: mockdata.PostResponse},
		{name: "BadId", request: "1s23", status: http.StatusBadRequest, result: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := runTestServer()
			defer ts.Close()
			hdlr, mSvc := initEnv()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.request}}
			if tt.status == http.StatusOK {
				i, _ := strconv.ParseInt(tt.request, 10, 64)
				mSvc.On("GetPost", i).
					Return(tt.result, nil)
			}

			hdlr.GetPost(c)
			if tt.status == http.StatusOK {
				assert.Nilf(t, c.Errors, "c.Errors should be nil")
			} else {
				assert.NotNilf(t, c.Errors, "c.Errors shouldnt be nil")
			}
		})
	}
}

func TestGetAllPosts(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	hdlr, mSvc := initEnv()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	mSvc.On("GetAllPosts").
		Return(mockdata.PostResponses, nil)

	hdlr.GetAllPosts(c)
	assert.Nilf(t, c.Errors, "c.Errors should be nil")
}

func TestSoftDeleteComment(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	tests := []struct {
		name    string
		request string
		status  int
		result  *siogeneric.SuccessResponse
	}{
		{
			name:    "happy",
			request: "123",
			status:  http.StatusOK,
			result:  mockdata.SuccessResponseSuccess,
		},
		{name: "BadId", request: "1s23", status: http.StatusBadRequest, result: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := runTestServer()
			defer ts.Close()
			hdlr, mSvc := initEnv()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.request}}
			if tt.status == http.StatusOK {
				i, _ := strconv.ParseInt(tt.request, 10, 64)
				mSvc.On("SoftDeleteComment", i).
					Return(tt.result, nil)
			}

			hdlr.SoftDeleteComment(c)
			if tt.status == http.StatusOK {
				assert.Nilf(t, c.Errors, "c.Errors should be nil")
			} else {
				assert.NotNilf(t, c.Errors, "c.Errors shouldnt be nil")
			}
		})
	}
}

func TestSoftDeletePost(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	tests := []struct {
		name    string
		request string
		status  int
		result  *siogeneric.SuccessResponse
	}{
		{
			name:    "happy",
			request: "123",
			status:  http.StatusOK,
			result:  mockdata.SuccessResponseSuccess,
		},
		{name: "BadId", request: "1s23", status: http.StatusBadRequest, result: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := runTestServer()
			defer ts.Close()
			hdlr, mSvc := initEnv()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.request}}
			if tt.status == http.StatusOK {
				i, _ := strconv.ParseInt(tt.request, 10, 64)
				mSvc.On("SoftDeletePost", i).
					Return(tt.result, nil)
			}

			hdlr.SoftDeletePost(c)
			if tt.status == http.StatusOK {
				assert.Nilf(t, c.Errors, "c.Errors should be nil")
			} else {
				assert.NotNilf(t, c.Errors, "c.Errors shouldnt be nil")
			}
		})
	}
}

func parsePostResponse(t *testing.T, resp *http.Response) *dto.PostResponse {
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	postResponse := new(dto.PostResponse)
	err := json.NewDecoder(resp.Body).Decode(postResponse)
	if err != nil {
		t.Fatal(err)
	}

	return postResponse
}

func parseCommentResponse(t *testing.T, resp *http.Response) *dto.CommentResponse {
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	commentResponse := new(dto.CommentResponse)
	err := json.NewDecoder(resp.Body).Decode(commentResponse)
	if err != nil {
		t.Fatal(err)
	}

	return commentResponse
}

func parseExistsResponse(t *testing.T, resp *http.Response) *siogeneric.ExistsResponse {
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	existsResponse := new(siogeneric.ExistsResponse)
	err := json.NewDecoder(resp.Body).Decode(existsResponse)
	if err != nil {
		t.Fatal(err)
	}

	return existsResponse
}

func checkIfCorrectError(
	t *testing.T,
	resp *http.Response,
	expectedError string,
	expectedStatus int,
) {
	if resp.StatusCode != expectedStatus {
		t.Fatalf("Expected status code %d, got %d", expectedStatus, resp.StatusCode)
	}
	result := new(siogeneric.ErrorResponse)
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		t.Fatal(err)
	}

	if expectedError != result.Error {
		t.Fatalf("Expected result %s, got %s", expectedError, result.Error)
	}
}

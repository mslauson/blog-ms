package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"gitea.slauson.io/blog/blog-ms/constants"
	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/blog/blog-ms/handler"
	"gitea.slauson.io/blog/blog-ms/integration/mockdata"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
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

func TestGetPost(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	token, err := sioUtils.NewTokenClient().CreateToken()
	if err != nil {
		t.Fatal(err)
	}

	id := createdPosts[0].ID
	errTests := []struct {
		name    string
		request string
		status  int
		err     string
	}{
		{
			name:    "Not Found",
			request: "123",
			status:  http.StatusNotFound,
			err:     "asdf",
		},

		{
			name:    "BadId",
			request: "1s23",
			status:  http.StatusBadRequest,
			err:     constants.INVALID_ID,
		},
	}

	t.Run("Happy", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("GET", ts.URL+"/api/blog/v1/post/"+strconv.Itoa(int(id)), nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
			}
		}(resp.Body)

		result := parsePostResponse(t, resp)
		require.Equal(t, id, result)
	})

	for _, et := range errTests {
		t.Run(et.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(
				"GET",
				ts.URL+"/api/blog/v1/post/"+strconv.Itoa(int(id)),
				nil,
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

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
				}
			}(resp.Body)

			checkIfCorrectError(t, resp, et.err, et.status)
		})
	}
}

func TestGetAllPosts(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	token, err := sioUtils.NewTokenClient().CreateToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Parallel()

	req, err := http.NewRequest("GET", ts.URL+"/api/blog/v1/post", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(resp.Body)

	_ = parsePostResponse(t, resp)
}

func TestSoftDeleteComment(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	token, err := sioUtils.NewTokenClient().CreateToken()
	if err != nil {
		t.Fatal(err)
	}

	id := createdComments[0].ID
	idStr := strconv.Itoa(int(id))
	errTests := []struct {
		name   string
		id     string
		status int
		err    string
	}{
		{
			name:   "Not Found - Already Deleted",
			id:     idStr,
			status: http.StatusNotFound,
			err:    "asdf",
		},
		{
			name:   "Not Found",
			id:     "123",
			status: http.StatusNotFound,
			err:    "asdf",
		},

		{
			name:   "BadId",
			id:     "1s23",
			status: http.StatusBadRequest,
			err:    constants.INVALID_ID,
		},
	}

	t.Run("Happy", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest(
			"DELETE",
			ts.URL+"/api/blog/v1/post/comment/"+idStr,
			nil,
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

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
			}
		}(resp.Body)

		result := parsePostResponse(t, resp)
		require.Equal(t, id, result)
	})

	for _, et := range errTests {
		t.Run(et.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(
				"DELETE",
				ts.URL+"/api/blog/v1/post/comment/"+idStr,
				nil,
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

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
				}
			}(resp.Body)

			checkIfCorrectError(t, resp, et.err, et.status)
		})
	}
}

func TestSoftDeletePost(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()
	token, err := sioUtils.NewTokenClient().CreateToken()
	if err != nil {
		t.Fatal(err)
	}

	id := createdPosts[0].ID
	idStr := strconv.Itoa(int(id))
	errTests := []struct {
		name   string
		id     string
		status int
		err    string
	}{
		{
			name:   "Not Found - Already Deleted",
			id:     idStr,
			status: http.StatusNotFound,
			err:    "asdf",
		},
		{
			name:   "Not Found",
			id:     "123",
			status: http.StatusNotFound,
			err:    "asdf",
		},

		{
			name:   "BadId",
			id:     "1s23",
			status: http.StatusBadRequest,
			err:    constants.INVALID_ID,
		},
	}

	t.Run("Happy", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest(
			"DELETE",
			ts.URL+"/api/blog/v1/post/comment/"+idStr,
			nil,
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

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
			}
		}(resp.Body)

		result := parsePostResponse(t, resp)
		require.Equal(t, id, result)
	})

	for _, et := range errTests {
		t.Run(et.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(
				"DELETE",
				ts.URL+"/api/blog/v1/post/comment/"+idStr,
				nil,
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

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
				}
			}(resp.Body)

			checkIfCorrectError(t, resp, et.err, et.status)
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

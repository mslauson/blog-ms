package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"gitea.slauson.io/blog/post-ms/constants"
	"gitea.slauson.io/blog/post-ms/dto"
	"gitea.slauson.io/blog/post-ms/handler"
	"gitea.slauson.io/blog/post-ms/integration/mockdata"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"github.com/stretchr/testify/require"
)

var (
	userID          = int64(10000011)
	createdPosts    []*dto.PostResponse
	createdComments []*dto.CommentResponse
)

func runTestServer() *httptest.Server {
	return httptest.NewServer(handler.CreateRouter())
}

func TestCreatePost(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	cr := &dto.CreatePostRequest{
		Title:       "Title",
		Body:        "Test Body",
		CreatedByID: userID,
	}

	token, err := sioUtils.NewTokenClient().CreateToken()
	if err != nil {
		t.Fatal(err)
	}

	rJSON, err := json.Marshal(cr)
	if err != nil {
		t.Fatal(err)
	}
	sr := strings.NewReader(string(rJSON))
	req, err := http.NewRequest("POST", ts.URL+"/api/post/v1", sr)
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
	require.Equal(t, cr.Title, pr.Title)
	require.Equal(t, cr.Body, pr.Body)
	require.Equal(t, cr.CreatedByID, pr.CreatedByID)
}

func TestCreatePost_Err(t *testing.T) {
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
				CreatedByID: userID,
			},
			status: http.StatusBadRequest,
			err:    "post already exists",
		},
		{
			name: "Bad Title",
			req: &dto.CreatePostRequest{
				Title:       mockdata.LongTitle,
				Body:        "Test Body",
				CreatedByID: userID,
			},
			status: http.StatusBadRequest,
			err:    constants.TITLE_TOO_LONG,
		},
		{
			name: "No Title",
			req: &dto.CreatePostRequest{
				Body:        "Test Body",
				CreatedByID: userID,
			},
			status: http.StatusBadRequest,
			err:    "Key: 'CreatePostRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag",
		},
		{
			name: "Bad Body",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				Body:        mockdata.LongBody,
				CreatedByID: userID,
			},
			status: http.StatusBadRequest,
			err:    constants.BODY_TOO_LONG,
		},
		{
			name: "No Body",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				CreatedByID: userID,
			},
			status: http.StatusBadRequest,
			err:    "Key: 'CreatePostRequest.Body' Error:Field validation for 'Body' failed on the 'required' tag",
		},
		{
			name: "Missing CreatedByID",
			req: &dto.CreatePostRequest{
				Title: "Title",
				Body:  "Test Body",
			},
			status: http.StatusBadRequest,
			err:    "Key: 'CreatePostRequest.CreatedByID' Error:Field validation for 'CreatedByID' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := runTestServer()
			defer ts.Close()
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
			req, err := http.NewRequest("POST", ts.URL+"/api/post/v1", sr)
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
				UpdatedByID: userID,
			},
		},
		{
			name: "success Update Title",
			req: &dto.UpdatePostRequest{
				Title:       "Title",
				UpdatedByID: userID,
			},
		},
		{
			name: "success Update Body",
			req: &dto.UpdatePostRequest{
				Body:        "Test Body",
				UpdatedByID: userID,
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
				ts.URL+"/api/post/v1/"+strconv.Itoa(int(id)),
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
			require.Equal(t, tt.req.Title, pr.Title)
			require.Equal(t, tt.req.Body, pr.Body)
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
			name: "Bad Title",
			req: &dto.UpdatePostRequest{
				Title:       mockdata.LongTitle,
				Body:        "Test Body",
				UpdatedByID: userID,
			},
			status: http.StatusBadRequest,
			err:    constants.TITLE_TOO_LONG,
		},
		{
			name: "Bad Body",
			req: &dto.UpdatePostRequest{
				Title:       "Title",
				Body:        mockdata.LongBody,
				UpdatedByID: userID,
			},
			status: http.StatusBadRequest,
			err:    constants.BODY_TOO_LONG,
		},
		{
			name: "Bad - No updates passed",
			req: &dto.UpdatePostRequest{
				UpdatedByID: userID,
			},
			status: http.StatusBadRequest,
			err:    constants.POST_UPDATE_INVALID,
		},
		{
			name: "Missing UpdatedByID",
			req: &dto.UpdatePostRequest{
				Title: "Title",
				Body:  "Test Body",
			},
			status: http.StatusBadRequest,
			err:    "Key: 'UpdatePostRequest.UpdatedByID' Error:Field validation for 'UpdatedByID' failed on the 'required' tag",
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
				ts.URL+"/api/post/v1/"+strconv.Itoa(int(id)),
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

func TestUpdatePost_ErrIDNotNumerical(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	ur := &dto.UpdatePostRequest{
		Title:       "asdf",
		UpdatedByID: userID,
	}
	token, err := sioUtils.NewTokenClient().CreateToken()
	if err != nil {
		t.Fatal(err)
	}

	rJSON, err := json.Marshal(ur)
	if err != nil {
		t.Fatal(err)
	}
	sr := strings.NewReader(string(rJSON))
	req, err := http.NewRequest(
		"PATCH",
		ts.URL+"/api/post/v1/asdf8",
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
	checkIfCorrectError(t, resp, "asdfj", http.StatusBadRequest)
}

func TestAddComment(t *testing.T) {
	postID := createdPosts[0].ID

	ts := runTestServer()
	defer ts.Close()
	acr := &dto.AddCommentRequest{
		Content: "Test Content",
		UserID:  userID,
		PostID:  postID,
	}

	token, err := sioUtils.NewTokenClient().CreateToken()
	if err != nil {
		t.Fatal(err)
	}

	defer ts.Close()

	rJSON, err := json.Marshal(acr)
	if err != nil {
		t.Fatal(err)
	}
	sr := strings.NewReader(string(rJSON))
	req, err := http.NewRequest("POST", ts.URL+"/api/post/v1/comment", sr)
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
	require.Equal(t, acr.Content, cr.Content)
	require.Equal(t, acr.UserID, cr.UserID)
	require.Equal(t, acr.PostID, cr.PostID)
}

func TestAddComment_Err(t *testing.T) {
	postID := createdPosts[0].ID
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
				UserID:  userID,
				PostID:  postID,
			},
			status: http.StatusBadRequest,
			err:    constants.COMMENT_TOO_LONG,
		},
		{
			name: "Missing Content",
			req: &dto.AddCommentRequest{
				UserID: userID,
				PostID: postID,
			},
			status: http.StatusBadRequest,
			err:    "Key: 'AddCommentRequest.Content' Error:Field validation for 'Content' failed on the 'required' tag",
		},
		{
			name: "Missing UserID",
			req: &dto.AddCommentRequest{
				Content: "Test Content",
				PostID:  userID,
			},
			status: http.StatusBadRequest,
			err:    "Key: 'AddCommentRequest.UserID' Error:Field validation for 'UserID' failed on the 'required' tag",
		},
		{
			name: "Missing PostID",
			req: &dto.AddCommentRequest{
				Content: "Test Content",
				UserID:  userID,
			},
			status: http.StatusBadRequest,
			err:    "Key: 'AddCommentRequest.PostID' Error:Field validation for 'PostID' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := sioUtils.NewTokenClient().CreateToken()
			if err != nil {
				t.Fatal(err)
			}

			ts := runTestServer()
			defer ts.Close()

			rJSON, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			sr := strings.NewReader(string(rJSON))
			req, err := http.NewRequest("POST", ts.URL+"/api/post/v1/comment", sr)
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

	ucr := &dto.UpdateCommentRequest{
		Content: "Test Content",
	}

	errTests := []struct {
		name   string
		req    *dto.UpdateCommentRequest
		err    string
		status int
	}{
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
			err:    "Key: 'UpdateCommentRequest.Content' Error:Field validation for 'Content' failed on the 'required' tag",
		},
	}

	t.Run("Happy", func(t *testing.T) {
		token, err := sioUtils.NewTokenClient().CreateToken()
		if err != nil {
			t.Fatal(err)
		}

		rJSON, err := json.Marshal(ucr)
		if err != nil {
			t.Fatal(err)
		}

		sr := strings.NewReader(string(rJSON))
		req, err := http.NewRequest(
			"PATCH",
			ts.URL+"/api/post/v1/comment/"+strconv.Itoa(int(id)),
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
		require.Equal(t, ucr.Content, cr.Content)
	})

	for _, et := range errTests {
		t.Run(et.name, func(t *testing.T) {
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
				ts.URL+"/api/post/v1/comment/"+strconv.Itoa(int(id)),
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

	t.Run("badID", func(t *testing.T) {
		token, err := sioUtils.NewTokenClient().CreateToken()
		if err != nil {
			t.Fatal(err)
		}

		rJSON, err := json.Marshal(ucr)
		if err != nil {
			t.Fatal(err)
		}
		sr := strings.NewReader(string(rJSON))
		req, err := http.NewRequest(
			"PATCH",
			ts.URL+"/api/post/v1/comment/"+strconv.Itoa(int(id)),
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
		require.Equal(t, ucr.Content, cr.Content)
	})
}

func TestGetPost(t *testing.T) {
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
			err:     "post not found",
		},

		{
			name:    "BadId",
			request: "1s23",
			status:  http.StatusBadRequest,
			err:     constants.INVALID_ID,
		},
	}

	t.Run("Happy", func(t *testing.T) {
		ts := runTestServer()
		defer ts.Close()
		req, err := http.NewRequest("GET", ts.URL+"/api/post/v1/"+strconv.Itoa(int(id)), nil)
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

		result := parsePostResponse(t, resp)
		require.Equal(t, id, result.ID)
	})

	for _, et := range errTests {
		t.Run(et.name, func(t *testing.T) {
			ts := runTestServer()
			defer ts.Close()
			req, err := http.NewRequest(
				"GET",
				ts.URL+"/api/post/v1/"+et.request,
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
					t.Fatal(err)
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

	req, err := http.NewRequest("GET", ts.URL+"/api/post/v1", nil)
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

	_ = parsePostResponses(t, resp)
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
		req, err := http.NewRequest(
			"DELETE",
			ts.URL+"/api/post/v1/comment/"+idStr,
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
				t.Fatal(err)
			}
		}(resp.Body)

		result := parseSuccessResponse(t, resp)
		require.Equal(t, true, result.Success)
	})

	for _, et := range errTests {
		t.Run(et.name, func(t *testing.T) {
			req, err := http.NewRequest(
				"DELETE",
				ts.URL+"/api/post/v1/comment/"+idStr,
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
					t.Fatal(err)
				}
			}(resp.Body)

			checkIfCorrectError(t, resp, et.err, et.status)
		})
	}
}

func TestSoftDeletePost(t *testing.T) {
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
		ts := runTestServer()
		defer ts.Close()

		req, err := http.NewRequest(
			"DELETE",
			ts.URL+"/api/post/v1/"+idStr,
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
				t.Fatal(err)
			}
		}(resp.Body)

		result := parseSuccessResponse(t, resp)
		require.Equal(t, true, result.Success)
	})

	for _, et := range errTests {
		t.Run(et.name, func(t *testing.T) {
			ts := runTestServer()
			defer ts.Close()
			req, err := http.NewRequest(
				"DELETE",
				ts.URL+"/api/post/v1/"+idStr,
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
					t.Fatal(err)
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

func parsePostResponses(t *testing.T, resp *http.Response) []*dto.PostResponse {
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	postResponses := make([]*dto.PostResponse, 0)
	err := json.NewDecoder(resp.Body).Decode(&postResponses)
	if err != nil {
		t.Fatal(err)
	}

	return postResponses
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

func parseSuccessResponse(t *testing.T, resp *http.Response) *siogeneric.SuccessResponse {
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	successResponse := new(siogeneric.SuccessResponse)
	err := json.NewDecoder(resp.Body).Decode(successResponse)
	if err != nil {
		t.Fatal(err)
	}

	return successResponse
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

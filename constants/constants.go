package constants

const (
	POST             = "post"
	POSTS            = "posts"
	COMMENT          = "comment"
	POST_EXISTS      = "post already exists"
	NO_POST_FOUND    = "no post found"
	NO_COMMENT_FOUND = "no comment found"

	TITLE_TOO_LONG      = "title is too long.  please choose another title with no more than 100 characters."
	BODY_TOO_LONG       = "post body is too long. please choose another post with no more than 10000 characters."
	COMMENT_TOO_LONG    = "comment is too long. please choose another comment with no more than 5000 characters."
	POST_UPDATE_INVALID = "title, body, or both must be provided to update a post"
)

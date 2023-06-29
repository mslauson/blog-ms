package constants

const (
	SELECT_ITEMS_POST    = `id, title, body, created_by_id, updated_by_id, posted_date, updated_date, deletion_date, soft_deleted`
	SELECT_ITEMS_COMMENT = `id, content, comment_date, user_id, post_id, soft_deleted, deletion_date`
)

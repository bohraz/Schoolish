package database

import (
	"root/internal/model"
)

func CreatePost(post model.Post) (int, error) {
	result, err := DB.Exec("INSERT INTO app.posts (title, content, userId) VALUES (?, ?, ?)", post.Title, post.Content, post.UserId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetPosts(amount int) ([]model.Post, error) {
	rows, err := DB.Query("SELECT * FROM app.posts LIMIT ?", amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.AnswerId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func GetPost(id int) (model.Post, error) {
	query := `
		SELECT app.posts.*, COUNT(app.comments.postId) 
		FROM app.posts
		LEFT JOIN app.comments ON app.posts.postId = app.comments.postId
		WHERE app.posts.postId = ?
		GROUP BY app.posts.postID
	`

	var post model.Post
	err := DB.QueryRow(query, id).Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.AnswerId, &post.Comments)
	if err != nil {
		return post, err
	}

	return post, nil
}

func GetCommentsByPostId(postId, limit, offset int) ([]model.Comment, error) {
	query := `
		SELECT c.*, u.firstName, u.lastName, u.username 
		FROM app.comments c
		INNER JOIN app.users u ON c.userId = u.userId
		WHERE c.postId = ?
		ORDER BY c.postId ASC
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, postId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]model.Comment, 0)
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.Id, &comment.Content, &comment.PostId, &comment.User.Id, &comment.ReplyId, &comment.User.FirstName, &comment.User.LastName, &comment.User.Handle)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func CreateComment(comment model.Comment) (int, error) {
	result, err := DB.Exec("INSERT INTO app.comments (userId, postId, content) VALUES (?, ?, ?)", comment.User.Id, comment.PostId, comment.Content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	post, err := GetPost(comment.PostId)
	if err != nil {
		return 0, err
	}

	model.PostsBroadcast <- post

	return int(id), nil
}
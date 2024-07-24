package database

import "root/internal/model"

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
	var post model.Post
	err := DB.QueryRow("SELECT * FROM app.posts WHERE postId = ?", id).Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.AnswerId)
	if err != nil {
		return post, err
	}

	return post, nil
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

	return int(id), nil
}
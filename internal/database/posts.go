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
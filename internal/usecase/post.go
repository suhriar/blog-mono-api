package usecase

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/suhriar/blog-mono-api/model"
)

func (u *postUsecase) CreatePost(ctx context.Context, userID int64, req model.CreatePostRequest) (err error) {
	postHashtags := strings.Join(req.PostHashtags, ",")

	now := time.Now()
	model := model.Post{
		UserID:       userID,
		PostTitle:    req.PostTitle,
		PostContent:  req.PostContent,
		PostHashtags: postHashtags,
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    strconv.FormatInt(userID, 10),
		UpdatedBy:    strconv.FormatInt(userID, 10),
	}

	_, err = u.postRepository.CreatePost(ctx, model)
	if err != nil {
		return err
	}
	return nil
}

func (u *postUsecase) GetPostByID(ctx context.Context, postID int64) (post model.GetPostResponse, err error) {
	postDetail, err := u.postRepository.GetPostByID(ctx, postID)
	if err != nil {
		return
	}

	likeCount, err := u.postRepository.CountLikeByPostID(ctx, postID)
	if err != nil {
		log.Error().Err(err).Msg("error count like to database")
		return
	}

	comments, err := u.postRepository.GetCommentsByPostID(ctx, postID)
	if err != nil {
		log.Error().Err(err).Msg("error get comment to database")
		return
	}

	post = model.GetPostResponse{
		PostDetail: model.PostDetail{
			ID:           postDetail.ID,
			UserID:       postDetail.UserID,
			Username:     postDetail.Username,
			PostTitle:    postDetail.PostTitle,
			PostContent:  postDetail.PostContent,
			PostHashtags: postDetail.PostHashtags,
			IsLiked:      postDetail.IsLiked,
		},
		LikeCount: likeCount,
		Comments:  comments,
	}

	return
}

func (u *postUsecase) GetAllPost(ctx context.Context, pageSize, pageIndex int) (posts model.GetAllPostResponse, err error) {
	limit := pageSize
	offset := pageSize * (pageIndex - 1)
	posts, err = u.postRepository.GetAllPost(ctx, limit, offset)
	if err != nil {
		return
	}
	return
}

func (u *postUsecase) CreateComment(ctx context.Context, postID, userID int64, request model.CreateCommentRequest) (err error) {
	now := time.Now()
	comment := model.Comment{
		PostID:         postID,
		UserID:         userID,
		CommentContent: request.CommentContent,
		CreatedAt:      now,
		UpdatedAt:      now,
		CreatedBy:      strconv.FormatInt(userID, 10),
		UpdatedBy:      strconv.FormatInt(userID, 10),
	}

	_, err = u.postRepository.CreateComment(ctx, comment)
	if err != nil {
		return err
	}

	return nil
}

func (u *postUsecase) UpsertUserActivity(ctx context.Context, postID, userID int64, request model.UserActivityRequest) (err error) {
	now := time.Now()
	userActivityReq := model.UserActivity{
		PostID:    postID,
		UserID:    userID,
		IsLiked:   request.IsLiked,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: strconv.FormatInt(userID, 10),
		UpdatedBy: strconv.FormatInt(userID, 10),
	}
	userActivity, err := u.postRepository.GetUserActivity(ctx, userActivityReq)
	if err != nil {
		log.Error().Err(err).Msg("error get user activity from database")
		return err
	}

	if userActivity.ID == 0 {
		// create user activity
		if !request.IsLiked {
			return errors.New("never liked this post")
		}
		_, err = u.postRepository.CreateUserActivity(ctx, userActivityReq)
	} else {
		// update user activity
		err = u.postRepository.UpdateUserActivity(ctx, userActivityReq)
	}
	if err != nil {
		return err
	}

	return nil
}

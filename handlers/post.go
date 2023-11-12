package handlers

import (
	"encoding/json"
	"net/http"
	"rest-ws/models"
	"rest-ws/repository"
	"rest-ws/server"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Id 	string  `json:"id"`
	PostContent string `json:"post_content"`
}

type PostUpdateResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}



func InsertPostHandler(s server.Server) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		value := r.Context().Value("userToken")
		if value == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, ok := value.(*jwt.Token)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok {
			userID := claims.UserId
			var req UpsertPostRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Check if PostContent is empty
            if req.PostContent == "" {
                http.Error(w, "Post content cannot be empty", http.StatusBadRequest)
                return
            }

			id, err := ksuid.NewRandom()
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			post := models.Post{
				ID: id.String(),
				PostContent: req.PostContent,
				UserId: userID,
			}

			err = repository.InsertPost(r.Context(), &post)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var postMessage = models.WebsocketMessage{
				Type: "post_created",
				Payload: post,
			}

			s.Hub().Broadcast(postMessage, nil)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(PostResponse{
				Message: "Post created successfully",
				Status:  true,
				Id: post.ID,
				PostContent: post.PostContent,
			})
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func GetPostByIdHandler (s server.Server) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		value := r.Context().Value("userToken")
		if value == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, ok := value.(*jwt.Token)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok {
			userID := claims.UserId
			vars := mux.Vars(r)
			postID := vars["id"]

			post, err := repository.GetPostById(r.Context(), postID)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if post.UserId != userID{
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(PostResponse{
				Message: "Post retrieved successfully",
				Status:  true,
				Id: post.ID,
				PostContent: post.PostContent,
			})
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		value := r.Context().Value("userToken")
		if value == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, ok := value.(*jwt.Token)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok {
			userID := claims.UserId
			vars := mux.Vars(r)
			postID := vars["id"]

			post, err := repository.GetPostById(r.Context(), postID)
			if err != nil{
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}

			if post.UserId != userID{
				http.Error(w, "Unauthorized over not user post", http.StatusUnauthorized)
				return
			}

			var req UpsertPostRequest
			err = json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Check if PostContent is empty
			if req.PostContent == "" {
				http.Error(w, "Post content cannot be empty", http.StatusBadRequest)
				return
			}

			post.PostContent = req.PostContent

			err = repository.UpdatePost(r.Context(), post)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Post updated successfully",
				Status:  true,
			})
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func DeletePostHandler(s server.Server) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		value := r.Context().Value("userToken")
		if value == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, ok := value.(*jwt.Token)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok {
			userID := claims.UserId
			vars := mux.Vars(r)
			postID := vars["id"]

			post, err := repository.GetPostById(r.Context(), postID)
			if err != nil{
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}

			if post.UserId != userID{
				http.Error(w, "Unauthorized over not user post", http.StatusUnauthorized)
				return
			}

			err = repository.DeletePost(r.Context(), postID, userID)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Post deleted successfully",
				Status:  true,
			})
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func ListPostsHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("pageSize")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		// Calculate the offset based on the page number and page size
		offset := (page - 1) * pageSize

		posts, err := repository.ListPost(r.Context(), int64(pageSize), int64(offset))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}
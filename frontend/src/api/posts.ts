import { client } from "./client";
import { Post } from "../types/post";

export const postApi = {
  // Get posts for a specific module
  getByModule: (moduleId: number) =>
    client<Post[]>(`/modules/${moduleId}/posts`),

  // Create a new post (Requires Auth)
  createPost: (data: { module_id: number; title: string; content: string }) =>
    client<Post>(`/api/modules/${data.module_id}/posts`, { body: data }),

  // Update an existing post
  update: (id: number, data: { title: string; content: string }) =>
    client(`/api/posts/${id}`, { method: "PUT", body: data }),

  // Delete a post (soft delete)
  delete: (id: number) => client(`/api/posts/${id}`, { method: "DELETE" }),

  // Get a single post with its details
  getPost: (postId: number) => client<Post>(`/posts/${postId}`),
};

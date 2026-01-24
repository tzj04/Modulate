import { client } from "./client";
import { Comment } from "../types/comment";

export const commentApi = {
  // Fetch the flat list of comments for a specific post
  getThread: (postId: number) =>
    client<Comment[]>(`/posts/${postId}/comments/thread`),

  // Create a new comment
  create: (
    postId: number,
    data: { content: string; parent_comment_id: number | null },
  ) => client<Comment>(`/api/posts/${postId}/comments`, { body: data }),
};

import { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { usePost } from "../../hooks/usePost";
import { useComments } from "../../hooks/useComments";
import { CommentList } from "../../components/comments/CommentList";
import { commentApi } from "../../api/comments";

export const PostDetail = () => {
  const { postId } = useParams<{ postId: string }>();
  const navigate = useNavigate();
  const numericPostId = postId ? parseInt(postId, 10) : undefined;

  const { post, loading: postLoading, error } = usePost(numericPostId);
  const {
    comments,
    loading: commentsLoading,
    refresh,
  } = useComments(numericPostId);
  const [newComment, setNewComment] = useState("");

  // Handles both main comments (parentId = null) and nested replies
  const handlePostComment = async (
    content: string,
    parentId: number | null = null,
  ) => {
    if (!numericPostId || !content.trim()) return;
    try {
      await commentApi.create(numericPostId, {
        content: content,
        parent_comment_id: parentId,
      });
      setNewComment("");
      refresh(); // Refresh thread after successful post
    } catch (err) {
      console.error("Failed to post comment:", err);
      alert("Could not post comment. Please try again.");
    }
  };

  if (postLoading)
    return <div className="post-detail-container">Loading thread...</div>;
  if (error)
    return (
      <div className="post-detail-container" style={{ color: "red" }}>
        {error}
      </div>
    );
  if (!post)
    return <div className="post-detail-container">Post not found.</div>;

  return (
    <div className="post-detail-container">
      <button onClick={() => navigate(-1)} className="back-btn">
        ← Back to list
      </button>

      <article className="post-article">
        <h1 className="post-title">{post.title}</h1>
        <div className="post-meta">
          Posted by <strong>{post.username}</strong> on{" "}
          {new Date(post.created_at).toLocaleString()}
        </div>
        <div className="post-content">{post.content}</div>
      </article>

      <section className="comment-section">
        <h3 style={{ marginBottom: "15px" }}>Discussion ({comments.length})</h3>

        {/* Top-level Comment Input */}
        <div className="comment-input-area">
          <textarea
            className="comment-textarea"
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder="Write a comment..."
          />
          <button
            className="comment-submit-btn"
            onClick={() => handlePostComment(newComment, null)}
          >
            Post Comment
          </button>
        </div>

        {commentsLoading ? (
          <p style={{ color: "#999" }}>Loading thread...</p>
        ) : (
          <CommentList
            comments={comments}
            onReplySubmit={handlePostComment} // Passes the handler down to recursive items
          />
        )}
      </section>
    </div>
  );
};

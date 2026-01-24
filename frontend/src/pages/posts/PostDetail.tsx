import { useState, useEffect, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { usePost } from "../../hooks/usePost";
import { useComments } from "../../hooks/useComments";
import { CommentList } from "../../components/comments/CommentList";
import { commentApi } from "../../api/comments";
import { postApi } from "../../api/posts";
import { useAuth } from "../../auth/useAuth";

export const PostDetail = () => {
  const { postId } = useParams<{ postId: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const numericPostId = postId ? parseInt(postId, 10) : undefined;

  const [showMenu, setShowMenu] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  // Close menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setShowMenu(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const {
    post,
    loading: postLoading,
    error,
    refresh: refreshPost,
  } = usePost(numericPostId);

  const {
    comments,
    loading: commentsLoading,
    refresh: refreshComments,
  } = useComments(numericPostId);

  const [isEditing, setIsEditing] = useState(false);
  const [editTitle, setEditTitle] = useState("");
  const [editContent, setEditContent] = useState("");
  const [newComment, setNewComment] = useState("");

  const handleStartEdit = () => {
    if (!post) return;
    setEditTitle(post.title);
    setEditContent(post.content);
    setIsEditing(true);
    setShowMenu(false); // Close menu if it was open
  };

  const handleSaveEdit = async () => {
    if (!numericPostId) return;
    try {
      await postApi.update(numericPostId, {
        title: editTitle,
        content: editContent,
      });
      setIsEditing(false);
      refreshPost();
    } catch (err) {
      alert("Failed to update post.");
    }
  };

  const handleDelete = async () => {
    if (!window.confirm("Remove the content of this post?")) return;

    try {
      await postApi.delete(numericPostId!);
      setShowMenu(false);

      refreshPost();
    } catch (err) {
      alert("Error removing post content.");
    }
  };

  const handlePostComment = async (
    content: string,
    parentId: number | null = null,
  ) => {
    if (!numericPostId || !content.trim()) return;
    try {
      await commentApi.create(numericPostId, {
        content,
        parent_comment_id: parentId,
      });
      setNewComment("");
      refreshComments();
    } catch (err) {
      alert("Could not post comment.");
    }
  };

  if (postLoading)
    return <div className="post-detail-container">Loading thread...</div>;
  if (error || !post)
    return (
      <div className="post-detail-container" style={{ color: "red" }}>
        {error || "Post not found"}
      </div>
    );

  const isEdited =
    post.updated_at &&
    new Date(post.updated_at).getTime() >
      new Date(post.created_at).getTime() + 1000;

  const isDeleted = post.title === "[Deleted]";

  return (
    <div className="post-detail-container">
      <button onClick={() => navigate(-1)} className="back-btn">
        ← Back
      </button>

      <article className="post-article" style={{ position: "relative" }}>
        {/* TOP RIGHT MENU - Author Only */}
        {user?.id === post.user_id && !isEditing && !isDeleted && (
          <div className="post-options-wrapper" ref={menuRef}>
            <button
              className="kebab-btn"
              onClick={() => setShowMenu(!showMenu)}
            >
              ⋮
            </button>
            {showMenu && (
              <div className="post-dropdown">
                <button onClick={handleStartEdit}>Edit Post</button>
                <button className="delete-btn-danger" onClick={handleDelete}>
                  Delete Post
                </button>
              </div>
            )}
          </div>
        )}

        {isEditing ? (
          <div className="edit-form">
            <input
              className="edit-title-input"
              value={editTitle}
              onChange={(e) => setEditTitle(e.target.value)}
            />
            <textarea
              className="edit-content-textarea"
              value={editContent}
              onChange={(e) => setEditContent(e.target.value)}
            />
            <div className="edit-actions">
              <button onClick={handleSaveEdit} className="save-btn">
                Save Changes
              </button>
              <button
                onClick={() => setIsEditing(false)}
                className="cancel-btn"
              >
                Cancel
              </button>
            </div>
          </div>
        ) : (
          <>
            <h1 className={isDeleted ? "post-title deleted" : "post-title"}>
              {post.title}
            </h1>
            <div className="post-meta">
              {/* Hide username if deleted, or keep it based on preference */}
              {isDeleted ? "Post deleted" : `Posted by ${post.username}`}
            </div>
            <div
              className={
                isDeleted ? "post-content deleted-text" : "post-content"
              }
            >
              {post.content}
            </div>
          </>
        )}
      </article>

      <section className="comment-section">
        <h3>Discussion ({comments.length})</h3>
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
          <p>Loading comments...</p>
        ) : (
          <CommentList comments={comments} onReplySubmit={handlePostComment} />
        )}
      </section>
    </div>
  );
};

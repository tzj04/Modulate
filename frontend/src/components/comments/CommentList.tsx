import { useState, useRef, useEffect } from "react";
import { Comment } from "../../types/comment";
import { useAuth } from "../../auth/useAuth";
import { commentApi } from "../../api/comments";

interface Props {
  comments: Comment[];
  onReplySubmit: (content: string, parentId: number) => Promise<void>;
  refreshThread: () => void; // reload data after edit/delete
}

export const CommentList = ({
  comments,
  onReplySubmit,
  refreshThread,
}: Props) => {
  return (
    <div className="comment-list">
      {comments.map((comment) => (
        <CommentItem
          key={comment.id}
          comment={comment}
          onReplySubmit={onReplySubmit}
          refreshThread={refreshThread}
        />
      ))}
    </div>
  );
};

const CommentItem = ({
  comment,
  onReplySubmit,
  refreshThread,
}: {
  comment: Comment;
  onReplySubmit: (content: string, parentId: number) => Promise<void>;
  refreshThread: () => void;
}) => {
  const { user } = useAuth();
  const [isReplying, setIsReplying] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [showMenu, setShowMenu] = useState(false);
  const [isCollapsed, setIsCollapsed] = useState(false);

  const [replyText, setReplyText] = useState("");
  const [editText, setEditText] = useState(comment.content);
  const menuRef = useRef<HTMLDivElement>(null);

  const isAuthor = user?.id === comment.user_id;
  const isDeleted = comment.is_deleted;
  const hasChildren = comment.children && comment.children.length > 0;
  const isEdited =
    comment.updated_at &&
    new Date(comment.updated_at).getTime() >
      new Date(comment.created_at).getTime() + 1000;

  // Close menu on click outside
  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node))
        setShowMenu(false);
    };
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  const handleUpdate = async () => {
    try {
      await commentApi.update(comment.id, editText);
      setIsEditing(false);
      refreshThread();
    } catch (err) {
      alert("Failed to edit");
    }
  };

  const handleDelete = async () => {
    if (!window.confirm("Delete this comment?")) return;
    try {
      await commentApi.delete(comment.id);
      refreshThread();
    } catch (err) {
      alert("Failed to delete");
    }
  };

  return (
    <div className={`comment-node ${isDeleted ? "ghosted" : ""}`}>
      <div className="comment-wrapper">
        {/* Collapse/Expand Button - Outside and to the left */}
        {hasChildren && (
          <button
            className="collapse-btn"
            onClick={() => setIsCollapsed(!isCollapsed)}
            title={isCollapsed ? "Expand thread" : "Collapse thread"}
          >
            {isCollapsed ? "▶" : "▼"}
          </button>
        )}
        {/* Placeholder when no collapse button to maintain alignment */}
        {!hasChildren && <div className="collapse-btn-placeholder"></div>}

        <div className="comment-card">
          {/* Comment Header with Meta and Kebab Menu */}
          <div className="comment-header">
            <div className="comment-meta">
              <strong>{isDeleted ? "[deleted]" : comment.username}</strong>
              {isEdited && <span className="edited-tag">(edited)</span>}
            </div>

            <div className="comment-header-right">
              <div className="comment-date">
                {new Date(comment.created_at).toLocaleDateString()}
              </div>

              {/* Kebab Menu for Author - Right side */}
              {isAuthor && !isDeleted && !isEditing && (
                <div className="comment-menu-wrapper" ref={menuRef}>
                  <button
                    className="comment-kebab-btn"
                    onClick={() => setShowMenu(!showMenu)}
                  >
                    ⋮
                  </button>
                  {showMenu && (
                    <div className="comment-dropdown">
                      <button
                        onClick={() => {
                          setIsEditing(true);
                          setShowMenu(false);
                        }}
                        className="comment-menu-item"
                      >
                        <img
                          src="/images/edit_logo.png"
                          alt="edit"
                          className="comment-menu-icon"
                        />
                        Edit
                      </button>
                      <button
                        className="comment-menu-item delete-text"
                        onClick={handleDelete}
                      >
                        <img
                          src="/images/delete_logo_red.png"
                          alt="delete"
                          className="comment-menu-icon"
                        />
                        Delete
                      </button>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>

          {isEditing ? (
            <div className="edit-comment-area">
              <textarea
                value={editText}
                onChange={(e) => setEditText(e.target.value)}
                className="reply-textarea"
              />
              <button onClick={handleUpdate} className="reply-submit-btn">
                Save
              </button>
              <button
                onClick={() => setIsEditing(false)}
                className="cancel-btn"
              >
                Cancel
              </button>
            </div>
          ) : !isDeleted ? (
            <p className="comment-content">{comment.content}</p>
          ) : (
            <p className="comment-content deleted-text">
              [this comment has been deleted by the author]
            </p>
          )}

          {/* Hide reply button if deleted */}
          {!isDeleted && (
            <button
              onClick={() => setIsReplying(!isReplying)}
              className="comment-reply-btn"
            >
              {isReplying ? "Cancel" : "Reply"}
            </button>
          )}

          {isReplying && (
            <div className="reply-input-container">
              <textarea
                className="reply-textarea"
                value={replyText}
                onChange={(e) => setReplyText(e.target.value)}
                placeholder="Write a reply..."
              />
              <button
                onClick={async () => {
                  await onReplySubmit(replyText, comment.id);
                  setReplyText("");
                  setIsReplying(false);
                }}
                className="reply-submit-btn"
              >
                Submit Reply
              </button>
            </div>
          )}
        </div>
      </div>

      {comment.children && comment.children.length > 0 && !isCollapsed && (
        <div className="comment-children">
          <CommentList
            comments={comment.children}
            onReplySubmit={onReplySubmit}
            refreshThread={refreshThread}
          />
        </div>
      )}
    </div>
  );
};

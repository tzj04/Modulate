import { useState } from "react";
import { Comment } from "../../types/comment";

interface Props {
  comments: Comment[];
  onReplySubmit: (content: string, parentId: number) => Promise<void>;
}

export const CommentList = ({ comments, onReplySubmit }: Props) => {
  return (
    <div className="comment-list">
      {comments.map((comment) => (
        <CommentItem
          key={comment.id}
          comment={comment}
          onReplySubmit={onReplySubmit}
        />
      ))}
    </div>
  );
};

const CommentItem = ({
  comment,
  onReplySubmit,
}: {
  comment: Comment;
  onReplySubmit: (content: string, parentId: number) => Promise<void>;
}) => {
  const [isReplying, setIsReplying] = useState(false);
  const [replyText, setReplyText] = useState("");

  const handleSubmit = async () => {
    if (!replyText.trim()) return; // Don't submit empty replies

    await onReplySubmit(replyText, comment.id);
    setReplyText("");
    setIsReplying(false);
  };

  return (
    <div className="comment-node">
      <div className="comment-card">
        <div className="comment-meta">
          <strong>{comment.username}</strong> •{" "}
          {new Date(comment.created_at).toLocaleDateString()}
        </div>

        <p className="comment-content">{comment.content}</p>

        <button
          onClick={() => setIsReplying(!isReplying)}
          className="comment-reply-btn"
        >
          {isReplying ? "Cancel" : "Reply"}
        </button>

        {isReplying && (
          <div className="reply-input-container">
            <textarea
              className="reply-textarea"
              value={replyText}
              onChange={(e) => setReplyText(e.target.value)}
              placeholder={`Replying to ${comment.username}...`}
            />
            <button onClick={handleSubmit} className="reply-submit-btn">
              Submit Reply
            </button>
          </div>
        )}
      </div>

      {/* Recursive rendering for children */}
      {comment.children && comment.children.length > 0 && (
        <div className="comment-children">
          <CommentList
            comments={comment.children}
            onReplySubmit={onReplySubmit}
          />
        </div>
      )}
    </div>
  );
};

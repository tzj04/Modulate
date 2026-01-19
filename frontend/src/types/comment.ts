export interface Comment {
  id: number;
  post_id: number;
  user_id: number;
  parent_comment_id: number | null;
  content: string;
  is_deleted: boolean;
  created_at: string;
  updated_at: string | null;
  children?: Comment[]; // Recursive nesting
  username?: string;
}

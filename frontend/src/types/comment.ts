export interface Comment {
  id: number;
  post_id: number;
  user_id: number;
  parent_comment_id?: number; // For nested comments
  content: string;
  is_deleted: boolean;
  created_at: string;
  updated_at: string;
  username?: string; // joined from the users table in the backend
}

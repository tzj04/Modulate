export interface Post {
  id: number;
  module_id: number;
  user_id: number;
  title: string;
  content: string;
  is_deleted: boolean;
  created_at: string;
  updated_at: string;
  username?: string; // Often joined from the users table in the backend
}

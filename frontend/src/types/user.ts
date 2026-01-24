export interface User {
  id: number;
  username: string;
  label: string | null;
  is_deleted: boolean;
  created_at: string;
}

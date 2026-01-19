// src/types/user.ts
export interface User {
  id: number;
  username: string;
  label: string | null; // Matches the 'label' TEXT column
  is_deleted: boolean; // Matches the 'is_deleted' BOOLEAN column
  created_at: string; // TIMESTAMPTZ comes over as an ISO string
}

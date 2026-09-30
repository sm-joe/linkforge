export type LinkStatus =
  | "active"
  | "disabled"
  | "expired";

export interface Link {
  id: string;
  short_code: string;
  short_url: string;
  destination: string;
  status: LinkStatus;
  created_at: string;
  expires_at?: string | null;
}
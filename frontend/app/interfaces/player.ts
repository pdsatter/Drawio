export interface Player {
  id?: string | null;
  guestId?: string;
  name: string;
  icon: string;
}

export interface GuestPlayer extends Player {
  id: null;
  guestId: string;
}

export interface LoggedInPlayer extends Player {
  id: string;
  guestId?: never;
}

export type AuthenticatedPlayer = GuestPlayer | LoggedInPlayer;
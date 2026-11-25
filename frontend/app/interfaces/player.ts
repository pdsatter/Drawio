export interface Player {
	name: string;
	icon: string;
}

export interface GuestPlayer extends Player {
	id: null;
	guestId: string;
}

export interface LoggedInPlayer extends Player {
  id: string;
}

export interface Player {
	name: string;
	icon: string;
}

export interface GuestPlayer extends Player {
  id: null;
}

export interface LoggedInPlayer extends Player {
  id: string;
}

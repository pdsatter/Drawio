interface Player {
	name: string;
	icon: string;
}

interface GuestPlayer extends Player {
  id: null;
}

interface LoggedInPlayer extends Player {
  id: string;
}

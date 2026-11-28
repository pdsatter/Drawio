export type Player = {
	id?: string;
	guestId?: string;
	name: string;
	icon: string;
};


export function getPlayerKey(player: Player): string {
	return player.id
		? `user:${player.id}`
		: `guest:${player.guestId}`;
}

import { getPlayerKey, type Player } from "./player";


export type LobbyState = "waiting" | "drawing" | "swapping";

export interface Lobby {
	host: Player;
	code: string;
	players: Player[];
	messages: Message[];
	drawer: string | null;
	state: LobbyState;
	turnEndsAt: number;
	swapEndsAt: number;
	timeLeft: number;
}

export interface Message {
	player: Player;
	message: string;
}

export function isCurrentDrawer(
	lobby: Lobby,
	player: Player
): boolean {
	return lobby.drawer === getPlayerKey(player);
}


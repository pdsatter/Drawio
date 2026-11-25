export interface Message {
  player: Player;
  message: string;
}

export interface Lobby {
  host: Player;
  code: string;
  players: Player[];
  messages: Message[];
}

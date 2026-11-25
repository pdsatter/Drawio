export interface Message {
  player: Player;
  message: string;
}

export interface Lobby {
  code: string;
  players: Player[];
  messages: Message[];
}

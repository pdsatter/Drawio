// Backend will control this, but for visual front end needs to know.
export const REWARDS = {
  easy: 800,
  medium: 900,
  hard: 1000,
} as const;

export type Difficulty = keyof typeof REWARDS;

export interface Category {
  name: string;
  words: Word[];
}

export interface Word {
  name: string;
  difficulty: Difficulty;
}

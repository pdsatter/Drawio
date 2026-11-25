import backendUrl from "./backendUrl";
import Player from "../interfaces/player";

export async function createGuestSession(): Promise<Player> {
  const response = await fetch(`${backendUrl}/auth/guest`, {
    method: "POST",
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error(`Failed to create guest session: ${response.status}`);
  }

  return response.json();
}

import backendUrl from "./backendUrl";
import Lobby from "../interfaces/lobby";


export async function createLobby(): Promise<Lobby> {
  const response = await fetch(`${backendUrl}/lobbies`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to create lobby: ${response.status}`);
  }

  return response.json();
}

export async function joinLobby(code: string): Promise<Lobby> {
  const response = await fetch(
    `${backendUrl}/lobbies/${encodeURIComponent(code)}/join`,
    {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    }
  );

  if (!response.ok) {
    throw new Error(`Failed to join lobby: ${response.status}`);
  }

  return response.json();
}

export async function leaveLobby(code: string): Promise<void> {
  const response = await fetch(
    `${backendUrl}/lobbies/${encodeURIComponent(code)}/leave`,
    {
      method: "POST",
      credentials: "include",
    }
  );

  if (!response.ok) {
    throw new Error(`Failed to leave lobby: ${response.status}`);
  }
}

export async function getLobby(code: string): Promise<Lobby> {
  const response = await fetch(
    `${backendUrl}/lobbies/${encodeURIComponent(code)}`,
    {
      method: "GET",
      credentials: "include",
    }
  );

  if (!response.ok) {
    throw new Error(`Failed to retrieve lobby: ${response.status}`);
  }

  return response.json();
}

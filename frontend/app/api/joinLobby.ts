import backendUrl from "./backendUrl";

export async function joinLobby(code: string) {
  try {
    const response = await fetch(`${backendUrl}/lobbies/${code}/join`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to join lobby: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error("Error joining lobby:", error);
    throw error;
  }
}

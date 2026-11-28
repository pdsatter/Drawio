import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import { getLobby, leaveLobby } from "../api/lobby";
import type { Lobby as LobbyType } from "../interfaces/lobby";

export default function Lobby() {
  const { roomCode } = useParams();
  const navigate = useNavigate();
  const [lobby, setLobby] = useState<LobbyType | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!roomCode) return;

    let active = true;

    async function refreshLobby() {
      try {
        if (!roomCode) return;
        const currentLobby = await getLobby(roomCode);
        if (active) {
          setLobby(currentLobby);
          setError("");
        }
      } catch (requestError) {
        if (active) {
          setError(requestError instanceof Error ? requestError.message : "Could not load lobby");
        }
      }
    }


    refreshLobby();
    const interval = window.setInterval(refreshLobby, 2000);

    return () => {
      active = false;
      window.clearInterval(interval);
    };
  }, [roomCode]);

  async function handleLeaveLobby() {
    if (!roomCode) return;

    try {
      await leaveLobby(roomCode);
      navigate("/");
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : "Could not leave lobby");
    }
  }

  async function handleStartGame() {
    if (!roomCode) return;

    try {
      const response = await fetch(`/api/lobbies/${roomCode}/start`, {
        method: "POST",
        credentials: "include",
      });

      if (!response.ok) {
        const data = await response.json().catch(() => null);
        throw new Error(data?.error ?? "Could not start game.");
      }

      navigate(`/game/${roomCode}`);
    } catch (requestError) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "Could not start game."
      );
    }
  }


  return (
    <main className="min-h-screen bg-[#F8F3E8] text-[#252525]">
      <header className="mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
        <Link to="/" className="text-2xl font-black tracking-tight">
          <span className="text-[#F28C45]">Draw</span> with Friends
          <span className="ml-1">✏️</span>
        </Link>

        <Link
          to="/avatar"
          className="flex items-center gap-2 rounded-full border-2 border-[#E5DCCB] bg-white px-4 py-2 font-bold transition hover:border-[#252525]"
        >
          <span className="text-xl">🐱</span>
          <span className="hidden sm:inline">Avatar</span>
        </Link>
      </header>

      <section className="mx-auto flex max-w-5xl flex-col items-center px-6 pb-20 pt-16 text-center">
        <h1 className="text-5xl font-black leading-tight tracking-tight sm:text-7xl">
          Lobby
        </h1>
        <p className="mt-4 text-lg text-[#777064]">Share this code with friends</p>

        <button
          type="button"
          onClick={() => navigator.clipboard?.writeText(lobby?.code ?? roomCode ?? "")}
          className="mt-6 rounded-full border-2 border-[#252525] bg-[#E9C84A] px-5 py-3 text-sm font-black shadow-[3px_3px_0px_#252525]"
          aria-label="Copy lobby code"
        >
          Game Code: <span className="font-bold">{lobby?.code ?? roomCode}</span>
        </button>

        {error && (
          <p role="alert" className="mt-6 rounded-xl border-2 border-red-300 bg-white px-4 py-3 font-bold text-red-700">
            {error}
          </p>
        )}

        <div className="mt-10 flex min-h-20 flex-wrap justify-center gap-4">
          {lobby?.players.map((player, index) => (
            <div
              key={player.guestId ?? player.id ?? `${player.name}-${index}`}
              className="flex h-16 w-16 items-center justify-center rounded-full border-2 border-[#252525] bg-[#78A89A] text-3xl"
              title={player.name}
            >
              {player.icon === "default" ? "👤" : player.icon}
            </div>
          ))}
        </div>

        <div className="mt-10 flex flex-wrap justify-center gap-4">
          <button
            type="button"
            onClick={handleLeaveLobby}
            className="rounded-full border-2 border-[#252525] bg-white px-6 py-3 font-black transition hover:-translate-y-1"
          >
            Leave Lobby
          </button>
          <button
            type="button"
            onClick={handleStartGame}
            disabled={!lobby || lobby.state !== "waiting"}
            className="rounded-full border-2 border-[#252525] bg-[#F28C45] px-6 py-3 font-black shadow-[3px_3px_0px_#252525] transition hover:-translate-y-1 disabled:opacity-50"
          >
            Start Game
          </button>
        </div>
      </section>
    </main>
  );
}

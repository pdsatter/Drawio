
import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router";
import { createLobby, joinLobby } from "../api/lobby";
import { useAuth } from "../context/AuthContext";

export function meta() {
  return [
    { title: "Draw with Friends!" },
    {
      name: "description",
      content: "Draw images with friends, and try to guess them as fast as possible.",
    },
  ];
}

export default function Home() {
  const navigate = useNavigate();
  const { player, loading: authLoading } = useAuth();
  const [joinCode, setJoinCode] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");

  async function handleCreateLobby() {
    setPending(true);
    setError("");

    try {
      const lobby = await createLobby();
      navigate(`/lobby/${lobby.code}`);
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : "Could not create lobby");
      setPending(false);
    }
  }

  async function handleJoinLobby(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    await submitJoinLobby();
  }

  async function submitJoinLobby() {
    const code = joinCode.trim().toUpperCase();
    if (!code) {
      setError("Enter a lobby code");
      return;
    }

    setPending(true);
    setError("");

    try {
      const lobby = await joinLobby(code);
      navigate(`/lobby/${lobby.code}`);
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : "Could not join lobby");
      setPending(false);
    }
  }

  return (
    <main className="min-h-screen bg-[#F8F3E8] text-[#252525]">
      <header className="mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
        <Link to="/" className="text-2xl font-black tracking-tight">
          <span className="text-[#F28C45]">Draw</span> with Friends
          <span className="ml-1">✏️</span>
        </Link>

        <nav className="flex items-center gap-3">
          <Link
            to="/avatar"
            className="flex items-center gap-2 rounded-full border-2 border-[#E5DCCB] bg-white px-4 py-2 font-bold transition hover:border-[#252525]"
          >
            <span className="text-xl">🐱</span>
            <span className="hidden sm:inline">Avatar</span>
          </Link>

          <button className="rounded-full bg-[#252525] px-5 py-2 font-bold text-white transition hover:bg-[#F28C45] hover:text-[#252525]">
            Sign in
          </button>
        </nav>
      </header>

      <section className="mx-auto flex max-w-5xl flex-col items-center px-6 pb-20 pt-16 text-center sm:pt-24">

        <div className="rotate-[-2deg] rounded-full border-2 border-[#252525] bg-[#E9C84A] px-5 py-2 text-sm font-black shadow-[3px_3px_0px_#252525]">
          ✏️ DRAW SOMETHING RIDICULOUS!
        </div>

        <h1 className="mt-10 text-5xl font-black leading-tight tracking-tight sm:text-7xl">
          Your imagination.
          <br />
          <span className="text-[#F28C45]">Their guess.</span>
        </h1>

        <p className="mt-6 max-w-lg text-lg leading-relaxed text-[#777064]">
          Grab your friends, pick up your imaginary pencils,
          and see who can draw and guess the fastest!
        </p>

        <div className="mt-14 grid w-full max-w-2xl gap-6 sm:grid-cols-2">
          <button
            type="button"
            onClick={handleCreateLobby}
            disabled={authLoading || !player || pending}
            className="group rounded-3xl border-2 border-[#252525] bg-[#F28C45] p-8 text-left shadow-[6px_6px_0px_#252525] transition hover:-translate-y-1 hover:shadow-[8px_8px_0px_#252525] disabled:cursor-not-allowed disabled:opacity-60"
          >
            <div className="text-5xl transition group-hover:rotate-[-10deg]">
              ✏️
            </div>
            <h2 className="mt-8 text-3xl font-black">Create Game</h2>
            <p className="mt-3 text-[#493526]">
              Start a room and invite your friends to play.
            </p>
            <div className="mt-8 font-black">
              {pending ? "Creating..." : "Let's play →"}
            </div>
          </button>

          <form
            onSubmit={handleJoinLobby}
            className="rounded-3xl border-2 border-[#252525] bg-white p-8 text-left shadow-[6px_6px_0px_#252525]"
          >
            <div className="text-5xl">🎯</div>
            <h2 className="mt-8 text-3xl font-black">Join Game</h2>
            <p className="mt-3 text-[#777064]">
              Enter the room code from your friend's lobby.
            </p>
            <div className="mt-8 flex items-center gap-2">
              <input
                value={joinCode}
                onChange={(event) => setJoinCode(event.target.value.replace(/[^a-zA-Z0-9]/g, "").slice(0, 6).toUpperCase())}
                placeholder="Lobby code"
                aria-label="Lobby code"
                disabled={pending}
                className="min-w-0 flex-1 rounded-xl border-2 border-[#252525] px-4 py-3 text-lg font-black uppercase outline-none focus:border-[#F28C45]"
              />
              <button
                type="submit"
                disabled={authLoading || !player || pending}
                aria-label="Join lobby"
                className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full border-2 border-[#252525] bg-[#F28C45] text-2xl font-black transition hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60"
              >
                →
              </button>
            </div>
          </form>
        </div>

        {error && (
          <p role="alert" className="mt-6 rounded-xl border-2 border-red-300 bg-white px-4 py-3 font-bold text-red-700">
            {error}
          </p>
        )}

        <div className="mt-14 flex w-full max-w-2xl flex-col items-center justify-between gap-5 rounded-2xl border-2 border-[#E5DCCB] bg-white p-6 sm:flex-row">

          <Link to="/avatar" className="flex items-center gap-4">
            <div className="flex h-14 w-14 items-center justify-center rounded-full border-2 border-[#252525] bg-[#78A89A] text-3xl">
              🐱
            </div>

            <div className="text-left">
              <p className="font-black">Your Character</p>
              <p className="text-sm text-[#777064]">
                Customize your avatar →
              </p>
            </div>
          </Link>

          <div className="hidden h-12 w-px bg-[#E5DCCB] sm:block" />

          <div className="flex items-center gap-3">
            <div className="text-3xl">🏆</div>
            <div className="text-left">
              <p className="text-sm text-[#777064]">Total Points</p>
              <p className="text-xl font-black">0 points</p>
            </div>
          </div>
        </div>

        <div className="mt-20 flex items-center gap-6 text-3xl opacity-60">
          <span className="rotate-[-15deg]">🖍️</span>
          <span className="rotate-[10deg]">🎨</span>
          <span className="rotate-[-8deg]">📝</span>
          <span className="rotate-[15deg]">✂️</span>
        </div>

        <p className="mt-5 text-sm font-medium text-[#A69B89]">
          Made for friends who can't draw.
        </p>
      </section>
    </main>
  );
}

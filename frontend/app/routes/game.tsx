
import { useCallback, useEffect, useRef, useState } from "react";
import { useParams, Link } from "react-router";
import { useAuth } from "~/context/AuthContext";
import type { Lobby } from "~/interfaces/lobby";
import type { Player } from "~/interfaces/player";


type GameResponse = {
    lobby: Lobby;
    animal?: string;
};

function getPlayerKey(player: Player): string {
    if (player.id) return `user:${player.id}`;
    return `guest:${player.guestId}`;
}

export default function Game() {
    const { roomCode } = useParams();
    const { player, loading } = useAuth();

    if (loading) {
        return <div>Loading player session...</div>;
    }
    if (!player) {
        return <div>Could not load player session. Please refresh.</div>;
    }

    const currentPlayer = player as Player | null;
    const [lobby, setLobby] = useState<Lobby | null>(null);

    const [targetAnimal, setTargetAnimal] = useState<string | null>(null);
    const [guess, setGuess] = useState("");
    const [secondsLeft, setSecondsLeft] = useState(0);
    const [points, setPoints] = useState(0);
    const [msgInput, setMsgInput] = useState("");
    const [error, setError] = useState<string | null>(null);

    const canvasRef = useRef<HTMLCanvasElement>(null);

    const isDrawer =
        !!lobby &&
        !!currentPlayer &&
        lobby.drawer === getPlayerKey(currentPlayer);

    const isDrawing = lobby?.state === "drawing";
    const isSwapping = lobby?.state === "swapping";

    const canDraw = isDrawer && isDrawing;
    const canGuess = !isDrawer && isDrawing && !!currentPlayer;

    const canChat =
        !!lobby &&
        !!currentPlayer &&
        (lobby.players.length === 1 || !isDrawer);


    // Fetch lobby information.
    const refreshLobby = useCallback(async () => {
        console.log("Game ID:", roomCode);
        if (!roomCode) return;

        try {
            const response = await fetch(
                `/api/lobbies/${roomCode}`,
                { credentials: "include" }
            );

            if (!response.ok) {
                throw new Error("Failed to retrieve lobby.");
            }

            const data: GameResponse = await response.json();

            setLobby(data.lobby ?? data);
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Failed to load lobby."
            );
        }
    }, [roomCode]);

    // Poll for lobby changes.
    useEffect(() => {
        refreshLobby();

        const interval = window.setInterval(refreshLobby, 1000);

        return () => window.clearInterval(interval);
    }, [refreshLobby]);

    // Fetch the secret animal only for the current drawer.
    useEffect(() => {
        if (!roomCode || !isDrawer || !isDrawing) {
            setTargetAnimal(null);
            return;
        }

        let cancelled = false;

        const fetchAnimal = async () => {
            try {
                const response = await fetch(
                    `/api/lobbies/${roomCode}/animal`,
                    { credentials: "include" }
                );

                if (!response.ok) {
                    throw new Error("Unable to retrieve drawing prompt.");
                }

                const data: { animal: string } = await response.json();

                if (!cancelled) {
                    setTargetAnimal(data.animal);
                }
            } catch (err) {
                if (!cancelled) {
                    setError(
                        err instanceof Error
                            ? err.message
                            : "Unable to retrieve drawing prompt."
                    );
                }
            }
        };

        fetchAnimal();

        return () => {
            cancelled = true;
        };
    }, [
        roomCode,
        isDrawer,
        isDrawing,
        lobby?.turnEndsAt
    ]);

    // Display-only countdown based on backend deadlines.
    useEffect(() => {
        if (!lobby) return;

        const updateCountdown = () => {
            const deadline =
                lobby.state === "drawing"
                    ? lobby.turnEndsAt
                    : lobby.state === "swapping"
                        ? lobby.swapEndsAt
                        : 0;

            if (!deadline) {
                setSecondsLeft(0);
                return;
            }

            setSecondsLeft(
                Math.max(
                    0,
                    Math.ceil((deadline - Date.now()) / 1000)
                )
            );
        };

        updateCountdown();

        const interval = window.setInterval(updateCountdown, 250);

        return () => window.clearInterval(interval);
    }, [lobby?.state, lobby?.turnEndsAt, lobby?.swapEndsAt]);

    // Drawing canvas.
    useEffect(() => {
        const canvas = canvasRef.current;
        if (!canvas) return;

        const ctx = canvas.getContext("2d");
        if (!ctx) return;

        ctx.lineWidth = 3;
        ctx.lineCap = "round";
        ctx.strokeStyle = "#000";

        let drawing = false;

        const start = (e: MouseEvent) => {
            if (!canDraw) return;

            drawing = true;
            ctx.beginPath();
            ctx.moveTo(e.offsetX, e.offsetY);
        };

        const draw = (e: MouseEvent) => {
            if (!drawing || !canDraw) return;

            ctx.lineTo(e.offsetX, e.offsetY);
            ctx.stroke();
        };

        const stop = () => {
            drawing = false;
        };

        canvas.addEventListener("mousedown", start);
        canvas.addEventListener("mousemove", draw);
        canvas.addEventListener("mouseup", stop);
        canvas.addEventListener("mouseleave", stop);

        return () => {
            canvas.removeEventListener("mousedown", start);
            canvas.removeEventListener("mousemove", draw);
            canvas.removeEventListener("mouseup", stop);
            canvas.removeEventListener("mouseleave", stop);
        };
    }, [canDraw]);

    const clearCanvas = () => {
        const canvas = canvasRef.current;

        if (!canvas) return;

        const ctx = canvas.getContext("2d");

        ctx?.clearRect(0, 0, canvas.width, canvas.height);
    };

    // Submit a guess. The backend validates the answer.
    const handleGuess = async () => {
        if (!roomCode || !canGuess || !guess.trim()) return;

        try {
            const response = await fetch(
                `/api/lobbies/${roomCode}/guess`,
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    credentials: "include",
                    body: JSON.stringify({
                        guess: guess.trim(),
                    }),
                }
            );

            if (!response.ok) {
                throw new Error("Failed to submit guess.");
            }

            const data = await response.json();

            if (data.correct) {
                setPoints(prev => prev + (data.points ?? 0));
                clearCanvas();
            }

            setGuess("");
            await refreshLobby();
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Guess failed."
            );
        }
    };

    // Submit chat message.
    const handleSendMessage = async () => {
        if (!roomCode || !canChat || !msgInput.trim()) return;

        try {
            const response = await fetch(
                `/api/lobbies/${roomCode}/messages`,
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    credentials: "include",
                    body: JSON.stringify({
                        message: msgInput.trim(),
                    }),
                }
            );

            if (!response.ok) {
                throw new Error("Unable to send message.");
            }

            setMsgInput("");
            await refreshLobby();
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Message failed."
            );
        }
    };

    const nextDrawer = lobby?.players.find(
        player => getPlayerKey(player) === lobby.drawer
    );

    if (error) {
        return (
            <main className="min-h-screen bg-[#F8F3E8] flex flex-col items-center justify-center">
                <h1 className="text-3xl font-black">Game Error</h1>
                <p className="mt-3">{error}</p>
                <button
                    onClick={() => {
                        setError(null);
                        refreshLobby();
                    }}
                    className="mt-5 rounded bg-[#F28C45] px-5 py-2 font-bold text-white"
                >
                    Retry
                </button>
            </main>
        );
    }

    if (!lobby || !currentPlayer) {
        return (
            <main className="min-h-screen bg-[#F8F3E8] flex items-center justify-center">
                <h1 className="text-2xl font-bold">Loading game...</h1>
                
            </main>
        );
    }

    // Waiting for the host to start.
    if (lobby.state === "waiting") {
        return (
            <main className="min-h-screen bg-[#F8F3E8] flex flex-col items-center justify-center">
                <h1 className="text-4xl font-black">Waiting for game...</h1>
                <p className="mt-4 text-lg">
                    Players: {lobby.players.length}
                </p>
            </main>
        );
    }

    // Drawer transition screen.
    if (isSwapping) {
        return (
            <main className="min-h-screen bg-[#F8F3E8] flex flex-col items-center justify-center text-[#252525]">
                <h1 className="text-5xl font-black">
                    Switching Drawer...
                </h1>

                <p className="mt-4 text-2xl font-bold">
                    {nextDrawer?.icon} {nextDrawer?.name} is up next!
                </p>

                <p className="mt-6 text-4xl font-black text-[#F28C45]">
                    {secondsLeft}
                </p>

                <p className="mt-2 text-gray-500">
                    Get ready!
                </p>
            </main>
        );
    }

    return (
        <main className="min-h-screen bg-[#F8F3E8] text-[#252525]">
            <header className="mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
                <Link to="/" className="text-2xl font-black tracking-tight">
                    <span className="text-[#F28C45]">Draw</span> with Friends
                    <span className="ml-1">✏️</span>
                </Link>

                <div className="flex items-center gap-5">
                    <span className="font-bold">
                        {currentPlayer.name}
                    </span>

                    <span className="rounded-full bg-[#252525] px-5 py-2 font-bold text-white">
                        {secondsLeft}s
                    </span>
                </div>
            </header>

            <section className="mx-auto flex max-w-5xl gap-6 px-6 pb-20">
                <div className="flex w-3/4 flex-col">
                    <h1 className="mb-2 text-5xl font-black leading-tight tracking-tight sm:text-7xl">
                        {isDrawer
                            ? `Draw a ${targetAnimal ?? "..."}`
                            : "Guess the animal!"}
                    </h1>

                    <canvas
                        ref={canvasRef}
                        width={800}
                        height={600}
                        className={`border-2 border-[#252525] bg-white ${
                            canDraw ? "cursor-crosshair" : "cursor-default"
                        }`}
                    />

                    {isDrawer ? (
                        <p className="mt-4 font-bold">
                            You are drawing!
                        </p>
                    ) : (
                        <div className="mt-4 flex items-center gap-2">
                            <input
                                type="text"
                                value={guess}
                                onChange={e => setGuess(e.target.value)}
                                onKeyDown={e => {
                                    if (e.key === "Enter") handleGuess();
                                }}
                                placeholder="Guess the animal"
                                className="flex-1 rounded border px-3 py-2"
                                disabled={!canGuess}
                            />

                            <button
                                onClick={handleGuess}
                                disabled={!canGuess || !guess.trim()}
                                className="rounded bg-[#F28C45] px-4 py-2 font-bold text-white disabled:cursor-not-allowed disabled:opacity-50"
                            >
                                Guess
                            </button>
                        </div>
                    )}

                    <div className="mt-4 font-bold">
                        Your points: {points}
                    </div>
                </div>

                <aside className="flex w-80 flex-col border-l">
                    <div className="mb-2 flex items-center px-2">
                        <span className="text-lg font-bold">
                            Chat
                        </span>
                    </div>

                    <div className="relative mb-2 flex-1 overflow-y-auto border-2 border-[#252525] bg-white p-2">
                        {lobby.messages.map((message, index) => (
                            <div
                                key={index}
                                className="mb-2 flex items-start gap-2"
                            >
                                <div className="flex-shrink-0 text-2xl">
                                    {message.player.icon}
                                </div>

                                <div className="flex-1">
                                    <span className="font-bold text-gray-500">
                                        {message.player.name}:
                                    </span>{" "}
                                    {message.message}
                                </div>
                            </div>
                        ))}
                    </div>

                    <div className="flex gap-2 px-2 pb-2">
                        <input
                            type="text"
                            value={msgInput}
                            onChange={e => setMsgInput(e.target.value)}
                            onKeyDown={e => {
                                if (e.key === "Enter") handleSendMessage();
                            }}
                            placeholder="Type message"
                            disabled={!canChat}
                            className="flex-1 rounded border px-3 py-2 disabled:bg-gray-100"
                        />

                        <button
                            onClick={handleSendMessage}
                            disabled={!canChat || !msgInput.trim()}
                            className="rounded bg-[#F28C45] px-4 py-2 font-bold text-white disabled:cursor-not-allowed disabled:opacity-50"
                        >
                            Send
                        </button>
                    </div>
                </aside>
            </section>
        </main>
    );
}

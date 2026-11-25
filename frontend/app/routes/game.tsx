import { useEffect, useRef, useState } from "react";
import { useParams, Link, useNavigate } from "react-router";

// Simple list of animals for the game
const ANIMALS = ["cat", "dog", "elephant", "lion", "tiger"]; // could extend later

export default function Game() {
    const { gameId } = useParams(); // expects :gameId in route
    const navigate = useNavigate();
    const [targetAnimal, setTargetAnimal] = useState<string | null>(null);
    const [guess, setGuess] = useState<string>("");
    const [secondsLeft, setSecondsLeft] = useState<number>(30);
    const [isDrawer, setIsDrawer] = useState<boolean>(true);
    const [points, setPoints] = useState<number>(0);
	// Player data for chat - dummy players since backend not ready yet
	const DRAWER_PLAYER = { name: "Drawer", icon: "👤" };
	const GUESSER_PLAYER = { name: "Guesser", icon: "🎯" };
	
	const [messages, setMessages] = useState<{ id: number; text: string; player: typeof DRAWER_PLAYER | typeof GUESSER_PLAYER }[]>([]);
    const [msgInput, setMsgInput] = useState<string>("");
    const canvasRef = useRef<HTMLCanvasElement>(null);

    // Pick a random animal when component mounts or after a correct guess
    const pickRandomAnimal = (exclude: string | null = null) => {
        if (ANIMALS.length === 1) {
            setTargetAnimal(ANIMALS[0]);
            return;
        }
        let idx: number;
        do {
            idx = Math.floor(Math.random() * ANIMALS.length);
        } while (exclude && ANIMALS[idx] === exclude);
        setTargetAnimal(ANIMALS[idx]);
    };
    useEffect(() => {
        pickRandomAnimal();
    }, []);

    // Timer effect
    useEffect(() => {
        if (secondsLeft <= 0) {
            setSecondsLeft(0);
            setIsDrawer(false);
            return;
        }
        const timer = setInterval(() => setSecondsLeft(prev => prev - 1), 1000);
        return () => clearInterval(timer);
    }, [secondsLeft]);

    // Basic drawing logic: set up mouse events
    useEffect(() => {
        const canvas = canvasRef.current;
        if (!canvas) return;
        const ctx = canvas.getContext("2d");
        if (!ctx) return;
        ctx.lineWidth = 3;
        ctx.strokeStyle = "#000";
        let drawing = false;
        const start = (e: MouseEvent) => {
            drawing = true;
            ctx.beginPath();
            ctx.moveTo(e.offsetX, e.offsetY);
        };
        const draw = (e: MouseEvent) => {
            if (!drawing) return;
            ctx.lineTo(e.offsetX, e.offsetY);
            ctx.stroke();
        };
        const stop = () => {
            drawing = false;
        };
        canvas.addEventListener("mousedown", start as any);
        canvas.addEventListener("mousemove", draw as any);
        canvas.addEventListener("mouseup", stop as any);
        canvas.addEventListener("mouseleave", stop as any);
        return () => {
            canvas.removeEventListener("mousedown", start as any);
            canvas.removeEventListener("mousemove", draw as any);
            canvas.removeEventListener("mouseup", stop as any);
            canvas.removeEventListener("mouseleave", stop as any);
        };
    }, []);

    const handleGuess = () => {
        if (!targetAnimal) return;
        if (guess.trim().toLowerCase() === targetAnimal?.toLowerCase()) {
            // Correct guess
            const canvas = canvasRef.current;
            if (canvas) {
                const ctx = canvas.getContext("2d");
                ctx?.clearRect(0, 0, canvas.width, canvas.height);
            }
            // Calculate points: start 100, lose 100/30 per second elapsed
            const earned = Math.max(0, Math.round(100 - (30 - secondsLeft) * (100 / 30)));
            setPoints(prev => prev + earned);
            pickRandomAnimal(targetAnimal);
            // Switch roles: the guesser becomes drawer next round
            setIsDrawer(true);
        }
        setGuess("");
    };

    return (
        <main className="min-h-screen bg-[#F8F3E8] text-[#252525]">
            <header className="mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
                <Link to="/" className="text-2xl font-black tracking-tight">
                    <span className="text-[#F28C45]">Draw</span> with Friends
                    <span className="ml-1">✏️</span>
                </Link>
                <button className="rounded-full bg-[#252525] px-5 py-2 font-bold text-white transition hover:bg-[#F28C45] hover:text-[#252525]">
                    Sign in
                </button>
            </header>
<section className="mx-auto flex max-w-5xl px-6 pt-0 pb-20">
		<div className="flex flex-col w-3/4">
			<h1 className="text-5xl font-black leading-tight tracking-tight sm:text-7xl mb-2">
				Draw a {targetAnimal ?? "animal"}
			</h1>
			<canvas
				ref={canvasRef}
				width={800}
				height={600}
				className="border-2 border-[#252525] bg-white"
			/>
			<div className="mt-4 flex items-center gap-2">
				<input
					type="text"
					value={guess}
					onChange={(e) => setGuess(e.target.value)}
					placeholder="Guess animal (guesser)"
					className="flex-1 rounded border px-3 py-2"
					disabled={!isDrawer}
				/>
				<button
					onClick={handleGuess}
					disabled={!isDrawer}
					className="rounded bg-[#F28C45] px-4 py-2 font-bold text-white disabled:opacity-50 disabled:cursor-not-allowed"
				>
					Guess
				</button>
			</div>
		</div>
		<aside className="w-80 flex flex-col border-l">
			<div className="flex items-center mb-2 px-2">
				<span className="text-lg font-bold text-[#252525]">Chat</span>
			</div>
			<div className="flex-1 overflow-y-auto border-2 border-[#252525] p-2 bg-white mb-2 relative">
				{messages.map((m) => (
					<div key={m.id} className="mb-2 flex items-start gap-2">
						<div className="flex-shrink-0 text-2xl">{m.player.icon}</div>
						<div className="flex-1">
							<span className="font-bold text-gray-500">{m.player.name}:</span> {m.text}
						</div>
					</div>
				))}
			</div>
			<div className="flex gap-2 px-2 pb-2">
				<input
					type="text"
					value={msgInput}
					onChange={(e) => setMsgInput(e.target.value)}
					placeholder="Type message"
					disabled={!isDrawer}
					className="flex-1 rounded border px-3 py-2 disabled:bg-gray-100"
				/>
				<button
					onClick={() => {
						if (!msgInput) return;
						setMessages((prev) => [...prev, { id: Date.now(), text: msgInput, player: GUESSER_PLAYER }]);
						setMsgInput("");
					}}
					disabled={!isDrawer}
					className="rounded bg-[#F28C45] px-4 py-2 font-bold text-white disabled:opacity-50 disabled:cursor-not-allowed"
				>
					Send
				</button>
			</div>
		</aside>
            </section>
        </main>
    );
}

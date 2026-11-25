import { useEffect, useState } from "react";
import { Link } from "react-router";

// Dummy player data structure
interface Player {
	id: number;
	name: string;
	icon: string;
}

// Generate a simple 6‑digit game code
const generateGameCode = (): string => {
	return Math.floor(100000 + Math.random() * 900000).toString();
};

export default function CreateGame() {
	const [code, setCode] = useState<string>(generateGameCode());
	const [players, setPlayers] = useState<Player[]>([
		{ id: 1, name: "You", icon: "👤" },
	]);

	// Simulate other players joining after a delay
	useEffect(() => {
		const interval = setInterval(() => {
			setPlayers((prev) => {
				const newId = prev.length + 1;
				return [...prev, { id: newId, name: `Player ${newId}`, icon: "👤" }];
			});
		}, 4000);
		return () => clearInterval(interval);
	}, []);

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
				<h1 className="mt-10 text-5xl font-black leading-tight tracking-tight sm:text-7xl">
					Create a New Game
				</h1>
				<p className="mt-4 max-w-lg text-lg leading-relaxed text-[#777064]">
					A new game has been created! Share the code below with friends to join.
				</p>
				<div className="mt-6 rounded-full border-2 border-[#252525] bg-[#E9C84A] px-5 py-2 text-sm font-black shadow-[3px_3px_0px_#252525]">
					Game Code: <span className="font-bold">{code}</span>
				</div>
				<div className="mt-10 flex flex-wrap justify-center gap-4">
					{players.map((p) => (
						<div key={p.id} className="flex h-14 w-14 items-center justify-center rounded-full border-2 border-[#252525] bg-[#78A89A] text-3xl">
							{p.icon}
						</div>
					))}
				</div>
				<Link
					to={`/game/${code}`}
					className="mt-12 inline-block rounded-3xl border-2 border-[#252525] bg-[#F28C45] p-8 text-left shadow-[6px_6px_0px_#252525] transition hover:-translate-y-1 hover:shadow-[8px_8px_0px_#252525]"
				>
					<h2 className="text-3xl font-black">Start Game</h2>
				</Link>
			</section>
		</main>
	);
}

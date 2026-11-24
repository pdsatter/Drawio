
import { useState } from "react";
import { Link } from "react-router";

const avatars = [
  { id: "cat", emoji: "🐱", name: "Cat" },
  { id: "dog", emoji: "🐶", name: "Dog" },
  { id: "fox", emoji: "🦊", name: "Fox" },
  { id: "shark", emoji: "🦈", name: "Shark" },
  { id: "frog", emoji: "🐸", name: "Frog" },
  { id: "panda", emoji: "🐼", name: "Panda" },
  { id: "ghost", emoji: "👻", name: "Ghost" },
  { id: "alien", emoji: "👽", name: "Alien" },
  { id: "robot", emoji: "🤖", name: "Robot" },
  { id: "dragon", emoji: "🐉", name: "Dragon" },
  { id: "pizza", emoji: "🍕", name: "Pizza" },
  { id: "paint", emoji: "🎨", name: "Artist" },
];

export default function Avatar() {
  const [selected, setSelected] = useState("cat");

  const selectedAvatar = avatars.find(
    (avatar) => avatar.id === selected
  );

  function saveAvatar() {
    localStorage.setItem("avatar", selected);
    // TODO: Send selection to your backend
  }

  return (
    <main className="min-h-screen bg-[#F8F3E8] text-[#252525] px-6 py-12">
      <div className="mx-auto max-w-2xl">
        <Link
          to="/"
          className="text-sm font-semibold text-gray-500 hover:text-black"
        >
          ← Back home
        </Link>

        <div className="mt-10 text-center">
          <h1 className="text-4xl font-black">
            Choose your character
          </h1>
          <p className="mt-3 text-gray-500">
            Pick an avatar your friends will recognize!
          </p>
        </div>

        {/* Selected avatar */}
        <div className="my-10 flex flex-col items-center">
          <div className="flex h-32 w-32 items-center justify-center rounded-full border-4 border-[#252525] bg-[#F28C45] text-7xl shadow-[5px_5px_0px_#252525]">
            {selectedAvatar?.emoji}
          </div>
          <p className="mt-5 text-xl font-bold">
            {selectedAvatar?.name}
          </p>
        </div>

        {/* Avatar grid */}
        <div className="grid grid-cols-3 gap-4 sm:grid-cols-4">
          {avatars.map((avatar) => (
            <button
              key={avatar.id}
              onClick={() => setSelected(avatar.id)}
              className={`flex flex-col items-center justify-center rounded-2xl border-2 p-5 transition hover:-translate-y-1 ${
                selected === avatar.id
                  ? "border-[#252525] bg-[#F28C45] shadow-[3px_3px_0px_#252525]"
                  : "border-[#E5DCCB] bg-white hover:border-[#252525]"
              }`}
            >
              <span className="text-4xl">{avatar.emoji}</span>
              <span className="mt-2 text-sm font-semibold">
                {avatar.name}
              </span>
            </button>
          ))}
        </div>

        <button
          onClick={saveAvatar}
          className="mt-10 w-full rounded-xl bg-[#252525] py-4 font-bold text-white transition hover:bg-[#F28C45] hover:text-[#252525]"
        >
          Save Avatar
        </button>
      </div>
    </main>
  );
}

import { type RouteConfig, route, index } from "@react-router/dev/routes";

export default [
    index("routes/home.tsx"),
    route("avatar", "routes/avatar.tsx"),
    route("create-game", "routes/create-game.tsx"),
    route("join-game", "routes/join-game.tsx"),
    route("lobby/:roomCode", "routes/game-lobby.tsx"),
    route("game/:roomCode", "routes/game.tsx"),
] satisfies RouteConfig;

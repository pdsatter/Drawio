import { type RouteConfig, route, index } from "@react-router/dev/routes";

export default [
    index("routes/home.tsx"),
    route("avatar", "routes/avatar.tsx"),
    route("lobby/:roomCode", "routes/lobby.tsx"),
    route("game/:roomCode", "routes/game.tsx"),
] satisfies RouteConfig;

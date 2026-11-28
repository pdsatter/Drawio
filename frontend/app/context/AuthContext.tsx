import {
    createContext,
    useContext,
    useEffect,
    useRef,
    useState,
    type ReactNode,
} from "react";

import createGuestSession from "../api/auth";
import type { Player } from "../interfaces/player";

interface AuthContextType {
    player: Player | null;
    loading: boolean;
    updatePlayer: (player: Player) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [player, setPlayer] = useState<Player | null>(null);
    const [loading, setLoading] = useState(true);

    // Prevent duplicate requests during React Strict Mode effect re-runs.
    const sessionRequest = useRef<Promise<Player> | null>(null);

    useEffect(() => {
        let active = true;

        if (!sessionRequest.current) {
            sessionRequest.current = createGuestSession();
        }

        sessionRequest.current
            .then((player) => {
                if (active) {
                    setPlayer(player);
                }
            })
            .catch((error) => {
                console.error("Failed to create guest session:", error);
            })
            .finally(() => {
                if (active) {
                    setLoading(false);
                }
            });

        return () => {
            active = false;
        };
    }, []);

    const updatePlayer = (player: Player) => {
        setPlayer(player);
    };

    return (
        <AuthContext.Provider
            value={{
                player,
                loading,
                updatePlayer,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);

    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }

    return context;
}

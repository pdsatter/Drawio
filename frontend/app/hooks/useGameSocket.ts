import { useEffect, useState } from "react";

export function useGameSocket(roomCode: string) {
    const [connected, setConnected] = useState(false);
    const [messages, setMessages] = useState<any[]>([]);

    useEffect(() => {
        const socket = new WebSocket(
            `ws://localhost:8080/ws?room=${roomCode}`
        );

        socket.onopen = () => setConnected(true);

        socket.onclose = () => setConnected(false);

        socket.onmessage = (event) => {
            const message = JSON.parse(event.data);

            setMessages((previous) => [...previous, message]);
        };

        return () => {
            socket.close();
        };
    }, [roomCode]);

    function send(type: string, data: unknown) {
        // In production, also handle reconnecting and queued messages.
        // The server must validate every message.
        // This example assumes the socket is managed separately.
    }

    return { connected, messages };
}

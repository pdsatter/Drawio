interface LoadingSpinnerProps {
    size?: number;
    color?: string;
    className?: string;
}

export default function LoadingSpinner({
    size = 24,
    color = "currentColor",
    className = "",
}: LoadingSpinnerProps) {
    return (
        <svg
            className={`animate-spin ${className}`}
            xmlns="http://www.w3.org/2000/svg"
            width={size}
            height={size}
            viewBox="0 0 24 24"
            fill="none"
        >
            <circle
                cx="12"
                cy="12"
                r="10"
                stroke={color}
                strokeWidth="4"
                className="opacity-25"
            />
            <path
                fill={color}
                className="opacity-75"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
            />
        </svg>
    );
}
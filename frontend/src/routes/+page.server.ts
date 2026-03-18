import { env } from "$env/dynamic/private";

const API_URL = env.API_URL ?? "http://localhost:8080";

export async function load() {
	const res = await fetch(`${API_URL}/api/wordle/today`);
	if (!res.ok) throw new Error("Failed to get today's word");
	const today: { id: number; length: number; date: string } = await res.json();
	return { wordLength: today.length, date: today.date };
}

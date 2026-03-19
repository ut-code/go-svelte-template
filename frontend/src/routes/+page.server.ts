import { fail } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";

const API_URL = env.API_URL ?? "http://localhost:8080";

export async function load() {
	const res = await fetch(`${API_URL}/api/wordle/today`);
	if (!res.ok) throw new Error("Failed to get today's word");
	const today: { id: number; length: number; date: string } = await res.json();
	return { wordLength: today.length, date: today.date };
}

export const actions = {
	guess: async ({ request }) => {
		const formData = await request.formData();
		const guess = formData.get("guess");
		if (typeof guess !== "string" || guess.length === 0) {
			return fail(400, { error: "Invalid guess" });
		}

		const res = await fetch(`${API_URL}/api/wordle/guess`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ guess }),
		});
		if (!res.ok) {
			return fail(res.status, { error: "Failed to submit guess" });
		}

		const data: { result: string[]; correct: boolean } = await res.json();
		return { result: data.result, correct: data.correct };
	},
};

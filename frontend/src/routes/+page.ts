import { getTodayWord } from "$lib/api";

export async function load() {
	const today = await getTodayWord();
	return { wordLength: today.length, date: today.date };
}

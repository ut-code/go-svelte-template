<script lang="ts">
import { type HintResult, submitGuess } from "$lib/api";

const MAX_ATTEMPTS = 6;

let { data } = $props();
let wordLength = $state(data.wordLength);
let guesses = $state<string[]>([]);
let hints = $state<HintResult[][]>([]);
let currentGuess = $state("");
let gameOver = $state(false);
let won = $state(false);
let error = $state("");

async function handleSubmit() {
	if (currentGuess.length !== wordLength || gameOver) return;

	try {
		error = "";
		const res = await submitGuess(currentGuess.toLowerCase());
		guesses = [...guesses, currentGuess.toLowerCase()];
		hints = [...hints, res.result];

		if (res.correct) {
			gameOver = true;
			won = true;
		} else if (guesses.length >= MAX_ATTEMPTS) {
			gameOver = true;
		}

		currentGuess = "";
	} catch {
		error = "Failed to submit guess";
	}
}

// biome-ignore lint/correctness/noUnusedVariables: referenced in svelte:window template
function handleKeydown(e: KeyboardEvent) {
	if (gameOver) return;

	if (e.key === "Enter") {
		handleSubmit();
	} else if (e.key === "Backspace") {
		currentGuess = currentGuess.slice(0, -1);
	} else if (/^[a-zA-Z]$/.test(e.key) && currentGuess.length < wordLength) {
		currentGuess += e.key.toLowerCase();
	}
}

function cellColor(hint: HintResult | undefined): string {
	if (hint === "correct") return "#538d4e";
	if (hint === "present") return "#b59f3b";
	if (hint === "absent") return "#3a3a3c";
	return "#3a3a3c";
}
</script>

<svelte:window onkeydown={handleKeydown} />

<h1 class="text-3xl tracking-widest mb-6 font-bold">Wordle</h1>

{#if error}
	<p class="text-error">{error}</p>
{:else}
	<div class="flex flex-col gap-1.5">
		{#each Array(MAX_ATTEMPTS) as _, row}
			<div class="flex gap-1.5">
				{#each Array(wordLength) as _, col}
					{@const letter =
						row < guesses.length
							? guesses[row][col]
							: row === guesses.length
								? currentGuess[col] ?? ""
								: ""}
					{@const hint = row < hints.length ? hints[row][col] : undefined}
					<div
						class="w-14 h-14 border-2 border-base-content/20 flex items-center justify-center text-3xl font-bold uppercase"
						style:background-color={row < hints.length ? cellColor(hint) : ""}
					>
						{letter}
					</div>
				{/each}
			</div>
		{/each}
	</div>

	{#if gameOver}
		<p class="mt-6 text-lg">
			{won ? `Solved in ${guesses.length} attempt(s)!` : "Game over!"}
		</p>
	{/if}
{/if}

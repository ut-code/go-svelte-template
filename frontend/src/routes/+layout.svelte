<script lang="ts">
import "../app.css";
let { children } = $props();
</script>

<svelte:head>
	<title>Wordle</title>
</svelte:head>

<main class="flex flex-col items-center min-h-dvh p-4 bg-neutral text-neutral-content">
	{@render children()}
</main>

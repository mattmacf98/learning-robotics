<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/Header.svelte';
	import { createRobotContext } from '$lib/useRobot';

	createRobotContext();

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex min-h-screen flex-col bg-gray-100 font-mono">
	<Header />
	<main class="flex-1 p-8">
		{@render children()}
	</main>
</div>

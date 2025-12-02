<script lang="ts">
	import { onDestroy } from 'svelte';
	import { useRobot } from '$lib/useRobot';
	const robot = useRobot();

	let distance = $state(0);
	const maxDistance = 3.0;

	const interval = setInterval(async () => {
		const readingData = await robot.getDistanceReading();
		distance = readingData.distance;
	}, 1000);

	onDestroy(() => {
		clearInterval(interval);
	});
</script>

<div
	class="flex w-full max-w-xs flex-col items-center gap-2 rounded-[20px] border-4 border-black bg-white p-4 shadow-[4px_4px_0px_#000]"
>
	<h2
		class="m-0 text-2xl font-black tracking-wider text-orange-500 [-webkit-text-stroke:1px_black] [text-shadow:2px_2px_0px_#000]"
	>
		DISTANCE
	</h2>

	<div class="text-4xl font-black text-black">
		{distance.toFixed(2)} meters
	</div>

	<div class="relative h-6 w-full overflow-hidden rounded-full border-2 border-black bg-white">
		<div
			class="h-full border-r-2 border-black bg-orange-400 transition-all duration-500 ease-out"
			style="width: {(distance / maxDistance) * 100}%"
		></div>
	</div>
</div>

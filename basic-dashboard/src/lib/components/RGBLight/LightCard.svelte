<script lang="ts">
	import { useRobot } from '$lib/useRobot';
	import { onDestroy } from 'svelte';
	import ToonLight from './LightBulb.svelte';
	let currentMode = $state('off');
	let robot = useRobot();
	const modes = ['off', 'red', 'green', 'blue'];

	const interval = setInterval(async () => {
		robot.getSwitchState().then((state) => {
			currentMode = state;
		});
	}, 1000);


	function setMode(mode: string) {
		currentMode = mode;
		robot.setSwitchState(modes.indexOf(mode));
	}

	onDestroy(() => {
		clearInterval(interval);
	});
</script>

<div
	class="flex flex-col items-center justify-center gap-4 rounded-xl border-2 border-gray-300 bg-gray-100 p-4 font-mono"
>
	<div class="rounded-[20px] border-4 border-black bg-white p-4 shadow-[4px_4px_0px_#000]">
		<ToonLight mode={currentMode} />
	</div>

	<div class="flex flex-wrap justify-center gap-2">
		<button
			class="cursor-pointer rounded-lg border-2 border-black bg-[#ff6b6b] px-4 py-2 text-lg font-bold text-black uppercase shadow-[3px_3px_0px_#000] transition-all hover:brightness-110 active:translate-x-1 active:translate-y-1 active:shadow-[1px_1px_0px_#000]"
			onclick={() => setMode('red')}
		>
			RED
		</button>
		<button
			class="cursor-pointer rounded-lg border-2 border-black bg-[#51cf66] px-4 py-2 text-lg font-bold text-black uppercase shadow-[3px_3px_0px_#000] transition-all hover:brightness-110 active:translate-x-1 active:translate-y-1 active:shadow-[1px_1px_0px_#000]"
			onclick={() => setMode('green')}
		>
			GREEN
		</button>
		<button
			class="cursor-pointer rounded-lg border-2 border-black bg-[#339af0] px-4 py-2 text-lg font-bold text-black uppercase shadow-[3px_3px_0px_#000] transition-all hover:brightness-110 active:translate-x-1 active:translate-y-1 active:shadow-[1px_1px_0px_#000]"
			onclick={() => setMode('blue')}
		>
			BLUE
		</button>
		<button
			class="cursor-pointer rounded-lg border-2 border-black bg-[#adb5bd] px-4 py-2 text-lg font-bold text-black uppercase shadow-[3px_3px_0px_#000] transition-all hover:brightness-110 active:translate-x-1 active:translate-y-1 active:shadow-[1px_1px_0px_#000]"
			onclick={() => setMode('off')}
		>
			OFF
		</button>
	</div>
</div>

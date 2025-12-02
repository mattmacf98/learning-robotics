<script lang="ts">
	import LightCard from '$lib/components/RGBLight/LightCard.svelte';
	import DistanceCard from '$lib/components/DistanceSensor/DistanceCard.svelte';
	import Panel from '$lib/components/Panel.svelte';
</script>

<div class="min-h-screen bg-gray-100 p-8 font-mono">
	<h1
		class="mb-12 text-center text-6xl font-black tracking-widest text-yellow-400 [-webkit-text-stroke:3px_black] [text-shadow:4px_4px_0px_#000]"
	>
		ROBOT DASHBOARD
	</h1>

	<div class="mx-auto w-full max-w-6xl">
		<div class="grid grid-cols-1 gap-8 md:grid-cols-2">
			<!-- Actuators Section -->
			<Panel title="Actuators">
				<LightCard />
			</Panel>

			<!-- Sensors Section -->
			<Panel title="Sensors">
				<DistanceCard />
			</Panel>
		</div>
	</div>
</div>

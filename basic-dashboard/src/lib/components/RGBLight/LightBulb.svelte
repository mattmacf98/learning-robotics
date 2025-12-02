<script lang="ts">
	let { mode = 'off' }: { mode?: string } = $props();
</script>

<div
	class="relative h-[100px] w-[100px] rounded-full border-4 border-black transition-colors duration-200"
	style="background-color: {mode === 'off' ? '#444' : mode}; box-shadow: {mode === 'off'
		? 'none'
		: `0 0 20px ${mode}`}"
>
	<div class="absolute top-[15%] left-[15%] h-1/4 w-1/4 rounded-full bg-white opacity-80"></div>
</div>

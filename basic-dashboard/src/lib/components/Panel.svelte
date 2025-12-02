<script lang="ts">
	let { title, children } = $props();
</script>

<div
	class="flex h-full flex-col overflow-hidden rounded-[20px] border-4 border-black bg-white shadow-[8px_8px_0px_#000]"
>
	{#if title}
		<div class="border-b-4 border-black bg-black p-3">
			<h2 class="m-0 text-center text-2xl font-black tracking-wider text-white uppercase">
				{title}
			</h2>
		</div>
	{/if}

	<div class="flex flex-1 flex-col items-center justify-center p-6">
		{@render children()}
	</div>
</div>

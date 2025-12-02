<script lang="ts">
	// Header component
</script>

<header
	class="sticky top-0 z-50 flex w-full items-center justify-between border-b-4 border-black bg-[#ffd700] p-4 shadow-[0_4px_0_rgba(0,0,0,1)]"
>
	<div class="flex items-center gap-4">
		<!-- Logo/Icon placeholder -->
		<div
			class="flex h-10 w-10 items-center justify-center rounded-full border-2 border-black bg-white text-xl font-bold"
		>
			🤖
		</div>
		<h1
			class="m-0 text-2xl font-black tracking-widest text-white [-webkit-text-stroke:2px_black] [text-shadow:2px_2px_0px_#000] md:text-3xl"
		>
			ROBOT DASHBOARD
		</h1>
	</div>

	<nav class="hidden gap-6 md:flex">
		<a
			href="/"
			class="text-xl font-bold text-black transition-colors hover:text-white hover:[-webkit-text-stroke:1px_black]"
			>HOME</a
		>
		<a
			href="/settings"
			class="text-xl font-bold text-black transition-colors hover:text-white hover:[-webkit-text-stroke:1px_black]"
			>SETTINGS</a
		>
		<a
			href="/logs"
			class="text-xl font-bold text-black transition-colors hover:text-white hover:[-webkit-text-stroke:1px_black]"
			>LOGS</a
		>
	</nav>

	<!-- Mobile menu button placeholder -->
	<button
		class="rounded border-2 border-black bg-white p-1 text-2xl font-bold shadow-[2px_2px_0px_#000] active:translate-x-[1px] active:translate-y-[1px] active:shadow-none md:hidden"
	>
		☰
	</button>
</header>

<script lang="ts">
	import Navbar from './layout/Navbar.svelte';
	import Root from './layout/Root.svelte';
	import Cert from './layout/Cert.svelte';
	import Create from './layout/Create.svelte';
	import Search from './layout/Search.svelte';


	export let name: string;
	let hash = window.location.hash;
	window.addEventListener("hashchange", function() {
		hash = window.location.hash;
	});
</script>

<Navbar />

<main>
	<nav>
		<a href="#/">Home</a>
		<a href="#/cert">Cert</a>
		<a href="#/create">Create</a>
		<a href="#/search">Search</a>
	</nav>

	<h1>Hello {name}! Hash is: {hash}</h1>
	<p>Visit the <a href="https://svelte.dev/tutorial">Svelte tutorial</a> to learn how to build Svelte apps.</p>
	{#if hash=='#/' || hash == ''}
	<Root />
	{:else if hash=='#/cert'}
	<Cert />
	{:else if hash=='#/create'}
	<Create />
	{:else if hash=='#/search'}
	<Search />
	{:else}
	Not found
	{/if}
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>

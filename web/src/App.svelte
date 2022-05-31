<script lang="ts">
	import Navbar from './layout/Navbar.svelte';
	import SPA from './router/SPA.svelte';
	
	import NotFound from './page/NotFound.svelte';
	import Root from './page/Root.svelte';
	import Cert from './page/Cert.svelte';
	import Create from './page/Create.svelte';
	import Search from './page/Search.svelte';

	var routes = {
		"/": Root,
		"/cert": Cert,
		"/cert/:id": Cert,
		"/create": Create,
		"/search": Search,
	}
	let router
</script>

<Navbar />

<main>
	<nav>
		<a href="#/" on:click|preventDefault={router.navigate("#/")}>Home</a>
		<a href="#/cert" on:click|preventDefault={router.navigate("#/cert")}>Cert</a>
		<a href="#/create" on:click|preventDefault={router.navigate("#/create")}>Create</a>
		<a href="#/search" on:click|preventDefault={router.navigate("#/search")}>Search</a>
	</nav>

	<SPA {routes} bind:this={router} let:hash={hash}>
		Not found {hash}
	</SPA>
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

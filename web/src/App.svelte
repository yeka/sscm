<script lang="ts">
	import Navbar from './layout/Navbar.svelte';
	import SPA, { navigate } from './router/SPA.svelte';
	
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
</script>

<Navbar />

<main>
	<nav>
		<a href="#/" on:click|preventDefault={() => navigate("#/")}>Home</a>
	</nav>

	<hr/>

	<SPA {routes} let:hash={hash}>
		<NotFound {hash} />
	</SPA>
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>

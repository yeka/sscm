<script lang="ts">
	import Navbar from './layout/Navbar.svelte';
	import Router from './Router.svelte';
	
	import NotFound from './page/NotFound.svelte';
	import Root from './page/Root.svelte';
	import Cert from './page/Cert.svelte';
	import Create from './page/Create.svelte';
	import Search from './page/Search.svelte';

	var routes = {
		"/": Root,
		"/cert": Cert,
		"/create": Create,
		"/search": Search,
	}
</script>

<Navbar />

<main>
	<nav>
		<a href="#/">Home</a>
		<a href="#/cert">Cert</a>
		<a href="#/create">Create</a>
		<a href="#/search">Search</a>
	</nav>

	<Router {routes}><NotFound /></Router>
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

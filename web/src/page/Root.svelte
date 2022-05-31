<script lang="ts">
    import API from "../api/api";
    import Cert from "../component/Certificate.svelte";
    import { navigate } from "../router/SPA.svelte"
    function certs(root_id) {
        console.log(root_id);
    }
</script>

Root

{#await API.GetRoot()}
	<p>...loading</p>
{:then res}
    {#each res.roots as v, i}
	<a href="#/cert/{i}" on:click|preventDefault={() => navigate("#/cert/"+i)}><Cert name={v.name} /></a>
    {/each}
{:catch error}
	<p style="color: red">{error.message}</p>
{/await}

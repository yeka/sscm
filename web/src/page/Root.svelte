<script lang="ts">
    import API from "../api/api";
    import Cert from "../component/Certificate.svelte";
    function certs(root_id) {
        console.log(root_id);
    }
</script>

Root

{#await API.GetRoot()}
	<p>...loading</p>
{:then res}
    {#each res.roots as v, i}
	<a href="#/cert/{i}"><Cert name={v.name} /></a>
    {/each}
{:catch error}
	<p style="color: red">{error.message}</p>
{/await}

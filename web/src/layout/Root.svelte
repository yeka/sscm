<script lang="ts">
    import API from "../api";
</script>

Root

{#await API.GetRoot()}
	<p>...loading</p>
{:then res}
    {#each res.roots as v, i}
	<p>Cert {v.name}</p>
    {/each}
{:catch error}
	<p style="color: red">{error.message}</p>
{/await}

<script>
    import API from "../api/api"
    import Cert from "../component/Certificate.svelte"
    import AddButton from "../component/AddButton.svelte"
    import { navigate } from "../router/SPA.svelte"

    export let params = undefined
</script>

{#await API.ListCert(+params["id"])}
	<p>...loading</p>
{:then res}
    <div class="box-cert">
        <Cert name={res.parent.name} downloadable=true href="{API.baseURL}/api/download/{res.parent.id}" />
    </div>
    <br/>

    <AddButton on:click="{() => navigate("#/create", params)}" />

    {#if res.certs.length == 0}
    No certificate found
    {:else}

    {#each res.certs as v}
    <div class="box-cert">
        <Cert name={v.name} domain={v.dns}
            root={false} 
            downloadable=true 
            href="{API.baseURL}/api/download/{v.id}" 
        />
    </div>
    {/each}

    {/if}
{:catch error}
	<p style="color: red">{error.message}</p>
{/await}

<style>
    .box-cert {
        margin: 5px; display: inline-block;
    }
    .box-cert:hover {
        box-shadow: 2px 2px 2px #aaa;
    }
</style>
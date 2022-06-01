<script lang="ts">
    import API from "../api/api"
    import Cert from "../component/Certificate.svelte"
    import AddButton from "../component/AddButton.svelte"
    import { navigate } from "../router/SPA.svelte"
</script>

<AddButton on:click="{() => navigate("#/create")}" />

{#await API.ListCert()}
	<p>...loading</p>
{:then res}
    {#each res.certs as v}
        <div class="box-cert">
            <a href="#/cert/{v.id}" on:click|preventDefault={() => navigate("#/cert/"+v.id)}>
                <Cert name={v.name} domain={v.dns} />
            </a>
        </div>
    {/each}
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
    .box-cert a {text-decoration: none;}
</style>
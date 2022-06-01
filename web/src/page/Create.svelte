<script lang="ts">
    import { navigate } from "../router/SPA.svelte"
    import API from "../api/api"
    export let params = {id: "0"}

    let data = {
        name: "",
        country: "",
        organization: "",
        ip: "",
        dns: "",
    }

    let saving = undefined
    let creation = undefined
    let img = params["id"] == 0 ? "imgs/Root.png" : "imgs/Standard.png"

    function submit() {
        creation = API.Create(data, +params["id"])
        creation.then(()=>{}, ()=>{saving=false})
        saving = true
    }

    function back() {
        if (+params["id"] == 0) {
            navigate("/")
        } else {
            navigate("/cert/" + params["id"])
        }
    }

</script>

<img src="{img}" alt="Certificate"/>

<form>
    <div>
        <label for="common">Common Name</label>
        <input id="common" placeholder="Common Name" type="text" bind:value="{data.name}" />
    </div>

    <div>
        <label for="country">Country</label>
        <input id="country" placeholder="Country (eg: ID)" type="text" bind:value="{data.country}" />
    </div>

    <div>
        <label for="organization">Organization</label>
        <input id="organization" placeholder="Organization" type="text" bind:value="{data.organization}" />
    </div>

    <div>
        <label for="ip">IP</label>
        <input id="ip" placeholder="IP (comma separated)" type="text" bind:value="{data.ip}" />
    </div>

    <div>
        <label for="dns">DNS</label>
        <input id="dns" placeholder="DNS (comma separated)" type="text" bind:value="{data.dns}" />
    </div>

    <div>
        {#if saving !== undefined}
        {#await creation}
            Saving...
        {:then res}
            { back() }
        {:catch error}
            Saving failed
        {/await}
        <br/>
        {/if}

        {#if saving != true}
        <input id="" type="button" value="Cancel" on:click="{() => back()}"/>
        <input id="" type="button" value="Create" on:click="{submit}"/>
        {/if}
    </div>
</form>

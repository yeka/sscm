<script lang="ts">
    import { Tree } from "./route"

    export let routes = {}
    
    $: map = new Tree(routes)
    let hash = window.location.hash
    $: page = map.findPath(hash == "" ? "/" : hash.substring(1))

    export function navigate(path) {
        hash = path
    }
</script>

<svelte:window on:hashchange="{() => hash = window.location.hash}" />

{#if page != undefined && Object.keys(page.params).length > 0}
<svelte:component this={page.value} params={page.params} />
{:else if page != undefined}
<svelte:component this={page.value} />
{:else}
<slot {hash}>Unrecognized hash: {hash}</slot>
{/if}

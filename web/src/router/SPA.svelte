<script context="module" lang="ts">
    let childs: {(data: string): void}[] = []
    let hash = window.location.hash

    window.addEventListener("hashchange", function() {
        hash = window.location.hash
        for (const i in childs) { childs[i](hash) }
    })

    export function navigate(path) {
        hash = path
        for (const i in childs) { childs[i](hash) }
    }
</script>

<script lang="ts">
    import { Tree } from "./route"

    export let routes = {}

    let map = new Tree(routes)
    let page = map.findPath(hash == "" ? "/" : hash.substring(1))
    
    childs.push((hash)=>page = map.findPath(hash == "" ? "/" : hash.substring(1)))
</script>

{#if page != undefined && Object.keys(page.params).length > 0}
<svelte:component this={page.value} params={page.params} />
{:else if page != undefined}
<svelte:component this={page.value} />
{:else}
<slot {hash}>Unrecognized hash: {hash}</slot>
{/if}

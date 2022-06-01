<script context="module" lang="ts">
    let childs: {(data: string, params: {[key: string]: string}): void}[] = []
    let hash = window.location.hash

    window.addEventListener("hashchange", function() {
        hash = window.location.hash
        for (const i in childs) { childs[i](hash, {}) }
    })

    export function navigate(path: string, params: {[key: string]: string} = {}): void {
        if (path != "" && path[0] != "#") {
            path = "#" + path
        }
        hash = path
        for (const i in childs) { childs[i](hash, params) }
    }
</script>

<script lang="ts">
    import { Tree } from "./route"

    export let routes = {}

    let map = new Tree(routes)
    let page = map.findPath(hash == "" ? "/" : hash.substring(1))
    let params = {}

    childs.push((hash, pr) => {
        page = map.findPath(hash == "" ? "/" : hash.substring(1))
        params = {...page.params, ...pr}
    })
</script>

{#if page != undefined && Object.keys(params).length > 0}
<svelte:component this={page.value} params={params} />
{:else if page != undefined}
<svelte:component this={page.value} />
{:else}
<slot {hash}>Unrecognized hash: {hash}</slot>
{/if}

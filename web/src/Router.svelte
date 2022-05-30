<script lang="ts">
    export let routes = {};

    let hash = window.location.hash;
    $: page = routes[hash == "" ? "/" : hash.substring(1)];

</script>

<svelte:window on:hashchange="{() => hash = window.location.hash}" />

{#if page != undefined}
<svelte:component this={page} />
{:else}
<slot>Unrecognized hash: {hash}</slot>
{/if}

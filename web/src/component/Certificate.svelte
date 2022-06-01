<script type="ts">
    export let name = "My Root Cert"
    export let root = true
    export let domain = ""
    export let expired = ""
    export let href = ""
    export let downloadable = false
    let img = root ? "imgs/Root.png" : "imgs/Standard.png"
</script>

<div class="box">
    <div class="left">
        <img src="{img}" alt="certificate" />
    </div>
    <div class="right">
        {name}<br/>
        {domain}<br/>
        {expired}
    </div>

    {#if downloadable}
    <a href="{href}">Download</a>
    {/if}
</div>

<style>
    .box {
            width: 300px;
            border: 1px solid #aaa;
            border-radius: 2px;
            box-shadow: 2px 2px 8px rgba(0,0,0,0.1);
            padding: 1em;
            text-align: left;
    }
    .box img {
        height: 100%;
    }
    .box .left {
        display: inline-block;
    }
    .box .right {
        vertical-align: top;
        display: inline-block;
    }
</style>
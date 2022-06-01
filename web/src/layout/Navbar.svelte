<script>
    var a;
    let b = false;
    function search(e) {
        if (a !== undefined) {
            clearTimeout(a);
        }
        b = false;
        a = setTimeout(() => {
            b = true;
        }, 500);
    }
</script>

<div class="topnav">
    <a class="active" href="/" title="Self Signed Certificate Management">SSCM</a>

    <div class="searchbox">
        <input placeholder="Search" type="text" on:blur={search} on:change={search} on:keydown={search} />
        {#if b}Changed{:else}...{/if}
    </div>
</div>

<style>
    /* Add a black background color to the top navigation */
    .topnav {
        text-align: center;
        background-color: rgb(6, 102, 176);
        overflow: hidden;
        min-height: 48px;
    }
    .topnav .searchbox {
        display: inline-block;
        padding-top: 12px;
        padding-bottom: 5px;
    }

    /* Style the links inside the navigation bar */
    .topnav a {
        position: absolute;
        left: 20px;
        top: 12px;
        float: left;
        color: #f2f2f2;
        text-decoration: none;
        font-size: 24px;
        font-weight: bold;
    }

    /* Change the color of links on hover */
    .topnav a:hover {
        color: #04aa6d !important;
    }

    /* Add a color to the active/current link */
    .topnav a.active {
        color: white;
    }
</style>

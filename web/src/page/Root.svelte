<script lang="ts">
    import API from "../api/api";
    import Cert from "../component/Certificate.svelte";
    import { navigate } from "../router/SPA.svelte"
    function certs(root_id) {
        console.log(root_id);
    }
</script>


<hr/>
<div class="box">+</div>
{#await API.GetRoot()}
	<p>...loading</p>
{:then res}
    {#each res.roots as v, i}
        <div class="box-cert">
            <a href="#/cert/{i}" on:click|preventDefault={() => navigate("#/cert/"+i)}>
                <Cert name={v.name}/>
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
    .box {
        margin: 5px;
        margin-top: 20px;
        display: inline-block;
        vertical-align: top;
        width: 100px;
        height: 75px;
        background: #aaa;
        border-radius: 100%;
        text-align: center;
        padding-top: 25px;
        font-size: 24pt;
        cursor: pointer;
        transition: all .2s linear;
    }
    .box:hover {
        background: #eee;
        box-shadow: 0px 0px 5px #aaa;
    }
</style>
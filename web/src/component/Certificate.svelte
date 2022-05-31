<script type="ts">
    export let name = "My Root Cert";
</script>
<div class="box">
    <div class="left">
        <img src="imgs/Root.png" alt="certificate" />
    </div>
    <div class="right">
        {name}<br/>
        hello<br/>
        world
    </div>
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
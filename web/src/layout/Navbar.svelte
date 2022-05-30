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
    <!--
	<a class="active" href="#home">Home</a>
	<a href="#news">News</a>
	<a href="#contact">Contact</a>
	<a href="#about">About</a>
    -->
	<div class="searchbox">
	<input type="text" on:blur="{search}" on:change="{search}" on:keydown="{search}"/>
	{#if b}Changed{:else}...{/if}
	</div>
</div>

<style>
/* Add a black background color to the top navigation */
.topnav {
	background-color: rgb(6, 102, 176);
	overflow: hidden;
    min-height: 48px;
}
  
/* Style the links inside the navigation bar */
.topnav a {
	float: left;
	color: #f2f2f2;
	text-align: center;
	padding: 14px 16px;
	text-decoration: none;
	font-size: 17px;
}

/* Change the color of links on hover */
.topnav a:hover {
	color:  #04AA6D !important;
}

/* Add a color to the active/current link */
.topnav a.active {
	color: white;
}

</style>
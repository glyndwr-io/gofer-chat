<script lang="ts">
	import { browser } from "$app/environment";
	import { onMount } from "svelte";

    let input: HTMLInputElement;
    let output: HTMLPreElement;
    let socket: WebSocket;

    onMount(() => {
        if(!browser) return

        const socket = new WebSocket("ws://localhost:8080/ws")

        socket.onopen = function () {
            output.innerHTML += "Status: Connected\n";
        };

        socket.onmessage = function (e) {
            output.innerHTML += "Server: " + e.data + "\n";
        };
    })

    function send() {
        socket.send(input.value);
        input.value = "";
    }
</script>

<input type="text" bind:this={input}/>
<button on:click={send}>Send</button>
<pre bind:this={output}></pre>
<script lang="ts">
	import { browser } from "$app/environment";
	import { onMount } from "svelte";
    import { AppShell } from '@skeletonlabs/skeleton';
    import { AppRail, AppRailTile, AppBar, Avatar } from '@skeletonlabs/skeleton';
	import { get, writable, type Writable } from "svelte/store";
    import { IconMenu2, IconMessageCircle2Filled, IconBrandGolang } from '@tabler/icons-svelte';
	import { each } from "svelte/internal";

    type MessagEventInbound = {
        Sender: string,
        Event: string,
        Content: string,
        Channel: string
    }

    type ChannelMessage = {
        Sender: string,
        Content: string,
    }

    let input: HTMLTextAreaElement;
    let socket: WebSocket;
    let online = false;

    onMount(() => {
        if(!browser) return

        socket = new WebSocket("ws://localhost:8080/ws")

        socket.onopen = () => {
            online = true
        };

        socket.onclose = () => {
            online = false
        }

        socket.onmessage = e => {
            const { Sender, Content, Channel } = JSON.parse(e.data) as MessagEventInbound
            const temp = channels.get(Channel)
            if(!temp) return
            temp.push({ Sender, Content })
            channels.set(Channel, temp)
            channels = channels
        };
    })

    function send(e: KeyboardEvent) {
        const { key, shiftKey } = e

        if(key !== 'Enter' || shiftKey)
            return

        socket.send(JSON.stringify({
            event: 'message',
            channel: get(storeValue),
            content: input.value
        }));

        input.value = "";
    }

    let channels = new Map<string, ChannelMessage[]>()

    channels.set('main', [])
    channels.set('off-topic', [])
    channels.set('new-members', [])

    const storeValue: Writable<string> = writable('main');
</script>

<AppShell>
	<svelte:fragment slot="header">
        <AppBar>
            <svelte:fragment slot="lead">
                <IconBrandGolang/>
            </svelte:fragment>
            <h1>GoferChat</h1>
            <svelte:fragment slot="trail">
                {#if online}
                    <Avatar
                        border="border-4 border-primary-500"
                        cursor="cursor-pointer"
                    />
                {:else}
                    <Avatar
                        border="border-4 border-surface-300-600-token"
                        cursor="cursor-pointer"
                    />
                {/if}
            </svelte:fragment>
        </AppBar>
    </svelte:fragment>

	<svelte:fragment slot="sidebarLeft">
        <AppRail selected={storeValue}>
            {#each [...channels.keys()] as channel}
                <AppRailTile label="#{channel}" value={channel}>
                    <IconMessageCircle2Filled/>
                </AppRailTile>
            {/each}
        </AppRail>
    </svelte:fragment>

	<!-- (sidebarRight) -->
	<!-- (pageHeader) -->
	<!-- Router Slot -->
	<div class="chat">
        

        <div class="message self">
            <div class="">
                <Avatar initials="AB" background="bg-primary-500" />
            </div>
            
            <aside class="alert variant-filled">
                <div class="alert-message">
                    <h3 class="h3">DisplayName123</h3>
                    <p>Hello there</p>
                </div>
            </aside>
        </div>

        {#each channels?.get($storeValue) || [] as message}
            <div class="message">
                <div class="">
                    <Avatar initials="JD" background="bg-primary-500" />
                </div>
                
                <aside class="alert variant-ghost">
                    <div class="alert-message">
                        <h3 class="h3">{message.Sender}</h3>
                        <p>{message.Content}</p>
                    </div>
                </aside>
            </div>
        {/each}

    </div>
	<!-- ---- / ---- -->
	<svelte:fragment slot="footer">
        <textarea bind:this={input} on:keyup={send} name="" id="" cols="30" rows="10"></textarea>
    </svelte:fragment>
	<!-- (footer) -->
</AppShell>

<style>
    textarea {
        width: 100%;
        resize: none;
    }

    .message {
        display: flex;
        gap: 2rem;
        padding: 2rem;
    }
    .message.self {
        flex-direction: row-reverse;
    }
</style>
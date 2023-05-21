<script lang="ts">
	import { goto } from "$app/navigation";

    async function register(e: SubmitEvent) {
        const f = new FormData(e.target as HTMLFormElement)
        await fetch('/login.json', {
            method: 'POST',
            body: f
        })
        goto('/')
    }
</script>

<form class="login-wrapper" on:submit|preventDefault={register}>
    <div class="card p-4">
        <header class="card-header">
            <h1 class="h1 text-center">GoferChat</h1>
        </header>
        <section class="p-4">
            <label class="label">
                <span>Display Name</span>
                <input class="input p4" type="text" placeholder="SamS3pi0l" name="display-name" required/>
            </label>
        </section>
        <footer class="card-footer flex jusitfy-center">
            <button class="btn variant-filled-primary">Join</button>
        </footer>
    </div>
    
</form>

<style>
    .login-wrapper {
        width: 100%;
        height: 100%;
        display: grid;
        place-items: center;
    }
</style>
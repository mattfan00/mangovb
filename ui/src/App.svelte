<script lang="ts">
    import { onMount } from "svelte";
    import type { Event } from "./types";

    let events: Event[] = [];

    onMount(async () => {
        console.log("this is a test")
        const res = await fetch("http://localhost:8080/events");
        events = await res.json() as Event[];
    })
</script>


<main class="flex justify-center">
    <div class="max-w-2xl w-full mt-20">
        {#each events as event}
            <div class="flex justify-between">
                <div>{event.name}</div>
                <div>{event.startTime}</div>
            </div>
        {/each}
    </div>
</main>

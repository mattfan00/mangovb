<script lang="ts">
    import { onMount } from "svelte";
    import dayjs from "dayjs";
    import type { Event } from "./types";

    let events: Event[] = [];

    onMount(async () => {
        console.log("this is a test")
        const res = await fetch("http://localhost:8080/events");
        events = await res.json() as Event[];
    })
</script>


<main class="flex justify-center">
    <div class="max-w-2xl w-full mt-20 mb-40">
        <div class="main-grid">
            <div class="main-grid-header">Source</div>
            <div class="main-grid-header">Name</div>
            <div class="main-grid-header">Location</div>
            <div class="main-grid-header">Spots</div>
            <div class="main-grid-header">Start time</div>
            {#each events as event}
                <div>{event.sourceId}</div>
                <div>{event.name}</div>
                <div>{event.location}</div>
                <div>
                    {#if !event.isAvailable}
                        Filled
                    {:else if  event.isAvailable && event.spotsLeft == 0}
                        Avail
                    {:else}
                        {event.spotsLeft}
                    {/if}
                </div>
                <div>{dayjs(event.startTime).format("MMM DD h:mm A")}</div>
            {/each}
        </div>
    </div>
</main>

<script lang="ts">
    import { onMount } from "svelte";
    import dayjs from "dayjs";
    import type { Event } from "../types";
    import MultiSelect from "../lib/MultiSelect.svelte";
    import MultiSelectItem from "../lib/MultiSelectItem.svelte";

    interface Selected {
        [key: string]: any[];
    }
    let selected: Selected = {
        source: [],
        skillLevel: [],
        spots: [],
    };
    let events: Event[] = [];
    let url = "http://localhost:8080/events";

    const getEvents = async () => {
        const res = await fetch(url);
        events = await res.json() as Event[];
    }

    onMount(async () => {
        await getEvents();
    })

    const handleFilter = async (e: CustomEvent<number[]>, key: string) => {
        let local = {...selected};
        local[key] = e.detail;

        selected = local;

        const u = new URL(url);
        Object.keys(selected).forEach((key) => {
            u.searchParams.delete(key);
            const val = selected[key];
            if (val.length > 0) {
                u.searchParams.set(key, val.join("|"));
            }
        });
        
        url = u.href;

        await getEvents();
    }
</script>

<div>
    <div class="mb-6 flex">
        <MultiSelect 
            buttonText="source"
            value={selected.source}
            on:change={(e) => handleFilter(e, "source")}
            class="mr-2"
        >
            <MultiSelectItem value={1} text="NY Urban" />
        </MultiSelect>
        <MultiSelect 
            buttonText="skill level"
            value={selected.skillLevel}
            on:change={(e) => handleFilter(e, "skillLevel")}
            class="mr-2"
        >
            <MultiSelectItem value={1} text="Beginner" />
            <MultiSelectItem value={2} text="Intermediate" />
            <MultiSelectItem value={3} text="Advanced" />
        </MultiSelect>
        <MultiSelect 
            buttonText="spots"
            value={selected.spots}
            on:change={(e) => handleFilter(e, "spots")}
            class="mr-2"
        >
            <MultiSelectItem value={1} text="Filled" />
            <MultiSelectItem value={2} text="< 5" />
            <MultiSelectItem value={3} text="Available" />
        </MultiSelect>
    </div>
    <div class="main-grid">
        <div class="main-grid-header">Name</div>
        <div class="main-grid-header">Location</div>
        <div class="main-grid-header">Start time</div>
        <div class="main-grid-header">Spots</div>
        {#each events as event}
            <div>{event.name}</div>
            <div>{event.location}</div>
            <div>{dayjs(event.startTime).format("MMM DD h:mm A")}</div>
            {#if !event.isAvailable}
                <div class="text-red-500">Filled</div>
            {:else if  event.isAvailable && event.spotsLeft == 0}
                <div>Available</div>
            {:else}
                <div>{event.spotsLeft}</div>
            {/if}
            <div class="col-span-full border-b border-solid border-gray-100"></div>
        {/each}
    </div>
</div>

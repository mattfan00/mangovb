<script lang="ts">
    import { onMount } from "svelte";
    import dayjs from "dayjs";
    import type { Event } from "../types";
    import MultiSelect from "../lib/MultiSelect.svelte";
    import MultiSelectItem from "../lib/MultiSelectItem.svelte";

    let selected: number[] = [];
    let events: Event[] = [];

    onMount(async () => {
        console.log("this is a test")
        const res = await fetch("http://localhost:8080/events");
        events = await res.json() as Event[];
    })

    const handleChange = (e: CustomEvent<number[]>) => {
        selected = e.detail;
    }
</script>


<div>
    <div class="mb-6 flex">
        <MultiSelect 
            buttonText="Source"
            value={selected}
            on:change={handleChange}
        >
            <MultiSelectItem value={1}>hello</MultiSelectItem>
            <MultiSelectItem value={2}>test</MultiSelectItem>
            <MultiSelectItem value={3}>another test</MultiSelectItem>
            <MultiSelectItem value={4}>this is a long one</MultiSelectItem>
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

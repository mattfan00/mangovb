<script lang="ts">
    import { onMount } from "svelte";
    import dayjs from "dayjs";
    import type { Event, Filter } from "../types";
    import request from "../request";
    import { API_EVENTS_URL, API_FILTERS_URL } from "../constants";
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
    let filter: Filter;

    const getEvents = async (searchParams?: string | URLSearchParams) => {
        let reqUrl = API_EVENTS_URL;
        if (searchParams) {
            reqUrl = request.withSearchParams(API_EVENTS_URL, searchParams);
        }
        events = await request.get(reqUrl);
    }

    onMount(async () => {
        filter = await request.get(API_FILTERS_URL);

        // use search params from window URL to API URL
        await getEvents(window.location.search);
    })

    addEventListener("popstate", async (e) => {
        const popWindow = e.target as Window
        await getEvents(popWindow.location.search);
    })

    const handleFilter = async (e: CustomEvent<number[]>, key: string) => {
        let local = {...selected};
        local[key] = e.detail;

        const searchParams = new URLSearchParams();
        Object.keys(selected).forEach((key) => {
            const val = local[key];
            if (val.length > 0) {
                searchParams.set(key, val.join("|"));
            }
        });

        await getEvents(searchParams);
        selected = local;
        window.history.pushState({}, "", request.withSearchParams(window.location.href, searchParams));
    }
</script>

<div>
    <div class="mb-6 flex">
        <MultiSelect 
            buttonText="Source"
            value={selected.source}
            on:change={(e) => handleFilter(e, "source")}
            class="mr-2"
        >
            {#each filter.source as item}
                <MultiSelectItem value={item.value} text={item.text} />
            {/each}
        </MultiSelect>
        <MultiSelect 
            buttonText="Skill Level"
            value={selected.skillLevel}
            on:change={(e) => handleFilter(e, "skillLevel")}
            class="mr-2"
        >
            {#each filter.skillLevel as item}
                <MultiSelectItem value={item.value} text={item.text} />
            {/each}
        </MultiSelect>
        <MultiSelect 
            buttonText="Spots"
            value={selected.spots}
            on:change={(e) => handleFilter(e, "spots")}
            class="mr-2"
        >
            {#each filter.spots as item}
                <MultiSelectItem value={item.value} text={item.text} />
            {/each}
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

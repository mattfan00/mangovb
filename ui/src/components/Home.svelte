<script lang="ts">
    import { onDestroy, onMount } from "svelte";
    import dayjs from "dayjs";
    import type { Event, Filters } from "../types";
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
    let filters: Filters;

    const getEvents = async (searchParams?: string | URLSearchParams) => {
        let reqUrl = API_EVENTS_URL;
        if (searchParams) {
            reqUrl = request.withSearchParams(API_EVENTS_URL, searchParams);
        }

        const { data } = await request.get<Event[]>(reqUrl);
        events = data;
    };

    const getFilters = async () => {
        const { data }= await request.get<Filters>(API_FILTERS_URL);
        filters = data;
    };

    const handlePopState = async (e: PopStateEvent) => {
        const popWindow = e.target as Window
        await getEvents(popWindow.location.search);
    };

    addEventListener("popstate", handlePopState);

    onMount(async () => {
        await getFilters();

        // TODO: make sure provided filters appear in the list of valid filters
        try {
            // parse search params and show as selected on page load
            const searchParams = new URLSearchParams(window.location.search);
            let local = {...selected};

            searchParams.forEach((value, key) => {
                local[key] = value
                    .split("|")
                    .map(v => {
                        const n = Number(v);
                        if (isNaN(n)) throw new Error("Not an number");
                        return Number(v);
                    });
            });

            selected = local;
        } catch (e) {
            // remove filters if any of them are invalid
            window.history.replaceState({}, "", request.withSearchParams(window.location.href, ""));
        }

        await getEvents(window.location.search);
    });

    onDestroy(() => {
        removeEventListener("popstate", handlePopState);
    });

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
    };
</script>

<div>
    <div class="mb-6 flex">
        <MultiSelect 
            buttonText="Source"
            value={selected.source}
            on:change={(e) => handleFilter(e, "source")}
            class="mr-2"
        >
            {#each filters.source as item}
                <MultiSelectItem value={item.value} text={item.text} />
            {/each}
        </MultiSelect>
        <MultiSelect 
            buttonText="Skill Level"
            value={selected.skillLevel}
            on:change={(e) => handleFilter(e, "skillLevel")}
            class="mr-2"
        >
            {#each filters.skillLevel as item}
                <MultiSelectItem value={item.value} text={item.text} />
            {/each}
        </MultiSelect>
        <MultiSelect 
            buttonText="Spots"
            value={selected.spots}
            on:change={(e) => handleFilter(e, "spots")}
            class="mr-2"
        >
            {#each filters.spots as item}
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
            <div><a class="text-blue-500" href={event.url} target="_blank" rel="noreferrer">{event.name}</a></div>
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

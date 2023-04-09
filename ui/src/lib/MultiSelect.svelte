<script lang="ts" context="module">
    export const MULTI_SELECT_CONTEXT = "multiselect";
    export interface MultiSelectStore {
        selected: unknown[];
        select(value: unknown): void;
    }
</script>

<script lang="ts">
    import { setContext, createEventDispatcher } from "svelte";
    import { writable } from "svelte/store";
    import {
        Popover,
        PopoverButton,
        PopoverPanel,
        Transition,
    } from "@rgossiaux/svelte-headlessui";

    export let buttonText: string;
    export let value: unknown[];

    let dispatch = createEventDispatcher<{ change: any }>();

    let store = writable<MultiSelectStore>({
        selected: value,
        select(value: unknown) {
            dispatch("change", value);
        },
    });
    setContext(MULTI_SELECT_CONTEXT, store);

    $: store.update((obj) => {
        return {
            ...obj,
            selected: value,
        };
    });
</script>

<Popover class="relative">
    <div class={$$props.class}>
        <PopoverButton 
            class="rounded px-3 py-1.5 bg-gray-100 hover:bg-gray-200 transition-colors"
        >
            {buttonText}
        </PopoverButton>
    </div>
    <Transition
            leave="transition ease-in duration-100"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
    >
        <PopoverPanel 
            as="ul" 
            class="absolute z-10 mt-1 bg-white border border-gray-200 rounded drop-shadow px-1 py-1" 
            role="listbox"
        >
            <slot></slot>
        </PopoverPanel>
    </Transition>
</Popover>

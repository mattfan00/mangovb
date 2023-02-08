<script lang="ts">
    import { getContext } from "svelte";
    import type { Writable } from "svelte/store";
    import type { MultiSelectStore } from "./MultiSelect.svelte";
    import { MULTI_SELECT_CONTEXT } from "./MultiSelect.svelte";
    
    export let value: any;

    let element: HTMLElement;
    let store = getContext<Writable<MultiSelectStore>>(MULTI_SELECT_CONTEXT);

    const onClick = (_: any) => {
        let local: any[] = [...$store.selected];
        if (local.includes(value)) {
            local.splice(local.indexOf(value), 1);
        } else {
            local = [...local, value];
        }

        $store.select(local);
    }

    const onKeyDown = (event: KeyboardEvent) => {
        if (['Enter', 'Space'].includes(event.code)) {
			event.preventDefault();
            element.click();
		}
    }

</script>

<li 
    bind:this={element}
    class="whitespace-nowrap px-2 py-1 rounded hover:bg-gray-100 transition-colors cursor-pointer"
    role="option" 
    tabindex="0" 
    aria-selected="false"
    on:click={onClick}
    on:keydown={onKeyDown}
    on:keyup
    on:keypress
>
    <slot></slot>
</li>

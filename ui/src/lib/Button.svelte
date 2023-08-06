<script lang="ts">
    import classnames from "classnames";

    type buttonTypes = "default"|"discord";
    export let type: buttonTypes = "default";
    export let style = "";
    export let href = null;

    let classname = classnames(
        "rounded transition-colors px-3 py-1.5",
        { "border border-gray-200 hover:bg-gray-100 active:bg-gray-200": type === "default" },
        { "bg-[#5865F2] text-white": type === "discord" },
        style,
    );

    const tag = !href ? "button" : "a";
</script>

<svelte:element 
    this={tag} 
    {href} 
    class={classname}
    {...$$restProps}
>
    <slot></slot>
</svelte:element>

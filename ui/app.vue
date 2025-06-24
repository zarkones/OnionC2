<template>
    <NuxtLayout>
    </NuxtLayout>
</template>

<style>

.liquid-glass {
    /* Semi-transparent background for the glass effect */
    background-color: rgba(0, 0, 0, 0.1) !important;

    /* Blur the content behind the element */
    backdrop-filter: blur(10px);

    /* Optional: Subtle border for definition */
    border: 1px solid rgba(0, 0, 0, 0.5);

    /* Optional: Slight shadow for depth */
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);

    /* Optional: Rounded corners for a smoother look */
    border-radius: 10px;
}

/* Fallback for browsers that don’t support backdrop-filter */
@supports not (backdrop-filter: blur(10px)) {
    .liquid-glass {
        background-color: rgba(0, 0, 0, 0.5);
        /* Slightly more opaque */
    }
}
</style>
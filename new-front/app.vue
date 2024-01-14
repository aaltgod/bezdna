<script setup lang="ts">
const route = useRoute();

useHead({
  htmlAttrs: {
    lang: "en",
  },
  title: "Bezdna",
  meta: [{ name: "description", content: "Полезная тулза" }],
});
</script>

<template>
  <NuxtLoadingIndicator />
  <NuxtPage />
</template>

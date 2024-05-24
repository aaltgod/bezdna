<script lang="ts" setup>
const route = useRoute();

useHead({
  htmlAttrs: {
    lang: "en",
  },
  title: "ALPHA",
  meta: [{name: "description", content: "Полезная тулза"}],
});
</script>

<template>
  <NuxtLoadingIndicator/>
  <NuxtPage/>
</template>

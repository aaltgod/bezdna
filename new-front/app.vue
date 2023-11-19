<script setup lang="ts">
const route = useRoute()

useHead({
  htmlAttrs: {
    lang: 'en',
  },
  charset: 'utf-8',
  title: 'Bezdna',
  meta: [
    { name: 'description', content: 'Полезная тулза' },
  ]
})
</script>

<template>
  <NuxtLoadingIndicator />
  <UCard>
    <template #header>
      <nav class="flex align-center gap-4 p-4">
        <UButton>
        <NuxtLink to="/" class="n-link-base">
          Home
        </NuxtLink>
        </UButton>
        <UButton>
        <NuxtLink @click="$router.push(`/services/${(Math.random() * 100).toFixed()}`)" class="n-link-base">
          Services
        </NuxtLink>
        </UButton>

      </nav>
    </template>
    <NuxtPage />
    <template #footer>
    footer
    </template>
  </UCard>

</template>


<script setup>
const isOpen = ref(false)

const links = [{
  label: 'Profile',
  avatar: {
    src: 'https://avatars.githubusercontent.com/u/739984?v=4'
  },
  badge: 100
}, {
  label: 'Services',
  icon: 'i-heroicons-home',
  to: '/services/2'
}, {
  label: 'Vertical Navigation',
  icon: 'i-heroicons-chart-bar',
  to: '/navigation/vertical-navigation'
}, {
  label: 'Command Palette',
  icon: 'i-heroicons-command-line',
  to: '/navigation/command-palette'
}]
</script>

<template>
  <UVerticalNavigation :links="links" />
  <div>
    <UButton label="Open" @click="isOpen = true" />
    <UModal v-model="isOpen" fullscreen>
      <UCard
          :ui="{
          base: 'h-full flex flex-col',
          rounded: '',
          divide: 'divide-y divide-gray-100 dark:divide-gray-800',
          body: {
            base: 'grow'
          }
        }"
      >
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
              Modal
            </h3>
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" class="-my-1" @click="isOpen = false" />
          </div>
        </template>
        <Placeholder class="h-full" />
      </UCard>
    </UModal>
  </div>
</template>

<template>
  <UApp>
    <div class="shadow-sm px-4 py-3">
      <div class="container mx-auto flex items-center gap-4">
        <h1 class="text-2xl font-bold">Speech shadowing</h1>
        <Breadcrumb :links="breadcrumbLinks" />
      </div>
    </div>

    <UMain>
      <NuxtLayout>
        <NuxtPage />
      </NuxtLayout>
    </UMain>

    <UFooter />
  </UApp>
</template>

<script setup lang="ts">
const route = useRoute()

const breadcrumbLinks = computed(() => {
  const links: Array<{ label: string; to?: string }> = [{ label: 'Home', to: '/' }]
  
  if (route.path.startsWith('/speech-exercises')) {
    links.push({ label: 'Speech Exercises', to: '/speech-exercises' })
    
    if (route.params.id) {
      const title = route.query.title as string || `Exercise ${route.params.id}`
      links.push({ label: title })
    }
  }
  
  return links
})
</script>


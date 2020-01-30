<template>
  <div>
    <h1 class="mb-3">Platforms</h1>

    <div
      v-for="(platform, key) in sortedPlatforms"
      :key="key"
    >
      <router-link :to="{ name: 'tasks', params: { platform }}">{{ platform }}</router-link>
    </div>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator'

  import Api from '@/services/api'

  @Component
  export default class Platforms extends Vue {
    private platforms: string[] = []

    /**
     * Methods
     */
    private mounted() {
      this.setPlatforms()
    }

    private setPlatforms(): void {
      Api().get('/platforms').then(response => {
        this.platforms = response.data
      })
    }

    /**
     * Computed
     */
    get sortedPlatforms(): string[] {
      const platforms = this.platforms

      return platforms.sort((a: string, b: string) => a > b ? 1 : -1)
    }
  }
</script>

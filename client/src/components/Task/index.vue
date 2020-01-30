<template>
  <div>
    <h1 class="mb-3">Tasks</h1>
    <div
      v-for="(task, key) in sortedTasks"
      :key="key"
    >
      <router-link :to="getRouterTaskLink(task)">{{ task }}</router-link>
    </div>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator'

  import Api from '@/services/api'

  @Component
  export default class Tasks extends Vue {
    private tasks: string[] = []

    /**
     * Methods
     */
    private mounted() {
      this.setTasks()
    }

    private getRouterTaskLink(task: string): any {
      return {
        name: 'new_job',
        params: {
          platform: this.$route.params.platform,
          task,
        },
      }
    }

    private setTasks(): void {
      Api().get('/tasks').then(response => {
        this.tasks = response.data
      })
    }

    /**
     * Computed
     */
    get sortedTasks(): string[] {
      const tasks = this.tasks

      return tasks.sort((a: string, b: string) => a > b ? 1 : -1)
    }
  }
</script>

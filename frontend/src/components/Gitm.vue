<template>
  <div class="gitm">
    <div v-for="commit in commits" v-bind:key="commit.object.hash" v-text="commit">Initializing...</div>
    <button :disabled="noMoreCommits" @click="get">Get more commits</button>
  </div>
</template>

<script>
export default {
  name: 'Gitm',
  props: {
    msg: String,
  },
  data() {
    return {
      logIterator: {},
      commits: [],
    };
  },
  apollo: {
  },
  computed: {
    noMoreCommits() {
      return !(this.logIterator && this.logIterator.pointers && this.logIterator.pointers.length);
    },
  },
  methods: {
    async init() {
      const result = await this.$apollo.query({
        query: require('../graphql/InitLog.gql'),
      })
      this.logIterator = result.data.init;
    },
    async get() {
      if (!this.logIterator.pointers.length) return;
      const result = await this.$apollo.query({
        query: require('../graphql/GetLog.gql'),
        variables: { log_iterator: this.logIterator },
      })
      this.logIterator = result.data.get;
      const commits = this.logIterator.commits;
      this.commits = [...this.commits, ...commits];
      this.logIterator = { ...this.logIterator, commits: [] };
    },
  },
  created() {
    // eslint-disable-next-line no-console
    console.log(0);
    this.init().then(() => this.get());
  }
}
</script>

<style scoped>
</style>
